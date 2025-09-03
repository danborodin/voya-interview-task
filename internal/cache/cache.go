package cache

import (
	"errors"
)

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

var (
	ErrCacheMiss         = errors.New("no such key in cache")
	ErrInvalidCacheValue = errors.New("invalid cache value")
	ErrTTLExpired        = errors.New("ttl for this key/value expired")
)
