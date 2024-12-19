package database

import (
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"your-service-db:3306" envconfig:"DB_HOST"`
	Name     string `default:"your-db"              envconfig:"DB_NAME"`
	Username string `default:"your-db"              envconfig:"DB_USER"`
	Password string `default:"your-db"              envconfig:"DB_PASSWORD"`
	Charset  string `default:"utf8mb4"              envconfig:"DB_CHARSET"`
	ErrorLog bool   `default:"false"                envconfig:"DB_ENABLE_ERR_LOG"`
	// Reference https://www.alexedwards.net/blog/configuring-sqldb
	// Default MaxIdleConns in sql.DB is 2
	MaxIdleConns int `default:"2" envconfig:"DB_MAX_IDLE_CONNS"`
	// Default MaxOpenConns in sql.DB is 0 (unlimited)
	MaxOpenConns int `default:"0" envconfig:"DB_MAX_OPEN_CONNS"`
	// Default is 0, connections are not closed due to a connection's age.
	ConnMaxLifetime int64 `default:"0" envconfig:"DB_CONN_MAX_LIFETIME"`
}

func newConfig(name string) (*config, error) {
	cfg := &config{}
	err := envconfig.Process(name, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
