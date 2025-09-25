package config

import (
	"log"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No .env file found, using env vars: %v", err)
	}
}
