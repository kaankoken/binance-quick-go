package main

import (
	"testing"

	"github.com/kaankoken/binance-quick-go/helper"
)

const (
	fileSuccessName string = "api_key_success_test"
	fileErrorName   string = "api_key_error_test"
	extension       string = "yaml"
)

func TestReadTelegramKeysSuccess(t *testing.T) {
	expected := helper.ReadApiKey(fileSuccessName, extension)

	if expected["api"] == nil || expected["secret"] == nil {
		t.Fatal("Could not parse the YAML file")
	}
}

func TestReadTelegramKeysError(t *testing.T) {
	expected := helper.ReadApiKey(fileErrorName, extension)

	if expected["api"] != nil || expected["secret"] != nil {
		t.Fatal("Could parse the YAML file")
	}
}
