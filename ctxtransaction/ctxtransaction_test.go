package ctxtransaction_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/datngo2sgtech/go-packages/ctxtransaction"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//nolint:gochecknoglobals
var (
	testKey        = "testKey"
	testVal        = "testVal"
	testContextKey = ctxtransaction.ContextKey(testKey)
	errTest        = errors.New("test error")
)

func TestTransactionContext_Session(t *testing.T) {
	t.Parallel()
	t.Run("should return origin transaction injected to context", func(t *testing.T) {
		t.Parallel()
		originDB := &gorm.DB{
			Error: errTest,
		}
		ctxTransaction := ctxtransaction.NewWithConnection(nil)
		ctx := context.WithValue(
			context.Background(),
			ctxtransaction.ContextKey("HealthcareDbTx"),
			originDB,
		)
		db := ctxTransaction.Session(ctx)
		assert.NotNil(t, db)
		assert.Equal(t, originDB, db)
	})
	t.Run(
		"should return origin transaction injected to context with custom key",
		func(t *testing.T) {
			t.Parallel()
			originDB := &gorm.DB{
				Error: errTest,
			}
			ctxTransaction := ctxtransaction.NewWithKey(nil, testKey)
			ctx := context.WithValue(context.Background(), testContextKey, originDB)
			db := ctxTransaction.Session(ctx)
			assert.NotNil(t, db)
			assert.Equal(t, originDB, db)
		},
	)
}

func TestTransactionContext_SessionWithFallback(t *testing.T) {
	t.Parallel()
	t.Run("should return fallback db on no transaction found in context", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		ctxTransaction := ctxtransaction.NewWithConnection(nil)
		db := ctxTransaction.SessionWithFallback(context.Background(), fallbackDB)
		assert.NotNil(t, db)
		assert.Equal(t, fallbackDB, db)
	})
	t.Run("should return fallback db invalid key in context", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		originDB := &gorm.DB{
			Error: errTest,
		}
		ctxTransaction := ctxtransaction.NewWithConnection(originDB)
		ctx := context.WithValue(
			context.Background(),
			ctxtransaction.ContextKey("invalidKey"),
			originDB,
		)
		db := ctxTransaction.SessionWithFallback(ctx, fallbackDB)
		assert.NotNil(t, db)
		assert.Equal(t, fallbackDB, db)
	})
	t.Run("should return transaction db on transaction found in context", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		originDB := &gorm.DB{
			Error: errTest,
		}
		ctxTransaction := ctxtransaction.NewWithConnection(originDB)
		ctx := context.WithValue(
			context.Background(),
			ctxtransaction.ContextKey("HealthcareDbTx"),
			originDB,
		)
		db := ctxTransaction.SessionWithFallback(ctx, fallbackDB)
		assert.NotNil(t, db)
		assert.Equal(t, originDB, db)
	})
}

func TestTransactionContext_RollbackFromContext(t *testing.T) {
	t.Parallel()
	t.Run("should rollback transaction successfully", func(t *testing.T) {
		t.Parallel()
		conn, mock, gormDB := newMockDB(t)
		defer conn.Close()
		mock.ExpectBegin()
		mock.ExpectRollback()

		ctx := context.WithValue(context.Background(), testContextKey, "testVal")
		ctxTransaction := ctxtransaction.NewWithConnection(gormDB)
		beginCtx, err := ctxTransaction.Begin(ctx)
		require.NoError(t, err)
		err = ctxTransaction.RollbackFromContext(beginCtx)
		require.NoError(t, err)
		// assert if the transaction is the same
		val, ok := beginCtx.Value(testContextKey).(string)
		assert.True(t, ok)
		assert.Equal(t, testVal, val)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTransactionContext_CommitFromContext(t *testing.T) {
	t.Parallel()
	t.Run("should commit transaction successfully", func(t *testing.T) {
		t.Parallel()
		conn, mock, gormDB := newMockDB(t)
		defer conn.Close()
		mock.ExpectBegin()
		mock.ExpectCommit()

		ctx := context.WithValue(context.Background(), testContextKey, testVal)
		ctxTransaction := ctxtransaction.NewWithConnection(gormDB)
		beginCtx, err := ctxTransaction.Begin(ctx)
		require.NoError(t, err)
		err = ctxTransaction.CommitFromContext(beginCtx)
		require.NoError(t, err)
		// assert if the transaction is the same
		val, ok := beginCtx.Value(testContextKey).(string)
		assert.True(t, ok)
		assert.Equal(t, testVal, val)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTransactionContext_BeginWithConnection(t *testing.T) {
	t.Parallel()
	t.Run("should begin transaction successfully", func(t *testing.T) {
		t.Parallel()
		conn, mock, gormDB := newMockDB(t)
		defer conn.Close()
		mock.ExpectBegin()

		ctx := context.WithValue(context.Background(), testContextKey, testVal)
		ctxTransaction := ctxtransaction.NewWithConnection(nil)
		beginCtx, err := ctxTransaction.BeginWithConnection(ctx, gormDB)
		require.NoError(t, err)
		// assert if the transaction is the same
		val, ok := beginCtx.Value(testContextKey).(string)
		assert.True(t, ok)
		assert.Equal(t, testVal, val)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTransactionContext_Begin(t *testing.T) {
	t.Parallel()
	t.Run("should begin transaction successfully", func(t *testing.T) {
		t.Parallel()
		conn, mock, gormDB := newMockDB(t)
		defer conn.Close()
		mock.ExpectBegin()

		ctx := context.WithValue(context.Background(), testContextKey, testVal)
		ctxTransaction := ctxtransaction.NewWithConnection(gormDB)
		beginCtx, err := ctxTransaction.Begin(ctx)
		require.NoError(t, err)
		// assert if the transaction is the same
		val, ok := beginCtx.Value(testContextKey).(string)
		assert.True(t, ok)
		assert.Equal(t, testVal, val)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTransactionContext_MustHaveTransaction(t *testing.T) {
	t.Parallel()
	t.Run("should return nil error on exist transaction", func(t *testing.T) {
		t.Parallel()
		db := &gorm.DB{
			Error: errors.New("test error"),
		}
		ctx := context.WithValue(context.Background(), ctxtransaction.ContextKey("HealthcareDbTx"), db)
		ctxTransaction := ctxtransaction.NewWithConnection(db)
		err := ctxTransaction.MustHaveTransaction(ctx)
		require.NoError(t, err)
	})
	t.Run("should return error on non-exist transaction", func(t *testing.T) {
		t.Parallel()
		ctx := context.WithValue(context.Background(), testContextKey, testVal)
		ctxTransaction := ctxtransaction.NewWithConnection(nil)
		err := ctxTransaction.MustHaveTransaction(ctx)
		require.ErrorIs(t, ctxtransaction.ErrNoTransactionFound, err)
	})
}

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, *gorm.DB) {
	t.Helper()
	conn, mock, err := sqlmock.New()
	require.NoError(t, err)
	assert.NotNil(t, conn)
	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      conn,
		SkipInitializeWithVersion: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	require.NoError(t, err)
	assert.NotNil(t, conn)
	return conn, mock, gormDB
}
