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

	RedisHost string `env:"REDIS_HOST"`
	RedisPort string `env:"REDIS_PORT"`

	Stream      string `env:"STREAM"`
	StreamGroup string `env:"STREAM_GROUP"`
}

func LoadEnvConfig() Config {
	goDotEnvVariableString("APP_ENV")
	goDotEnvVariableString("APP_VERSION")

	goDotEnvVariableString("REDIS_HOST")
	goDotEnvVariableString("REDIS_PORT")

	goDotEnvVariableString("STREAM")
	goDotEnvVariableString("STREAM_GROUP")

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
