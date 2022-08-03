package observerbot

import (
	"context"
	"fmt"

	"log"
	"sync"

	"github.com/adshao/go-binance/v2"
	"github.com/adshao/go-binance/v2/futures"
	"github.com/kaankoken/binance-quick-go/helper"
	"github.com/kaankoken/binance-quick-go/observer-bot/models"
	"github.com/oklog/ulid/v2"
)

var (
	keys    *models.ApiKeyModel
	symbols []models.SymbolModel
)

const (
	fileName   string = "api_key"
	extension  string = "yaml"
	quoteAsset string = "USDT"
)

func init() {
	keys = models.ToModel(helper.ReadApiKey(fileName, extension))
	symbols = []models.SymbolModel{}
}

func getKlines(ctx context.Context, client *futures.Client, wg *sync.WaitGroup, ch chan<- []string, symbol string, interval string) {
	result, err := client.NewKlinesService().Symbol(symbol).
		Interval(interval).Do(ctx)

	fmt.Println(result)
	helper.CheckError(err)

	// TODO: pass result to channel
	// TODO create basic model for channel
	ch <- []string{"result"}
	defer wg.Done()
}

func calculateFiftyDaysOfVolume() {

}

func calculateMovingAverage() {

}

func Run(symbols []string, intervals [2]string) {
	ctx := context.Background()
	futuresClient := binance.NewFuturesClient(keys.Api, keys.Secret)

	channel := make(chan []string)

	var wg sync.WaitGroup

	for _, v := range symbols {
		for _, i := range intervals {
			wg.Add(1)
			go getKlines(ctx, futuresClient, &wg, channel, v, i)
		}
	}

	go func() {
		wg.Wait()
		close(channel)
	}()

	for kline := range channel {
		log.Println(kline)
	}
}

func GetLatestResult() {
	// TODO: get saved latest from DB
}

func GetSymbols() {
	//TODO: check if data exist in DB use it else get it
	client := binance.NewFuturesClient(keys.Api, keys.Secret)

	_symbols, err := client.NewExchangeInfoService().Do(context.Background())
	helper.CheckError(err)

	for _, symbol := range _symbols.Symbols {
		if symbol.QuoteAsset == quoteAsset {
			symbols = append(symbols, models.SymbolModel{Id: ulid.Make(), Symbol: symbol.Symbol})
		}
	}

	log.Println(symbols)
}

func SaveLatestResult() {
	// TODO: save into DB
}
