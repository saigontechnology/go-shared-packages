package cache

import (
	"context"
	"sync"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string, data interface{}, tag *string) error
	Set(ctx context.Context, key string, data interface{}, expire time.Duration, tag *string) error
	RemoveHashKey(ctx context.Context, key string) error
	HGet(ctx context.Context, key, field string, data interface{}) error
	HSet(ctx context.Context, key, field string, data interface{}, expire time.Duration) error
	DelKeysWithPattern(ctx context.Context, pattern string) error
	Del(ctx context.Context, key string) error
	// More functions will be added later on demand
}

var (
	once     sync.Once
	instance *provider
)

type Provider interface {
	RedisCache() Cache
	DumpCache() Cache
}

type provider struct {
	redis Cache
}

func GetProvider() Provider {
	once.Do(func() {
		instance = &provider{
			redis: newRedisCache(),
		}
	})

	return instance
}

func (p *provider) DumpCache() Cache {
	return &dumpCache{}
}

func (p *provider) RedisCache() Cache {
	return p.redis
}
