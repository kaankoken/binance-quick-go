package models

import (
	"encoding/json"

	"github.com/kaankoken/binance-quick-go/helper"
)

type KlineMainModel struct {
	Key  string
	Data *[]KlineModel
}

type KlineModel struct {
	Open      string
	Close     string
	High      string
	Low       string
	Volume    string
	CloseTime int64
}

type KlineArguments struct {
	Symbol   string
	Interval string
	Ch       chan<- KlineMainModel
	Limit    int
}

func ToKlineModel(data interface{}) *[]KlineModel {
	model := []KlineModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
