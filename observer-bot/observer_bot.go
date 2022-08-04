package observerbot

import (
	"context"
	"log"
	"time"

	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kaankoken/binance-quick-go/helper"
	"github.com/kaankoken/binance-quick-go/observer-bot/models"
)

var (
	keys      *models.ApiKeyModel
	symbols   *[]models.SymbolModel
	intervals []string
)

const (
	fileName   string = "api_key"
	extension  string = "yaml"
	quoteAsset string = "USDT"
	maxLimit   uint16 = 300
	ma200      uint16 = 200
	ma100      uint16 = 100
	ma50       uint16 = 50
	ma20       uint16 = 20
)

func init() {
	keys = models.ToApiModel(helper.ReadApiKey(fileName, extension))
	symbols = &[]models.SymbolModel{}
	intervals = []string{"1h", "4h"}
}

func getKlines(ctx context.Context, client *futures.Client, wg *sync.WaitGroup, arg models.KlineArguments) {
	result, err := client.NewKlinesService().Symbol(arg.Symbol).
		Interval(arg.Interval).Limit(arg.Limit).Do(ctx)

	helper.CheckError(err)

	key := arg.Symbol + "-" + arg.Interval

	// TODO: pass result to channel
	// TODO create basic model for channel
	arg.Ch <- models.KlineMainModel{Key: key, Data: models.ToKlineModel(result)}
	defer wg.Done()
}

func calculateFiftyDaysOfVolume() {

}

func calculateMovingAverage() {

}

func Run() {
	GetSymbols()
	ctx := context.Background()
	futuresClient := binance.NewFuturesClient(keys.Api, keys.Secret)

	channel := make(chan models.KlineMainModel)

	var wg sync.WaitGroup

	for _, v := range *symbols {
		wg.Add(1)
		go getKlines(ctx, futuresClient, &wg, models.KlineArguments{Ch: channel, Symbol: v.Symbol, Interval: intervals[0], Limit: int(maxLimit)})
	}

	for _, v := range *symbols {
		wg.Add(1)
		go getKlines(ctx, futuresClient, &wg, models.KlineArguments{Ch: channel, Symbol: v.Symbol, Interval: intervals[1], Limit: int(ma200)})
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
		log.Println(value.Key, time.UnixMilli((*value.Data)[0].CloseTime), time.UnixMilli((*value.Data)[len(*value.Data)-1].CloseTime))
	}
}

func GetLatestResult() {
	// TODO: get saved latest from DB
}

func GetSymbols() {
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

	log.Println(len(*symbols))
	// TODO: Save to sql from sql module
}

func SaveLatestResult() {
	// TODO: save into DB
}
