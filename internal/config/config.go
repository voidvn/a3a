package config

import "github.com/spf13/viper"

func Load() error {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	return viper.ReadInConfig()
}

func GetString(key, fallback string) string {
	if viper.GetString(key) == "" {
		return fallback
	}
	return viper.GetString(key)
}
