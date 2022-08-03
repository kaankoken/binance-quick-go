package main

import (
	"testing"

	"github.com/kaankoken/binance-quick-go/telegram-bot/models"
)

func TestToTelegramModelSuccess(t *testing.T) {
	var expected, result interface{}
	data := map[string]interface{}{
		"token":     "asdasdasdad",
		"channelId": 123182390123,
	}

	result = &models.TelegramKeyModel{Token: "asdasdasdasd", ChannelId: 12312312313}

	expected = models.ToModel(data)

	if expected.(*models.TelegramKeyModel) == result.(*models.TelegramKeyModel) {
		t.Fatal("Types are not EQUAL")
	}
}

func TestToTelegramModelError(t *testing.T) {
	data := map[string]interface{}{
		"token1":     "asdasdasdad",
		"channelId1": 123182390123,
	}

	expected := models.ToModel(data)

	if len(expected.Token) > 0 && expected.ChannelId != 0 {
		t.Fatal("Conversion is HAPPEND")
	}
}
