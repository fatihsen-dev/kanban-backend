package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port      string `mapstructure:"PORT"`
	DBUrl     string `mapstructure:"DB_URL"`
	JWTSecret string `mapstructure:"JWT_SECRET"`
	ClientUrl string `mapstructure:"CLIENT_URL"`
}

func Read() *AppConfig {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	viper.AddConfigPath("$PWD")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Error reading config: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("Error unmarshalling config: %w", err))
	}

	return &appConfig
}
