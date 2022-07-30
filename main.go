package main

import (
	observerbot "legolas/binance-trading/observer-bot"
	telegrambot "legolas/binance-trading/telegram-bot"
	"log"
)

func main() {
	bot := telegrambot.TelegramBot()

	log.Println(bot)

	symbols := observerbot.GetSymbols()
	intervals := [...]string{"1d", "4h"}

	observerbot.Run(symbols, intervals)
}
