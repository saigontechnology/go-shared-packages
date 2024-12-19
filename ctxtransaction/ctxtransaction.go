package ctxtransaction

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type ContextKey string

const (
	defaultTransactionKey ContextKey = "HealthcareDbTx"
)

var (
	ErrNoTransaction             = errors.New("no database transaction found in context")
	ErrDBConnectionHasNotBeenSet = errors.New("database connection has not been set")
)

//go:generate mockery --inpackage --name TransactionContext --filename ctxtransaction_mock.go
type TransactionContext interface {
	CommitFromContext(ctx context.Context) error
	RollbackFromContext(ctx context.Context) error
	Session(ctx context.Context) *gorm.DB
	SessionWithFallback(ctx context.Context, fallbackDB *gorm.DB) *gorm.DB
	BeginWithConnection(ctx context.Context, conn *gorm.DB) (context.Context, error)
	Begin(ctx context.Context) (context.Context, error)
	ShouldHaveConnection() error
	MustHaveTransaction(
		ctx context.Context,
	) error
}

var _ TransactionContext = (*TransactionInjector)(nil)

type TransactionInjector struct {
	transactionKey ContextKey
	db             *gorm.DB
}

func NewWithConnection(db *gorm.DB) *TransactionInjector {
	return &TransactionInjector{
		transactionKey: defaultTransactionKey,
		db:             db,
	}
}

func NewWithKey(db *gorm.DB, key string) *TransactionInjector {
	return &TransactionInjector{
		transactionKey: ContextKey(key),
		db:             db,
	}
}

func (tx *TransactionInjector) BeginWithConnection(
	ctx context.Context,
	conn *gorm.DB,
) (context.Context, error) {
	transaction := conn.Begin()
	if transaction.Error != nil {
		return ctx, transaction.Error
	}
	return context.WithValue(ctx, tx.transactionKey, transaction), nil
}

func (tx *TransactionInjector) Begin(ctx context.Context) (context.Context, error) {
	return tx.BeginWithConnection(ctx, tx.db)
}

func (tx *TransactionInjector) CommitFromContext(ctx context.Context) error {
	val, ok := ctx.Value(tx.transactionKey).(*gorm.DB)
	if !ok || val == nil {
		return ErrNoTransaction
	}
	return val.Commit().Error
}

func (tx *TransactionInjector) RollbackFromContext(ctx context.Context) error {
	val, ok := ctx.Value(tx.transactionKey).(*gorm.DB)
	if !ok || val == nil {
		return ErrNoTransaction
	}
	return val.Rollback().Error
}

func (tx *TransactionInjector) Session(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(tx.transactionKey).(*gorm.DB); ok {
		return tx
	}
	return tx.db
}

func (tx *TransactionInjector) SessionWithFallback(
	ctx context.Context,
	fallbackDB *gorm.DB,
) *gorm.DB {
	if tx, ok := ctx.Value(tx.transactionKey).(*gorm.DB); ok {
		return tx
	}
	return fallbackDB
}

func (tx *TransactionInjector) MustHaveTransaction(
	ctx context.Context,
) error {
	if tx, ok := ctx.Value(defaultTransactionKey).(*gorm.DB); ok && tx != nil {
		return nil
	}
	return ErrNoTransactionFound
}

func (tx *TransactionInjector) ShouldHaveConnection() error {
	if tx.db == nil {
		return ErrDBConnectionHasNotBeenSet
	}
	return nil
}
