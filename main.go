package main

import (
	observerbot "github.com/kaankoken/binance-quick-go/observer-bot"
	telegrambot "github.com/kaankoken/binance-quick-go/telegram-bot"
)

func main() {
	telegrambot.Run()
	observerbot.Run()
}
