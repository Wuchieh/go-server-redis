package redis

import "errors"

var (
	ErrClientNotInit = errors.New("redis client has not yet been initialized")
)
