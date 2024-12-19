package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/pkg/errors"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/datngo2sgtech/go-packages/must"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type MysqlMigration interface {
	Migrate(migration string, steps int) error
}

type mysqlMigration struct {
	m *migrate.Migrate
}

func NewMysqlMigration(tablePrefix string) MysqlMigration {
	cfg, err := newConfig()
	must.NotFail(err)
	// multiStatements=true is very important to execute a migration file with multi statements with an existing database connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local&multiStatements=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Name,
		cfg.Charset,
	)
	// Gorm can connect to database with password that contains reserved characters for URL.
	// Ref:
	//    https://github.com/golang-migrate/migrate
	//    https://en.wikipedia.org/wiki/Percent-encoding#Percent-encoding_reserved_characters
	gormDB, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
	})
	must.NotFail(err)
	// Migrate use standard sqlDB
	sqlDB, err := gormDB.DB()
	must.NotFail(err)
	driver, err := mysql.WithInstance(sqlDB, &mysql.Config{})
	must.NotFail(err)
	m, err := migrate.NewWithDatabaseInstance(cfg.SourceURL, cfg.Name, driver)
	must.NotFail(err)
	return &mysqlMigration{
		m: m,
	}
}

func (mm *mysqlMigration) Migrate(migration string, steps int) error {
	var migrateErr error
	switch migration {
	case "up":
		migrateErr = mm.m.Up()
		//nolint: gosimple
		break
	case "steps":
		if steps == 0 {
			return errors.New(
				"[Database migration] Steps must not be 0. Please use a positive number to migrate up, a negative number to migrate down.",
			)
		}
		migrateErr = mm.m.Steps(steps)
		//nolint: gosimple
		break
	default:
		return errors.New(
			"[Database migration] Migration value is not valid, please provide one of these values (up, steps)",
		)
	}
	if migrateErr != nil && !errors.Is(migrateErr, migrate.ErrNoChange) {
		return fmt.Errorf("[Database migration] Error: %w", migrateErr)
	}

	return nil
}
