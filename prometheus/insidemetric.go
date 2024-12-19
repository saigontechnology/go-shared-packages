package prometheus

import (
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/datngo2sgtech/go-packages/must"
)

const (
	vectorGroup                     = "group"
	vectorSegment                   = "segment"
	vectorGroupRepository           = "repository"
	nameInsideDurationSeconds       = "inside_duration_seconds"
	descriptionInsideDurationMetric = "Monitor inside latency"
)

var (
	insideMetricOnce     sync.Once
	insideMetricInstance *insideLatencyMetric
)

type InsideLatencyMetric interface {
	StartRepositoryTransaction(metricName string) InsideLatencyMetricTxn
}

type insideLatencyMetric struct {
	cfg           *insideMetricConfig
	latencyMetric *prometheus.HistogramVec
}

func GetInsideLatencyMetric() InsideLatencyMetric {
	insideMetricOnce.Do(func() {
		cfg, err := newInsideMetricConfig()
		must.NotFail(err)
		segmentLatency := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Metric.Namespace,
				Name:      fmt.Sprintf("%s_%s", cfg.Metric.MetricPrefix, nameInsideDurationSeconds),
				Help:      descriptionInsideDurationMetric,
				Buckets: prometheus.LinearBuckets(
					cfg.InsideLatencyBucketStart,
					cfg.InsideLatencyBucketWidth,
					cfg.InsideLatencyBucketCount,
				),
			},
			[]string{vectorGroup, vectorSegment},
		)
		prometheus.MustRegister(segmentLatency)
		insideMetricInstance = &insideLatencyMetric{
			cfg:           cfg,
			latencyMetric: segmentLatency,
		}
	})

	return insideMetricInstance
}

func (m *insideLatencyMetric) StartRepositoryTransaction(name string) InsideLatencyMetricTxn {
	if !m.cfg.InsideMetricEnabled {
		return &nullInsideLatencyMetricTxn{}
	}

	return &insideLatencyMetricTxn{
		group:         vectorGroupRepository,
		name:          name,
		start:         time.Now(),
		latencyMetric: m.latencyMetric,
	}
}

type InsideLatencyMetricTxn interface {
	End()
}

type nullInsideLatencyMetricTxn struct{}

func (tx *nullInsideLatencyMetricTxn) End() {}

type insideLatencyMetricTxn struct {
	group         string
	name          string
	start         time.Time
	latencyMetric *prometheus.HistogramVec
}

func (tx *insideLatencyMetricTxn) End() {
	if tx.name == "" {
		return
	}
	tx.latencyMetric.WithLabelValues(tx.group, tx.name).Observe(time.Since(tx.start).Seconds())
}
