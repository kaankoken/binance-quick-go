package observerbot

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/evsamsonov/trading-indicators/indicator"
	"github.com/evsamsonov/trading-timeseries/timeseries"
	"github.com/kaankoken/binance-quick-go/helper"
	"github.com/kaankoken/binance-quick-go/observer-bot/models"
	telegrambot "github.com/kaankoken/binance-quick-go/telegram-bot"
	"github.com/xuri/excelize/v2"

	log "github.com/sirupsen/logrus"
)

var (
	keys      *models.ApiKeyModel
	symbols   *[]models.SymbolModel
	intervals []string
	crossTime []int                          = []int{03, 07, 11, 15, 19, 23}
	volume    map[string]float64             = make(map[string]float64)
	results   map[string][]models.KlineModel = make(map[string][]models.KlineModel)
)

const (
	fileName      string = "api_key"
	extension     string = "yaml"
	quoteAsset    string = "USDT"
	averageVolume uint16 = 51
	ma14          uint16 = 15
	ma9           uint16 = 10
)

func init() {
	keys = models.ToApiModel(helper.ReadApiKey(fileName, extension))
	symbols = &[]models.SymbolModel{}
	intervals = []string{"1h", "4h", "1d"}
}

func isCrossTime() bool {
	_localTime := time.Now().Hour()

	for _, i := range crossTime {
		if _localTime == i {
			return true
		}
	}

	return false
}

func getKlines(ctx context.Context, client *futures.Client, wg *sync.WaitGroup, arg models.KlineArguments) {
	result, err := client.NewKlinesService().Symbol(arg.Symbol).
		Interval(arg.Interval).Limit(arg.Limit).Do(ctx)

	helper.CheckError(err)

	var key string
	if arg.Key_Write_Option {
		key = arg.Symbol + "-" + arg.Interval
	} else {
		key = arg.Symbol
	}

	arg.Ch <- models.KlineMainModel{Key: key, Data: models.ToKlineModel(result)}
	defer wg.Done()
}

func hourlyRequest(ctx context.Context, client *futures.Client, wg *sync.WaitGroup, arg models.KlineArguments) {
	for _, v := range *symbols {
		wg.Add(1)
		arg.Symbol = v.Symbol
		go getKlines(ctx, client, wg, arg)
	}
}

func fourHourlyRequest(ctx context.Context, client *futures.Client, wg *sync.WaitGroup, arg models.KlineArguments) {
	for _, v := range *symbols {
		wg.Add(1)
		arg.Symbol = v.Symbol

		go getKlines(ctx, client, wg, arg)
	}
}

