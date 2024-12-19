package prometheus

import "github.com/kelseyhightower/envconfig"

type metricConfig struct {
	Namespace    string `default:""         envconfig:"PROMETHEUS_NAMESPACE"`
	MetricPrefix string `default:"your_api" envconfig:"PROMETHEUS_METRIC_PREFIX"`
}

type handlerMetricConfig struct {
	Metric                    *metricConfig
	HandlerMetricEnabled      bool    `default:"true" envconfig:"PROMETHEUS_HANDLER_METRIC_ENABLED"`
	HandlerLatencyBucketStart float64 `default:"0.05" envconfig:"PROMETHEUS_HANDLER_LATENCY_BUCKET_START"`
	HandlerLatencyBucketWidth float64 `default:"0.05" envconfig:"PROMETHEUS_HANDLER_LATENCY_BUCKET_WIDTH"`
	HandlerLatencyBucketCount int     `default:"5"    envconfig:"PROMETHEUS_HANDLER_LATENCY_BUCKET_COUNT"`
}

type cacheMetricConfig struct {
	Metric                  *metricConfig
	CacheMetricEnabled      bool    `default:"false" envconfig:"PROMETHEUS_CACHE_METRIC_ENABLED"`
	CacheLatencyBucketStart float64 `default:"0.01"  envconfig:"PROMETHEUS_CACHE_LATENCY_BUCKET_START"`
	CacheLatencyBucketWidth float64 `default:"0.01"  envconfig:"PROMETHEUS_CACHE_LATENCY_BUCKET_WIDTH"`
	CacheLatencyBucketCount int     `default:"5"     envconfig:"PROMETHEUS_CACHE_LATENCY_BUCKET_COUNT"`
}

type insideMetricConfig struct {
	Metric                   *metricConfig
	InsideMetricEnabled      bool    `default:"true" envconfig:"PROMETHEUS_INSIDE_METRIC_ENABLED"`
	InsideLatencyBucketStart float64 `default:"0.05" envconfig:"PROMETHEUS_INSIDE_LATENCY_BUCKET_START"`
	InsideLatencyBucketWidth float64 `default:"0.05" envconfig:"PROMETHEUS_INSIDE_LATENCY_BUCKET_WIDTH"`
	InsideLatencyBucketCount int     `default:"3"    envconfig:"PROMETHEUS_INSIDE_LATENCY_BUCKET_COUNT"`
}

func newHandlerMetricConfig() (*handlerMetricConfig, error) {
	metricCfg := &metricConfig{}
	if err := envconfig.Process("", metricCfg); err != nil {
		return nil, err
	}
	cfg := &handlerMetricConfig{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	cfg.Metric = metricCfg
	return cfg, nil
}

func newCacheMetricConfig() (*cacheMetricConfig, error) {
	metricCfg := &metricConfig{}
	if err := envconfig.Process("", metricCfg); err != nil {
		return nil, err
	}
	cfg := &cacheMetricConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	cfg.Metric = metricCfg
	return cfg, nil
}

func newInsideMetricConfig() (*insideMetricConfig, error) {
	metricCfg := &metricConfig{}
	if err := envconfig.Process("", metricCfg); err != nil {
		return nil, err
	}
	cfg := &insideMetricConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}
	cfg.Metric = metricCfg
	return cfg, nil
}
