package telegrambot

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gcloudbucket "github.com/kaankoken/binance-quick-go/gcloud-bucket"
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
	keys = models.ToTelegramModel(helper.ReadApiKey(fileName, extension))
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

func SendFile(path string) {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)

	helper.CheckError(err)
	gcloudbucket.Uploader.UploadFile(*file, file.Name())
}
