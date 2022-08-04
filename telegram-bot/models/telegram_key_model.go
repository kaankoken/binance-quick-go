package models

import (
	"encoding/json"

	"github.com/kaankoken/binance-quick-go/helper"
)

type TelegramKeyModel struct {
	Token     string
	ChannelId int64
}

func ToTelegramModel(data interface{}) *TelegramKeyModel {
	model := TelegramKeyModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
