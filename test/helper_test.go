package main

import (
	"testing"

	"github.com/kaankoken/binance-quick-go/helper"
)

const (
	fileName  string = "telegram_key_success_test"
	extension string = "yaml"
)

func TestReadTelegramKeysSuccess(t *testing.T) {
	expected := helper.ReadTelegramKeys(fileName, extension)

	if expected["token"] == nil || expected["channelid"] == nil {
		t.Fatal("Could not parse the YAML file")
	}
}

func TestReadTelegramKeysError(t *testing.T) {
	expected := helper.ReadTelegramKeys("telegram_key_error_test", "yaml")

	if expected["token"] != nil || expected["channelid"] != nil {
		t.Fatal("Could parse the YAML file")
	}
}
