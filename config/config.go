package config

import (
	"github.com/spf13/viper"
)

func initConfig() {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetString(s string) string {
	return viper.GetString(s)
}

func GetInt(i string) int {
	return viper.GetInt(i)
}
