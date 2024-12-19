package database

import "gorm.io/gorm"

type dumpDB struct {
	db *gorm.DB
}

// GetDumpDB is used in test.
func GetDumpDB() Connector {
	return &dumpDB{}
}

func (d *dumpDB) Connect(name string) *gorm.DB {
	return d.db
}

func (d *dumpDB) SetReplicas(masterDB *gorm.DB, names []string) {}
