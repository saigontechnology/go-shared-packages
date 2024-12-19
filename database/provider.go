package database

import (
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/datngo2sgtech/go-packages/env"
	"github.com/datngo2sgtech/go-packages/must"
)

var (
	providerOnce     sync.Once
	providerInstance *provider
)

type Connector interface {
	Connect(name string) *gorm.DB
	SetReplicas(masterDB *gorm.DB, names []string)
}

type Provider interface {
	DB(name string) *gorm.DB
	SetReplicas(masterDB *gorm.DB, names []string)
}

// CloseDB closes database connection pool before exiting the main function.
func CloseDB(db *gorm.DB) {
	sqlDB, err := db.DB()
	must.NotFail(err)
	if err = sqlDB.Close(); err != nil {
		log.Println("!!! API could not close database connection pool")
	}
}

type provider struct{}

// GetProvider singleton implementation makes sure only one Provider is created to avoid duplicated database connection pools.
func GetProvider() Provider {
	providerOnce.Do(func() {
		providerInstance = &provider{}
	})

	return providerInstance
}

func (p *provider) DB(name string) *gorm.DB {
	if env.IsTestEnv() {
		return GetDumpDB().Connect(name)
	}
	return GetMysqlDB().Connect(name)
}

func (p *provider) SetReplicas(masterDB *gorm.DB, names []string) {
	if env.IsTestEnv() {
		GetDumpDB().SetReplicas(masterDB, names)
	}
	GetMysqlDB().SetReplicas(masterDB, names)
}
