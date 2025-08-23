package redis

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func getCtx(ctx context.Context) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx
}

type Cache[T any] struct {
	mux sync.RWMutex
	d   time.Duration
	t   *time.Timer
	key string
}

func (c *Cache[T]) Get(ctx context.Context) (T, error) {
	c.mux.RLock()
	defer c.mux.RUnlock()

	var result T

	ctx = getCtx(ctx)

	err := Use(func(client *redis.Client) error {
		record := client.Get(ctx, c.key)

		b, err := record.Bytes()
		if err != nil {
			return err
		}

		return json.Unmarshal(b, &result)
	}, false)

	return result, err
}

func (c *Cache[T]) SimpleGet() T {
	result, _ := c.Get(getCtx(nil))
	return result
}

func (c *Cache[T]) Set(ctx context.Context, value T) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	ctx = getCtx(ctx)

	return Use(func(client *redis.Client) error {
		marshal, err := json.Marshal(value)
		if err != nil {
			return err
		}
		if err := client.
			Set(ctx, c.key, marshal, c.d).
			Err(); err != nil {
			return err
		}
		return nil
	}, false)
}

func (c *Cache[T]) SimpleSet(value T) {
	_ = c.Set(getCtx(nil), value)
}

func (c *Cache[T]) Delete(ctx context.Context) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	ctx = getCtx(ctx)

	return Use(func(client *redis.Client) error {
		return client.Del(ctx, c.key).Err()
	}, false)
}

func (c *Cache[T]) SimpleDelete() {
	_ = c.Delete(getCtx(nil))
}

func NewCache[T any](key string, expiry time.Duration, data ...T) *Cache[T] {
	c := Cache[T]{
		d:   expiry,
		key: key,
	}

	if len(data) > 0 {
		c.SimpleSet(data[0])
	}

	return &c
}
