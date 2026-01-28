package config

import (
	"fmt"
	"reflect"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port       string `mapstructure:"PORT" validate:"required"`
	JWTSecret  string `mapstructure:"JWT_SECRET" validate:"required"`
	ClientUrl  string `mapstructure:"CLIENT_URL" validate:"required"`
	DBPort     string `mapstructure:"DB_PORT" validate:"required"`
	DBHost     string `mapstructure:"DB_HOST" validate:"required"`
	DBUser     string `mapstructure:"DB_USER" validate:"required"`
	DBPassword string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName     string `mapstructure:"DB_NAME" validate:"required"`
	DBUrl      string `mapstructure:"DB_URL"`
}

func Read() *AppConfig {
	_ = godotenv.Load()
	viper.AutomaticEnv()

	var cfg AppConfig
	BindAllEnv(&cfg)

	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("config unmarshal error: %w", err))
	}

	if err := validation.Validate(&cfg); err != nil {
		panic(fmt.Errorf("config validation error: %w", err))
	}

	cfg.DBUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	return &cfg
}

func BindAllEnv(cfg any) {
	t := reflect.TypeOf(cfg)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			viper.BindEnv(tag)
		}
	}
}
