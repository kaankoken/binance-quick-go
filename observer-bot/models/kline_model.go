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
	Open      float64 `json:",string"`
	Close     float64 `json:",string"`
	High      float64 `json:",string"`
	Low       float64 `json:",string"`
	Volume    float64 `json:",string"`
	CloseTime int64
}

type KlineArguments struct {
	Symbol           string
	Interval         string
	Ch               chan<- KlineMainModel
	Limit            int
	Key_Write_Option bool
}

func ToKlineModel(data interface{}) *[]KlineModel {
	model := []KlineModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
