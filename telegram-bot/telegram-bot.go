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
	keys = models.ToTelegramModel(helper.ReadTelegramKeys(fileName, extension))
}

func Run() {
	bot, err := tgbotapi.NewBotAPI(keys.Token)

	helper.CheckError(err)

	bot.Debug = true
	Bot = bot

	log.Printf("Authorized on account %s", bot.Self.UserName)
}

func SendNewMessage(message string) {
	newMessage := tgbotapi.NewMessage(keys.ChannelId, message)
	_, err := Bot.Send(newMessage)
	helper.CheckError(err)
}
