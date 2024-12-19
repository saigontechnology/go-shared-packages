package cache

import (
	"context"
	"log"
	"time"
)

type dumpCache struct{}

func (d *dumpCache) RemoveHashKey(ctx context.Context, key string) error {
	log.Println("[DumpCache.Get] nothing to do")
	return nil
}

func (d *dumpCache) HGet(ctx context.Context, key, field string, data interface{}) error {
	log.Println("[DumpCache.Get] nothing to do")
	return nil
}

func (d *dumpCache) HSet(
	ctx context.Context,
	key, field string,
	data interface{},
	expire time.Duration,
) error {
	log.Println("[DumpCache.Get] nothing to do")
	return nil
}

func (d *dumpCache) Remove(ctx context.Context, key string) error {
	log.Println("[DumpCache.Get] nothing to do")
	return nil
}

func (d *dumpCache) Get(ctx context.Context, key string, data interface{}, tag *string) error {
	log.Println("[DumpCache.Get] nothing to do")
	return nil
}

func (d *dumpCache) Del(ctx context.Context, key string) error {
	log.Println("[DumpCache.Del] nothing to do")
	return nil
}

func (d *dumpCache) Set(
	ctx context.Context,
	key string,
	data interface{},
	expire time.Duration,
	tag *string,
) error {
	log.Println("[DumpCache.Set] nothing to do")
	return nil
}

func (d *dumpCache) DelKeysWithPattern(ctx context.Context, pattern string) error {
	log.Println("[DumpCache.DelKeyWithPattern] nothing to do")
	return nil
}
