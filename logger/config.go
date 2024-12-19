package logger

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	EnvTest = "test"
	EnvDev  = "dev"
)

type loggerConfig struct {
	Env              string `envconfig:"ENV"                 default:"dev"`
	Enabled          bool   `envconfig:"LOG_ENABLED"         default:"true"`
	Level            int8   `envconfig:"LOG_LEVEL"           default:"0"`
	AppRole          string `envconfig:"APP_ROLE"            default:"please-give-me-a-name"`
	LogWithTimestamp bool   `envconfig:"LOG_WITH_TIMESTAMP"  default:"true"`
	EnableSampling   bool   `envconfig:"LOG_ENABLE_SAMPLING" default:"false"`
	EnableTracing    bool   `envconfig:"LOG_ENABLE_TRACING"  default:"false"`
}

func (c *loggerConfig) IsTestEnv() bool {
	return c.Env == EnvTest
}

func (c *loggerConfig) IsDevEnv() bool {
	return c.Env == EnvDev
}

func newLoggerConfig() (*loggerConfig, error) {
	cfg := &loggerConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

type RequestLogConfig struct {
	Enabled         bool     `envconfig:"REQUEST_LOG_ENABLED"       default:"false"`
	LoggingResponse bool     `envconfig:"REQUEST_LOG_WITH_RESPONSE" default:"false"`
	WhiteList       []string `envconfig:"REQUEST_LOG_WHITE_LIST"    default:"/metrics,/health-check"`
}

func NewRequestLogConfig() (*RequestLogConfig, error) {
	cfg := &RequestLogConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
