package test

import "github.com/kelseyhightower/envconfig"

type htConfig struct {
	TranslationsDir string `default:"../../translations" envconfig:"HT_TRANSLATIONS_DIR"`
}

func newHtConfig() (*htConfig, error) {
	cfg := &htConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
