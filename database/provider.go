package database

import (
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/saigontechnology/go-shared-packages/env"
	"github.com/saigontechnology/go-shared-packages/must"
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
	SetConnector(c Connector) Provider
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

type provider struct {
	connector Connector
}

// GetProvider singleton implementation makes sure only one Provider is created to avoid duplicated database connection pools.
func GetProvider() Provider {
	providerOnce.Do(func() {
		providerInstance = &provider{}
	})

	return providerInstance
}

func (p *provider) SetConnector(c Connector) Provider {
	p.connector = c
	return p
}

func (p *provider) setDefaultConnector() {
	p.connector = GetMysqlDB()
}

func (p *provider) DB(name string) *gorm.DB {
	if env.IsTestEnv() {
		return GetDumpDB().Connect(name)
	}
	if p.connector == nil {
		p.setDefaultConnector()
	}
	return p.connector.Connect(name)
}

func (p *provider) SetReplicas(masterDB *gorm.DB, names []string) {
	if env.IsTestEnv() {
		GetDumpDB().SetReplicas(masterDB, names)
	}
	p.connector.SetReplicas(masterDB, names)
}
