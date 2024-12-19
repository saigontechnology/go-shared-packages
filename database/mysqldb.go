package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/datngo2sgtech/go-packages/must"

	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
)

var (
	mysqlDBOnce     sync.Once
	mysqlDBInstance *mysqlDB
)

type mysqlDB struct {
	mu  sync.Mutex
	dbs map[string]*gorm.DB
}

// GetMysqlDB singleton implementation makes sure only one mysqlDB is created to avoid duplicated database connection pools.
func GetMysqlDB() Connector {
	mysqlDBOnce.Do(func() {
		mysqlDBInstance = &mysqlDB{
			dbs: make(map[string]*gorm.DB),
		}
	})

	return mysqlDBInstance
}

func (m *mysqlDB) Connect(name string) *gorm.DB {
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
		mysql.New(mysql.Config{DriverName: "nrmysql", DSN: m.dsnFromConfig(cfg)}),
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

func (m *mysqlDB) SetReplicas(masterDB *gorm.DB, names []string) {
	dialectors := make([]gorm.Dialector, len(names))
	for i, name := range names {
		cfg, err := newConfig(name)
		must.NotFail(err)
		dialectors[i] = mysql.New(mysql.Config{DriverName: "nrmysql", DSN: m.dsnFromConfig(cfg)})
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

func (m *mysqlDB) dsnFromConfig(cfg *config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Name,
		cfg.Charset,
	)
}
