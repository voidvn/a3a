package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENV"`
	HTTP        struct {
		Port         string        `mapstructure:"PORT"`
		ReadTimeout  time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
		WriteTimeout time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
		IdleTimeout  time.Duration `mapstructure:"HTTP_IDLE_TIMEOUT"`
	} `mapstructure:",squash"`
	Database struct {
		URL string `mapstructure:"DATABASE_URL"`
	} `mapstructure:",squash"`
	Admin struct {
		Enabled bool   `mapstructure:"ADMIN_ENABLED"`
		AppKey  string `mapstructure:"ADMIN_APP_KEY"`
	} `mapstructure:",squash"`
	JWT struct {
		Secret string `mapstructure:"JWT_SECRET"`
	} `mapstructure:",squash"`
}

func Load() (*Config, error) {
	// Set default values
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HTTP_READ_TIMEOUT", "30s")
	viper.SetDefault("HTTP_WRITE_TIMEOUT", "30s")
	viper.SetDefault("HTTP_IDLE_TIMEOUT", "120s")
	viper.SetDefault("ADMIN_ENABLED", true)

	// Load .env file if it exists
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func GetString(key, fallback string) string {
	if viper.GetString(key) == "" {
		return fallback
	}
	return viper.GetString(key)
}
