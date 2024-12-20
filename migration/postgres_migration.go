package migration

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/saigontechnology/go-shared-packages/must"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresMigration interface {
	Migrate(migration string, steps int) error
}

type postgresMigration struct {
	m *migrate.Migrate
}

func NewPostgresMigration(tablePrefix string) PostgresMigration {
	cfg, err := newConfig()
	must.NotFail(err)
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
	gormDB, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
	})
	must.NotFail(err)
	sqlDB, err := gormDB.DB()
	must.NotFail(err)
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	must.NotFail(err)
	m, err := migrate.NewWithDatabaseInstance(cfg.SourceURL, cfg.Name, driver)
	must.NotFail(err)
	return &postgresMigration{
		m: m,
	}
}

func (pm *postgresMigration) Migrate(migration string, steps int) error {
	var migrateErr error
	switch migration {
	case "up":
		migrateErr = pm.m.Up()
		break
	case "steps":
		if steps == 0 {
			return errors.New(
				"[Database migration] Steps must not be 0. Please use a positive number to migrate up, a negative number to migrate down.",
			)
		}
		migrateErr = pm.m.Steps(steps)
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
