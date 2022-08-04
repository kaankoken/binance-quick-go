package helper

import (
	"fmt"

	"github.com/spf13/viper"
)

func setupViper(fileName string, extension string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(extension)
	viper.AddConfigPath(".")
	viper.AddConfigPath("$PWD/")

	err := viper.ReadInConfig()
	CheckError(err)

	viper.ReadInConfig()
}

func CheckError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

func ReadApiKey(fileName string, extension string) map[string]interface{} {
	setupViper(fileName, extension)

	return viper.GetStringMap("key")
}

func ReadTelegramKeys(fileName string, extension string) map[string]interface{} {
	setupViper(fileName, extension)

	return viper.GetStringMap("key")
}
