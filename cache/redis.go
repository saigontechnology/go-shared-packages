package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"

	"github.com/datngo2sgtech/go-packages/list"
	"github.com/datngo2sgtech/go-packages/must"
	"github.com/datngo2sgtech/go-packages/prometheus"
)

type redisCache struct {
	client    *redis.Client
	metric    prometheus.CacheMetric
	namespace string
	scanCount int64
}

func newRedisCache() *redisCache {
	cfg, err := newConfig()
	must.NotFail(err)

	var tlsConfig *tls.Config
	if cfg.TLSEnabled {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return &redisCache{
		client: redis.NewClient(&redis.Options{
			Addr:      cfg.Host,
			DB:        cfg.Database,
			Password:  cfg.Password,
			TLSConfig: tlsConfig,
		}),
		metric:    prometheus.GetCacheMetric(),
		namespace: cfg.Namespace,
		scanCount: int64(cfg.ScanCount),
	}
}

func NewMiniRedisForTest(t *testing.T) Cache {
	t.Helper()

	mr := miniredis.RunT(t)
	return &redisCache{
		namespace: "test",
		metric:    prometheus.GetCacheMetric(),
		client:    redis.NewClient(&redis.Options{Addr: mr.Addr()}),
	}
}

func (r *redisCache) newrelicRedisSegment(
	ctx context.Context,
	operation string,
) *newrelic.DatastoreSegment {
	s := &newrelic.DatastoreSegment{
		StartTime: newrelic.FromContext(ctx).StartSegmentNow(),
		Product:   newrelic.DatastoreRedis,
		Operation: operation,
	}
	return s
}

func (r *redisCache) Get(ctx context.Context, key string, data interface{}, tag *string) error {
	s := r.newrelicRedisSegment(ctx, "Get")
	defer s.End()

	txn := r.metric.NewCacheGetLatencyTransaction(tag)
	defer txn.End()
	err := r.client.Get(ctx, r.makeAppCacheKey(key)).Scan(data)
	if err == nil {
		r.metric.CountCacheHit(tag)
	} else {
		r.metric.CountCacheMiss(tag)
	}
	return err
}

func (r *redisCache) Set(
	ctx context.Context,
	key string,
	data interface{},
	expire time.Duration,
	tag *string,
) error {
	s := r.newrelicRedisSegment(ctx, "Set")
	defer s.End()
	txn := r.metric.NewCacheSetLatencyTransaction(tag)
	defer txn.End()
	err := r.client.Set(ctx, r.makeAppCacheKey(key), data, expire).Err()
	if err != nil {
		log.Printf("[Cache] could not set for key %s. Error: %s", key, err.Error())
	}
	return err
}

func (r *redisCache) HGet(ctx context.Context, key, field string, data interface{}) error {
	s := r.newrelicRedisSegment(ctx, "HGet")
	defer s.End()
	err := r.client.HGet(ctx, r.makeAppCacheKey(key), field).Scan(data)
	return err
}

func (r *redisCache) HSet(
	ctx context.Context,
	key, field string,
	data interface{},
	expire time.Duration,
) error {
	s := r.newrelicRedisSegment(ctx, "HSet")
	defer s.End()
	err := r.client.HSet(ctx, r.makeAppCacheKey(key), field, data).Err()
	if err != nil {
		log.Printf("[Cache] could not set for key %s. Error: %s", key, err.Error())
	}

	err = r.client.Expire(ctx, r.makeAppCacheKey(key), expire).Err()
	if err != nil {
		log.Printf("[Cache] could not set expire key %s. Error: %s", key, err.Error())
	}

	return err
}

func (r *redisCache) RemoveHashKey(ctx context.Context, key string) error {
	s := r.newrelicRedisSegment(ctx, "RemoveHashKey")
	defer s.End()
	fieldKeys, err := r.client.HKeys(ctx, r.makeAppCacheKey(key)).Result()
	if err != nil {
		log.Printf("[Cache] delete key %s. Error: %s", key, err.Error())
	}

	if len(fieldKeys) > 0 {
		err = r.client.HDel(ctx, r.makeAppCacheKey(key), fieldKeys...).Err()

		if err != nil {
			log.Printf("[Cache] delete key %s. Error: %s", key, err.Error())
		}
	}

	return err
}

func (r *redisCache) DelKeysWithPattern(ctx context.Context, pattern string) error {
	// add namespace to pattern
	pattern = r.makeAppCacheKey(pattern)
	cursor := uint64(0)
	for {
		var keys []string
		var err error
		keys, cursor, err = r.client.Scan(ctx, cursor, pattern, r.scanCount).Result()
		if err != nil {
			return err
		}

		keyChunks := list.Chunk(keys, 1024)
		for _, keysInChunk := range keyChunks {
			//nolint: errcheck
			r.client.Del(ctx, keysInChunk...)
		}

		if cursor == 0 {
			// No more keys
			break
		}
	}
	return nil
}

func (r *redisCache) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, r.makeAppCacheKey(key)).Err()
}

func (r *redisCache) makeAppCacheKey(key string) string {
	return fmt.Sprintf("%s_%s", r.namespace, key)
}
