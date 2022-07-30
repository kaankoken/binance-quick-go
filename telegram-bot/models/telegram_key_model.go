package models

import (
	"encoding/json"
	helper "legolas/binance-trading/helper"
)

type TelegramKeyModel struct {
	Token     string
	ChannelId int64
}

func ToModel(data map[string]interface{}) *TelegramKeyModel {
	model := TelegramKeyModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
