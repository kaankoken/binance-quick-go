package models

import (
	"github.com/oklog/ulid/v2"
)

type SymbolModel struct {
	Id     ulid.ULID
	Symbol string
}
