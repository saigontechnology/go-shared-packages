package migration

import "github.com/kelseyhightower/envconfig"

type config struct {
	Host      string `default:"localhost:9091"                  envconfig:"MIGRATION_DB_HOST"`
	Name      string `default:"pet-db"                          envconfig:"MIGRATION_DB_NAME"`
	Username  string `default:"pet-db"                          envconfig:"MIGRATION_DB_USER"`
	Password  string `default:"pet-db"                          envconfig:"MIGRATION_DB_PASSWORD"`
	Charset   string `default:"utf8mb4"                         envconfig:"MIGRATION_DB_CHARSET"`
	SourceURL string `default:"file://./internal/db/migrations" envconfig:"MIGRATION_SOURCE_URL"`
}

func newConfig() (*config, error) {
	cfg := &config{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
