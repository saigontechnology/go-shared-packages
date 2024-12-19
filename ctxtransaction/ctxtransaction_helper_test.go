package ctxtransaction_test

import (
	"context"
	"errors"
	"testing"

	"github.com/datngo2sgtech/go-packages/ctxtransaction"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

//nolint:gochecknoglobals
var (
	testingKey        = "testingKey"
	testingContextKey = ctxtransaction.ContextKey(testingKey)
)

func TestNewClosureWithFallback(t *testing.T) {
	t.Parallel()
	t.Run("should return fallback db on non-existing transaction", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		closure := ctxtransaction.NewClosureWithFallback(testingKey, fallbackDB)
		assert.NotNil(t, closure)
		db := closure(context.Background())
		assert.Equal(t, db, fallbackDB)
	})
	t.Run("should return context transaction on existing key", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		transactionDB := &gorm.DB{
			Error: errors.New("test error"),
		}
		closure := ctxtransaction.NewClosureWithFallback(testingKey, fallbackDB)
		assert.NotNil(t, closure)
		ctx := context.WithValue(context.Background(), testingContextKey, transactionDB)
		db := closure(ctx)
		assert.Equal(t, transactionDB, db)
	})
	t.Run("should return fallback db on invalid transaction key", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		transactionDB := &gorm.DB{
			Error: errors.New("test error"),
		}
		closure := ctxtransaction.NewClosureWithFallback("testingKey1", fallbackDB)
		assert.NotNil(t, closure)
		ctx := context.WithValue(
			context.Background(),
			ctxtransaction.ContextKey("testingKey2"),
			transactionDB,
		)
		db := closure(ctx)
		assert.Equal(t, fallbackDB, db)
	})
}

func TestSessionFromContext(t *testing.T) {
	t.Parallel()
	t.Run("should return fallback db on non-existing transaction", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{
			Error: errors.New("test error"),
		}
		db := ctxtransaction.SessionFromContext(context.Background(), fallbackDB)
		assert.Equal(t, fallbackDB, db)
	})
	t.Run("should return context transaction on existing key", func(t *testing.T) {
		t.Parallel()
		fallbackDB := &gorm.DB{}
		transactionDB := &gorm.DB{
			Error: errors.New("test error"),
		}
		ctx := context.WithValue(
			context.Background(),
			ctxtransaction.ContextKey("HealthcareDbTx"),
			transactionDB,
		)
		db := ctxtransaction.SessionFromContext(ctx, fallbackDB)
		assert.Equal(t, transactionDB, db)
	})
}
