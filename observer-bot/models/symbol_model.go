package models

import (
	"encoding/json"

	"github.com/adshao/go-binance/v2/futures"
	"github.com/kaankoken/binance-quick-go/helper"
	"github.com/oklog/ulid/v2"
)

type SymbolDtoModel struct {
	Id     ulid.ULID
	Symbol string
}
type SymbolModel struct {
	Symbol string
}

func FilterSymbols(data *futures.ExchangeInfo, quoteAsset string) *[]SymbolDtoModel {
	symbols := make([]SymbolDtoModel, 0)

	for _, symbol := range data.Symbols {
		if symbol.QuoteAsset == quoteAsset {
			symbols = append(symbols, SymbolDtoModel{Id: ulid.Make(), Symbol: symbol.Symbol})
		}
	}

	return &symbols
}

func ToSymbolModel(data interface{}) *[]SymbolModel {
	model := []SymbolModel{}
	body, err := json.Marshal(data)
	helper.CheckError(err)

	err = json.Unmarshal(body, &model)
	helper.CheckError(err)

	return &model
}
