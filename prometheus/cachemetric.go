package prometheus

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/datngo2sgtech/go-packages/must"
)

const (
	vectorCacheHit  = "hit"
	vectorCacheMiss = "miss"
	vectorCacheGet  = "get"
	vectorCacheSet  = "set"

	vectorTag    = "tag"
	vectorStatus = "status"
	vectorOp     = "op"

	nameCacheHitTotal        = "cache_hit_total"
	descriptionCacheHit      = "Monitor cache hit by tag"
	nameCacheDurationSeconds = "cache_duration_seconds"
	descriptionCacheDuration = "Monitor cache latency by tag"
)

var (
	cacheMetricOnce     sync.Once
	cacheMetricInstance *cacheMetric
)

type CacheMetric interface {
	CountCacheHit(tag *string)
	CountCacheMiss(tag *string)
	NewCacheGetLatencyTransaction(tag *string) CacheLatencyMetricTxn
	NewCacheSetLatencyTransaction(tag *string) CacheLatencyMetricTxn
}

type cacheMetric struct {
	cfg           *cacheMetricConfig
	cacheLatency  *prometheus.HistogramVec
	cacheHitTotal *prometheus.CounterVec
}

func GetCacheMetric() CacheMetric {
	cacheMetricOnce.Do(func() {
		cfg, err := newCacheMetricConfig()
		must.NotFail(err)
		cacheHitTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.Metric.Namespace,
			Name:      fmt.Sprintf("%s_%s", cfg.Metric.MetricPrefix, nameCacheHitTotal),
			Help:      descriptionCacheHit,
		}, []string{vectorTag, vectorStatus})
		prometheus.MustRegister(cacheHitTotal)
		cacheLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Metric.Namespace,
			Name:      fmt.Sprintf("%s_%s", cfg.Metric.MetricPrefix, nameCacheDurationSeconds),
			Help:      descriptionCacheDuration,
			Buckets: prometheus.LinearBuckets(
				cfg.CacheLatencyBucketStart,
				cfg.CacheLatencyBucketWidth,
				cfg.CacheLatencyBucketCount,
			),
		}, []string{vectorTag, vectorOp})
		prometheus.MustRegister(cacheLatency)
		cacheMetricInstance = &cacheMetric{
			cfg:           cfg,
			cacheHitTotal: cacheHitTotal,
			cacheLatency:  cacheLatency,
		}
	})

	return cacheMetricInstance
}

func (m *cacheMetric) CountCacheHit(tag *string) {
	m.countCacheOp(tag, vectorCacheHit)
}

func (m *cacheMetric) CountCacheMiss(tag *string) {
	m.countCacheOp(tag, vectorCacheMiss)
}

func (m *cacheMetric) countCacheOp(tag *string, status string) {
	if !m.cfg.CacheMetricEnabled {
		return
	}
	if tag == nil {
		return
	}
	m.cacheHitTotal.WithLabelValues(*tag, status).Inc()
}

func (m *cacheMetric) NewCacheGetLatencyTransaction(tag *string) CacheLatencyMetricTxn {
	if !m.cfg.CacheMetricEnabled || tag == nil {
		return &nullCacheLatencyTxn{}
	}

	return &cacheLatencyMetricTxn{
		op:            vectorCacheGet,
		tag:           *tag,
		start:         time.Now(),
		latencyMetric: m.cacheLatency,
	}
}

func (m *cacheMetric) NewCacheSetLatencyTransaction(tag *string) CacheLatencyMetricTxn {
	if !m.cfg.CacheMetricEnabled || tag == nil {
		return &nullCacheLatencyTxn{}
	}

	return &cacheLatencyMetricTxn{
		op:            vectorCacheSet,
		tag:           *tag,
		start:         time.Now(),
		latencyMetric: m.cacheLatency,
	}
}

type CacheLatencyMetricTxn interface {
	End()
}

type nullCacheLatencyTxn struct{}

func (n *nullCacheLatencyTxn) End() {}

type cacheLatencyMetricTxn struct {
	op            string
	tag           string
	start         time.Time
	latencyMetric *prometheus.HistogramVec
}

func (tx *cacheLatencyMetricTxn) End() {
	tx.latencyMetric.WithLabelValues(tx.tag, tx.op).Observe(time.Since(tx.start).Seconds())
}
