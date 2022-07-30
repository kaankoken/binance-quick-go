package telegrambot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	helper "legolas/binance-trading/helper"
	model "legolas/binance-trading/telegram-bot/models"
)

var (
	keys *model.TelegramKeyModel
)

func init() {
	keys = model.ToModel(helper.ReadTelegramKeys())
}

func TestKeys() {
	log.Println(keys)
}

func TelegramBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(keys.Token)

	helper.CheckError(err)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// ignore any non-Message updates
		if update.Message == nil {
			continue
		}

		// ignore any non-command Messages
		if !update.Message.IsCommand() {
			continue
		}

		var text string
		switch update.Message.Command() {
		case "help":
			text = "I understand /sayhi and /status."
		case "sayhi":
			text = "Hi :)"
		case "status":
			text = "I'm ok."
		case "latest":
			text = "Get latest result"
		default:
			text = "I don't know that command"
		}

		SendNewMessage(bot, text)
	}

	return bot
}

func SendNewMessage(bot *tgbotapi.BotAPI, message string) {
	newMessage := tgbotapi.NewMessage(keys.ChannelId, "")

	_, err := bot.Send(newMessage)
	helper.CheckError(err)
}
