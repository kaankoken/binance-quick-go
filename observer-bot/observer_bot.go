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
)

var (
	keys *models.ApiKeyModel
)

const (
	fileName  string = "api_key"
	extension string = "yaml"
)

func init() {
	keys = models.ToModel(helper.ReadApiKey(fileName, extension))
}

func TestKeys() {
	log.Println(keys)
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

func GetSymbols() []string {
	//TODO: check if data exist in DB use it else get it
	client := binance.NewClient(keys.Api, keys.Secret)

	symbols, err := client.NewExchangeInfoService().Do(context.Background())
	helper.CheckError(err)

	log.Println(symbols)

	return []string{"sadahjsdkj"}
}

func SaveLatestResult() {
	// TODO: save into DB
}
