package config

import (
	"context"
	"os"

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
	goDotEnvVariableString("APP_ENV")
	goDotEnvVariableString("APP_VERSION")

	goDotEnvVariableString("LOG_LEVEL")

	goDotEnvVariableString("REDIS_HOST")
	goDotEnvVariableString("REDIS_PORT")

	goDotEnvVariableString("CONSUMER_GROUP")

	return buildConfig()
}

func goDotEnvVariableString(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}

func buildConfig() Config {
	var config Config
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
