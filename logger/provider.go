package logger

import (
	"sync"
)

var (
	providerOnce     sync.Once
	providerInstance *provider
)

type Provider interface {
	Logger() Logger
}

type provider struct {
	l Logger
}

// GetProvider singleton implementation makes sure only one Provider is created to avoid duplicated logger
func GetProvider() Provider {
	providerOnce.Do(func() {
		cfg, err := newLoggerConfig()
		if err != nil {
			panic(err)
		}
		var l Logger
		if !cfg.Enabled || cfg.IsTestEnv() {
			l = NewNoopLogger()
		} else {
			l, err = NewZapLogger(cfg)
			if err != nil {
				panic(err)
			}
		}

		providerInstance = &provider{
			l: l,
		}
	})

	return providerInstance
}

func (p *provider) Logger() Logger {
	return p.l
}
