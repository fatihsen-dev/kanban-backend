package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port      string `mapstructure:"port" yaml:"port"`
	DBUrl     string `mapstructure:"db_url" yaml:"db_url"`
	JWTSecret string `mapstructure:"jwt_secret" yaml:"jwt_secret"`
}

func Read() *AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$PWD/config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	var appConfig AppConfig
	err = viper.Unmarshal(&appConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshalling config: %w", err))
	}

	return &appConfig
}
