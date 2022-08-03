package models

import (
	"encoding/json"

	"github.com/kaankoken/binance-quick-go/helper"
)

type ApiKeyModel struct {
	Api    string
	Secret string
}

func ToModel(data interface{}) *ApiKeyModel {
	model := ApiKeyModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
