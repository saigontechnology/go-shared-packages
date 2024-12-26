package authmiddleware

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	IdentityServer string `default:"identity:8081" envconfig:"IDENTITY_SERVER"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
