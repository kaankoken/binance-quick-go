package models

import (
	"encoding/json"
	helper "legolas/binance-trading/helper"
)

type ApiKeyModel struct {
	Api    string
	Secret string
}

func ToModel(data map[string]string) *ApiKeyModel {
	model := ApiKeyModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
