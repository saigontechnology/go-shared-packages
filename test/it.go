package test

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/datngo2sgtech/go-packages/cache"
	"github.com/datngo2sgtech/go-packages/fixture"
	"github.com/datngo2sgtech/go-packages/must"
)

const (
	TypeItOptionDatabaseInitDir = "ItOptionDatabaseInitDir"
)

// IT is a utility to implement integration test.
type IT interface {
	UT
	DB() *gorm.DB
	FixtureStore() fixture.Store
	Cache() cache.Cache
}

// it is an Integration Test utility with Sqlite database.
type it struct {
	*ut
	cfg          *itConfig
	dbFile       string
	db           *gorm.DB
	fixtureStore fixture.Store
	cache        cache.Cache
}

type ItOption interface {
	Type() string
	Value() string
}

type ItOptionDatabaseInitDir struct {
	InitDir string
}

func (o *ItOptionDatabaseInitDir) Type() string {
	return TypeItOptionDatabaseInitDir
}

func (o *ItOptionDatabaseInitDir) Value() string {
	return o.InitDir
}

func NewIT(t *testing.T, options ...ItOption) IT {
	t.Helper()

	return newIT(t, options...)
}

func newIT(t *testing.T, options ...ItOption) *it {
	t.Helper()

	it := &it{
		ut: newUT(t),
	}
	cfg, err := newItConfig()
	must.NotFail(err)
	parseItOptions(cfg, options...)
	it.cfg = cfg
	it.cache = cache.NewMiniRedisForTest(t)
	it.fixtureStore = fixture.NewStore()
	it.createSqliteDB()
	it.initDatabase()
	// Register a cleanup function to close a database connection and delete the database file when a test finished
	t.Cleanup(it.cleanUp)
	return it
}

func parseItOptions(cfg *itConfig, options ...ItOption) {
	for _, option := range options {
		if option.Type() == TypeItOptionDatabaseInitDir {
			cfg.SqliteTestDatabaseInitDir = option.Value()
		}
	}
}

func (i *it) DB() *gorm.DB {
	return i.db
}

func (i *it) FixtureStore() fixture.Store {
	return i.fixtureStore
}

func (i *it) Cache() cache.Cache {
	return i.cache
}

func (i *it) createSqliteDB() {
	i.generateDBFile()
	err := os.MkdirAll(filepath.Dir(i.dbFile), 0o770)
	must.NotFail(err)
	f, err := os.Create(i.dbFile)
	must.NotFail(err)
	defer func() {
		if err := f.Close(); err != nil {
			log.Println("!!! Could not close test Sqlite database file")
		}
	}()
	db, err := gorm.Open(sqlite.Open(i.dbFile), &gorm.Config{})
	must.NotFail(err)
	i.db = db
}

func (i *it) generateDBFile() {
	// runFolder is used to avoid conflict when multi runs of a test are triggered at the same time
	runFolder := make([]rune, 5)
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	for i := 0; i < 5; i++ {
		runFolder[i] = letterRunes[rand.Intn(len(letterRunes))] //nolint:gosec
	}
	testName := strings.ToLower(i.t.Name())
	testName = strings.ReplaceAll(testName, " ", "_")
	i.dbFile = fmt.Sprintf(
		"%s/%s_%s.db",
		i.cfg.SqliteTestDatabaseDir,
		string(runFolder),
		testName,
	)
}

func (i *it) initDatabase() {
	schemaFilePath := filepath.Join(
		i.cfg.SqliteTestDatabaseInitDir,
		i.cfg.SqliteTestDatabaseSchemaFile,
	)
	schemaBytes, err := os.ReadFile(schemaFilePath)
	must.NotFail(err)
	if len(schemaBytes) > 0 {
		schemaTx := i.db.Exec(string(schemaBytes))
		must.NotFail(schemaTx.Error)
	}
	initFilePath := filepath.Join(
		i.cfg.SqliteTestDatabaseInitDir,
		i.cfg.SqliteTestDatabaseInitialDataFile,
	)
	initialDataBytes, err := os.ReadFile(initFilePath)
	must.NotFail(err)
	if len(initialDataBytes) > 0 {
		initTx := i.db.Exec(string(initialDataBytes))
		must.NotFail(initTx.Error)
	}
}

// Close a database connection pool and delete the database file when a test finished.
func (i *it) cleanUp() {
	sqlDB, err := i.db.DB()
	must.NotFail(err)
	err = sqlDB.Close()
	must.NotFail(err)
	err = os.Remove(i.dbFile)
	must.NotFail(err)
}
