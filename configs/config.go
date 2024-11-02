package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	LoggerLevel   string `env:"LOGGER_LEVEL"`
	IsDevelopment bool   `env:"IS_DEVELOPMENT"`

	HttpPort int `env:"HTTP_PORT"`

	RedisHost        string `env:"REDIS_HOST"`
	RedisPort        int    `env:"REDIS_PORT"`
	RedisPoolSize    int    `env:"REDIS_POOL_SIZE"`
	RedisIdleTimeout int    `env:"REDIS_IDLE_TIMEOUT"`
	RedisPoolTimeout int    `env:"REDIS_POOL_TIMEOUT"`
}

func BuildConfig() Config {
	var config Config

	godotenv.Load(".env")
	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		panic(err)
	}

	return config
}
