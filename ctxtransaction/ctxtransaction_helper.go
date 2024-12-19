package ctxtransaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/datngo2sgtech/go-packages/logger"

	"gorm.io/gorm"
)

type GetSessionFunc func(ctx context.Context) *gorm.DB

var (
	ErrNoTransactionFound = errors.New("no transaction found")
)

// used without initialization of injector struct.
func NewClosureWithFallback(
	transactionKey string,
	fallbackDB *gorm.DB,
) GetSessionFunc {
	return func(ctx context.Context) *gorm.DB {
		if tx, ok := ctx.Value(ContextKey(transactionKey)).(*gorm.DB); ok {
			return tx
		}
		return fallbackDB
	}
}

// used with default transaction key.
func SessionFromContext(
	ctx context.Context,
	fallbackDB *gorm.DB,
) *gorm.DB {
	if tx, ok := ctx.Value(defaultTransactionKey).(*gorm.DB); ok {
		return tx
	}
	return fallbackDB
}

func MustHaveTransaction(
	ctx context.Context,
) error {
	if tx, ok := ctx.Value(defaultTransactionKey).(*gorm.DB); ok && tx != nil {
		return nil
	}
	return ErrNoTransactionFound
}

func RecoverAndRollback(ctx context.Context, txCtx TransactionContext, lg logger.Logger) {
	if r := recover(); r != nil {
		lg.Error(ctx, fmt.Sprintf("[Panic Recover] %v\n%v", r, string(debug.Stack())))
	}
	if errRb := txCtx.RollbackFromContext(ctx); errRb != nil && !errors.Is(errRb, sql.ErrTxDone) {
		lg.Error(ctx, fmt.Sprintf("[ContextTransaction] Could not rollback transaction. Error: %s", errRb.Error()))
	}
}
