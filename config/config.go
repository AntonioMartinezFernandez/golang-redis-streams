package config

import (
	"context"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	AppEnv     string `env:"APP_ENV"`
	AppVersion string `env:"APP_VERSION"`

	LogLevel string `env:"LOG_LEVEL"`

	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`

	ConsumerGroup string `env:"CONSUMER_GROUP"`
}

func LoadEnvConfig() Config {
	var config Config

	godotenv.Load(".env")
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
