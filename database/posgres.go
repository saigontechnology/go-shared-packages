package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/saigontechnology/go-shared-packages/must"
)

var (
	postgresDBOnce     sync.Once
	postgresDBInstance *postgresDB
)

type postgresDB struct {
	mu  sync.Mutex
	dbs map[string]*gorm.DB
}

// GetPostgresDB singleton implementation makes sure only one postgresDB is created to avoid duplicated database connection pools.
func GetPostgresDB() Connector {
	postgresDBOnce.Do(func() {
		postgresDBInstance = &postgresDB{
			dbs: make(map[string]*gorm.DB),
		}
	})

	return postgresDBInstance
}

func (m *postgresDB) Connect(name string) *gorm.DB {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.dbs[name]
	if ok {
		return m.dbs[name]
	}

	cfg, err := newConfig(name)
	must.NotFail(err)
	// Setting up gorm config
	gormConfig := gorm.Config{
		// We should monitor service performance first then decide whether we disable default transaction or not
		// SkipDefaultTransaction: true,
	}
	if !cfg.ErrorLog {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{DriverName: "pgx", DSN: m.dsnFromConfig(cfg)}),
		&gormConfig,
	)
	must.NotFail(err)
	sqlDB, err := gormDB.DB()
	must.NotFail(err)
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))

	m.dbs[name] = gormDB
	return m.dbs[name]
}

func (m *postgresDB) SetReplicas(masterDB *gorm.DB, names []string) {
	dialectors := make([]gorm.Dialector, len(names))
	for i, name := range names {
		cfg, err := newConfig(name)
		must.NotFail(err)
		dialectors[i] = postgres.New(postgres.Config{DriverName: "pgx", DSN: m.dsnFromConfig(cfg)})
	}

	err := masterDB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: dialectors,
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))
	must.NotFail(err)
}

func (m *postgresDB) dsnFromConfig(cfg *config) string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Ho_Chi_Minh",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Name,
		cfg.Port,
	)
}
