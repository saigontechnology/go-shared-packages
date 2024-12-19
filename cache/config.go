package cache

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host       string `default:"your-service-redis:6379" envconfig:"CACHE_HOST"`
	Database   int    `default:"0"                       envconfig:"CACHE_DATABASE"`
	Password   string `default:""                        envconfig:"CACHE_PASSWORD"`
	Namespace  string `default:"your_service"            envconfig:"CACHE_NAMESPACE"`
	TLSEnabled bool   `default:"false"                   envconfig:"CACHE_TLS_ENABLED"`
	ScanCount  int    `default:"5000"                    envconfig:"CACHE_SCAN_COUNT"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
