package redis

import "github.com/redis/go-redis/v9"

type Config struct {
	Enable   bool   `mapstructure:"enable"`
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func (c Config) toOption() *redis.Options {
	return &redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	}
}

func GetDefaultConfig() Config {
	return Config{
		Enable:   true,
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}
}