func CalculateFiftyDaysOfAverageVolume() {
	getSymbols()

	ctx := context.Background()
	futuresClient := binance.NewFuturesClient(keys.Api, keys.Secret)

	channel := make(chan models.KlineMainModel)

	var wg sync.WaitGroup
	arg := models.KlineArguments{Ch: channel, Interval: intervals[2], Limit: int(averageVolume), Key_Write_Option: false}

	for _, v := range *symbols {
		wg.Add(1)
		arg.Symbol = v.Symbol
		go getKlines(ctx, futuresClient, &wg, arg)
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	currentTime := time.Now().UnixMilli()
	for value := range channel {
		if len(*value.Data) > 0 && currentTime-(*value.Data)[len(*value.Data)-1].CloseTime < 0 {
			*value.Data = (*value.Data)[:len(*value.Data)-1]
		}

		var _volume float64
		for _, v := range *value.Data {
			_volume += float64(v.Volume)
		}

		volume[value.Key] = _volume / 50
	}
}

func filterIntervals(data map[string][]models.KlineModel, f func(string) bool) map[string][]models.KlineModel {
	fltd := make(map[string][]models.KlineModel, 0)

	for key, value := range data {
		if f(key) {
			fltd[key] = value
		}
	}

	return fltd
}

func calculateMovingAverage(interval string) {
	res := filterIntervals(results, func(s string) bool {
		return strings.Contains(s, interval)
	})

	ema14 := make(map[string]float64, 0)
	ema9 := make(map[string]float64, 0)

	for key, value := range res {
		series := timeseries.New()

		for _, v := range value {
			candle := timeseries.NewCandle(time.UnixMilli(v.CloseTime))
			candle.Open = v.Open
			candle.Close = v.Close
			candle.High = v.High
			candle.Low = v.Low
			candle.Volume = int64(v.Volume)

			err := series.AddCandle(candle)
			if err != nil {
				log.Info("Failed to add candle: %v\n", err)
			}
		}

		atrIndicator14, err := indicator.NewExponentialMovingAverage(series, 2)
		helper.CheckError(err)

		newTs, err := series.RemoveCandles(5, nil)
		helper.CheckError(err)

		atrIndicator9, err := indicator.NewExponentialMovingAverage(newTs, 2)
		helper.CheckError(err)

		ema14[key] = atrIndicator14.Calculate(1)
		ema9[key] = atrIndicator9.Calculate(1)
	}

	message := "-----------------------------------------------------" + "\n" + interval + "\n" + "Coin\t+Volume\tEMA14\tEMA9\n"

	for i, v := range ema14 {
		message += i + " " + fmt.Sprintf("%f", volume[i]) + " " + fmt.Sprintf("%f", v) + " " + fmt.Sprintf("%f", ema9[i]) + "\n"
	}

	message += "-----------------------------------------------------\n"

	telegrambot.SendNewMessage(message)

	//createXlxsFile(ema14, ema9, interval)
}

func Run() {
	getSymbols()

	if volume == nil || len(volume) > 0 {
		CalculateFiftyDaysOfAverageVolume()
	}

	ctx := context.Background()
	futuresClient := binance.NewFuturesClient(keys.Api, keys.Secret)

	channel := make(chan models.KlineMainModel)

	var wg sync.WaitGroup

	if isCrossTime() {
		hourlyRequest(ctx, futuresClient, &wg, models.KlineArguments{Ch: channel, Interval: intervals[0], Limit: int(ma14), Key_Write_Option: true})
		fourHourlyRequest(ctx, futuresClient, &wg, models.KlineArguments{Ch: channel, Interval: intervals[1], Limit: int(ma14), Key_Write_Option: true})
	} else {
		hourlyRequest(ctx, futuresClient, &wg, models.KlineArguments{Ch: channel, Interval: intervals[0], Limit: int(ma14), Key_Write_Option: true})
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	currentTime := time.Now().UnixMilli()
	// Removing last element if last candle stick is not closed
	for value := range channel {
		if len(*value.Data) > 0 && currentTime-(*value.Data)[len(*value.Data)-1].CloseTime < 0 {
			*value.Data = (*value.Data)[:len(*value.Data)-1]
		}

		results[value.Key] = *value.Data
	}

	if isCrossTime() {
		// 1 HR
		calculateMovingAverage(intervals[0])
		// 4 HR
		calculateMovingAverage(intervals[1])
	} else {
		// 1 HR
		calculateMovingAverage(intervals[0])
	}
}

func GetLatestResult() {
	// TODO: get saved latest from DB
}

func getSymbols() {
	if symbols != nil && len(*symbols) > 0 {
		return
	}
	/*
		TODO: if data exist in db
		call dto model to model
		symbols = models.ToSymbolModel()
		return
	*/

	// TODO: if does not exist
	client := binance.NewFuturesClient(keys.Api, keys.Secret)

	_symbols, err := client.NewExchangeInfoService().Do(context.Background())
	helper.CheckError(err)

	filteredSymbols := models.FilterSymbols(_symbols, quoteAsset)
	symbols = models.ToSymbolModel(filteredSymbols)
	// TODO: Save to sql from sql module
}

func createXlxsFile(ema14 map[string]float64, ema9 map[string]float64, interval string) {
	f := excelize.NewFile()
	categories := map[string]string{
		"A0": "Symbol", "B0": "Interval", "C0": "50-Days-Average-Volume",
		"D0": "EMA14", "E0": "EMA9"}

	for k, v := range categories {
		f.SetCellValue("Sheet1", k, v)
	}

	i := 1
	for k, v := range ema14 {
		f.SetCellValue("Sheet1", "A"+string(rune(i)), k)
		f.SetCellValue("Sheet1", "C"+string(rune(i)), volume[k])
		f.SetCellValue("Sheet1", "D"+string(rune(i)), v)
		f.SetCellValue("Sheet1", "E"+string(rune(i)), ema9[k])

		i++
	}
	f.SetCellValue("Sheet1", "B1", interval)

	err := f.SaveAs(string(rune(time.Now().UnixMilli())) + interval + ".xlsx")
	helper.CheckError(err)
}
