package prometheus

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/datngo2sgtech/go-packages/must"
)

var (
	apiMetricOnce     sync.Once
	apiMetricInstance *apiMetric
)

const (
	vectorAPI                        = "api"
	vectorCode                       = "code"
	nameHTTPAPIDurationSeconds       = "http_api_duration_seconds"
	descriptionHTTPAPIDuration       = "Monitor API latency"
	nameHTTPResponseCodeTotal        = "http_api_response_code_total"
	descriptionHTTPResponseCodeTotal = "Monitor response code of APIs"
)

type APIMetric interface {
	IsEnabled() bool
	CountResponseCode(name string, status int)
	NewLatencyMetricTransaction(name string) APILatencyMetricTxn
}

type apiMetric struct {
	cfg               *handlerMetricConfig
	handlerLatency    *prometheus.HistogramVec
	responseCodeTotal *prometheus.CounterVec
}

func NewAPIMetric() APIMetric {
	apiMetricOnce.Do(func() {
		cfg, err := newHandlerMetricConfig()
		must.NotFail(err)
		handlerLatency := prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: cfg.Metric.Namespace,
				Name: fmt.Sprintf(
					"%s_%s",
					cfg.Metric.MetricPrefix,
					nameHTTPAPIDurationSeconds,
				),
				Help: descriptionHTTPAPIDuration,
				Buckets: prometheus.LinearBuckets(
					cfg.HandlerLatencyBucketStart,
					cfg.HandlerLatencyBucketWidth,
					cfg.HandlerLatencyBucketCount,
				),
			},
			[]string{vectorAPI},
		)
		prometheus.MustRegister(handlerLatency)
		responseCodeTotal := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: cfg.Metric.Namespace,
				Name:      fmt.Sprintf("%s_%s", cfg.Metric.MetricPrefix, nameHTTPResponseCodeTotal),
				Help:      descriptionHTTPResponseCodeTotal,
			},
			[]string{vectorAPI, vectorCode},
		)
		prometheus.MustRegister(responseCodeTotal)
		apiMetricInstance = &apiMetric{
			cfg:               cfg,
			handlerLatency:    handlerLatency,
			responseCodeTotal: responseCodeTotal,
		}
	})

	return apiMetricInstance
}

func (m *apiMetric) IsEnabled() bool {
	return m.cfg.HandlerMetricEnabled
}

func (m *apiMetric) CountResponseCode(name string, status int) {
	if !m.cfg.HandlerMetricEnabled {
		return
	}
	if name == "" {
		return
	}
	m.responseCodeTotal.WithLabelValues(name, strconv.Itoa(status)).Inc()
}

func (m *apiMetric) NewLatencyMetricTransaction(name string) APILatencyMetricTxn {
	if !m.cfg.HandlerMetricEnabled {
		return &nullAPILatencyMetricTxn{}
	}
	return &apiLatencyMetricTxn{
		name:          name,
		start:         time.Now(),
		latencyMetric: m.handlerLatency,
	}
}

type APILatencyMetricTxn interface {
	End()
}

type nullAPILatencyMetricTxn struct{}

func (tx *nullAPILatencyMetricTxn) End() {}

type apiLatencyMetricTxn struct {
	name          string
	start         time.Time
	latencyMetric *prometheus.HistogramVec
}

func (tx *apiLatencyMetricTxn) End() {
	if tx.name == "" {
		return
	}
	tx.latencyMetric.WithLabelValues(tx.name).Observe(time.Since(tx.start).Seconds())
}

func GinAPIMetricMiddleware(metric APIMetric) gin.HandlerFunc {
	return func(c *gin.Context) {
		if metric.IsEnabled() {
			// Need an improvement to have a better name
			name := fmt.Sprintf("%s %s", c.Request.Method, getHandlerName(c, true))
			txn := metric.NewLatencyMetricTransaction(name)
			defer func() {
				txn.End()
				metric.CountResponseCode(name, c.Writer.Status())
			}()
		}
		c.Next()
	}
}

type handlerNamer interface {
	HandlerName() string
}

func getHandlerName(c handlerNamer, useNewNames bool) string {
	if useNewNames {
		if fp, ok := c.(interface{ FullPath() string }); ok {
			return fp.FullPath()
		}
	}
	return c.HandlerName()
}
