package telegrambot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kaankoken/binance-quick-go/helper"
	"github.com/kaankoken/binance-quick-go/telegram-bot/models"
)

var (
	keys *models.TelegramKeyModel
	Bot  *tgbotapi.BotAPI
)

const (
	fileName  string = "telegram_key"
	extension string = "yaml"
)

func init() {
	keys = models.ToModel(helper.ReadTelegramKeys(fileName, extension))
}

func Run() {
	bot, err := tgbotapi.NewBotAPI(keys.Token)

	helper.CheckError(err)

	bot.Debug = true
	Bot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	botInitializationMessage := `
		Current bot commands
			/help
			/sayhi
			/status
			/latest
	`
	SendNewMessage(botInitializationMessage)

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

		SendNewMessage(text)
	}

}

func SendNewMessage(message string) {
	newMessage := tgbotapi.NewMessage(keys.ChannelId, message)
	_, err := Bot.Send(newMessage)
	helper.CheckError(err)
}
