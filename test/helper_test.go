package main

import (
	"testing"

	"github.com/kaankoken/binance-quick-go/helper"
)

const (
	fileSuccessName string = "telegram_key_success_test"
	fileErrorName   string = "telegram_key_success_test"
	extension       string = "yaml"
)

func TestReadTelegramKeysSuccess(t *testing.T) {
	expected := helper.ReadTelegramKeys(fileSuccessName, extension)

	if expected["token"] == nil || expected["channelid"] == nil {
		t.Fatal("Could not parse the YAML file")
	}
}

func TestReadTelegramKeysError(t *testing.T) {
	expected := helper.ReadTelegramKeys(fileErrorName, extension)

	if expected["token"] != nil || expected["channelid"] != nil {
		t.Fatal("Could parse the YAML file")
	}
}
