package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func Setup(cfg Config, options ...func(*redis.Options)) error {
	if !cfg.Enable {
		return nil
	}

	option := cfg.toOption()

	for _, fn := range options {
		fn(option)
	}

	rdb = redis.NewClient(option)

	return rdb.Ping(context.Background()).Err()
}

func Use(fn func(client *redis.Client) error, ignore bool) error {
	if fn == nil {
		return nil
	}

	if !ignore && rdb == nil {
		return ErrClientNotInit
	}

	return fn(rdb)
}
