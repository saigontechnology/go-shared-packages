package test

import "github.com/kelseyhightower/envconfig"

type itConfig struct {
	SqliteTestDatabaseDir             string `default:"/tmp/pet/sqlite"   envconfig:"SQLITE_TEST_DATABASE_DIR"`
	SqliteTestDatabaseInitDir         string `default:"../../test/sqlite" envconfig:"SQLITE_TEST_DATABASE_INIT_DIR"`
	SqliteTestDatabaseSchemaFile      string `default:"schema.sql"        envconfig:"SQLITE_TEST_DATABASE_SCHEMA_FILE"`
	SqliteTestDatabaseInitialDataFile string `default:"initial_data.sql"  envconfig:"SQLITE_TEST_DATABASE_INITIAL_DATA_FILE"`
}

func newItConfig() (*itConfig, error) {
	cfg := &itConfig{}
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
