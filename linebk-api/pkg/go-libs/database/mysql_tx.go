package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// TxManager defines the interface for managing MySQL database transactions with sqlx.
type TxManager interface {
	ExecInTx(
		ctx context.Context,
		fn func(ctx context.Context, tx *sqlx.Tx) error,
		opts ...TxOption,
	) error
	ExecInTxWithRetry(
		ctx context.Context,
		maxRetries int,
		fn func(ctx context.Context, tx *sqlx.Tx) error,
		opts ...TxOption,
	) error
	ExecSavepoint(
		ctx context.Context,
		tx *sqlx.Tx,
		name string,
		fn func(ctx context.Context) error,
	) error
}

type sqlxTxManager struct {
	db *sqlx.DB
}

// NewTxManager creates a new transaction manager backed by a sqlx.DB.
func NewTxManager(db *sqlx.DB) TxManager {
	return &sqlxTxManager{db: db}
}

// TxConfig holds transaction options and execution timeouts.
type TxConfig struct {
	Options sql.TxOptions
	Timeout time.Duration
}

// TxOption configures database transaction settings.
type TxOption func(*TxConfig)

// WithIsolationLevel sets the isolation level for the transaction.
func WithIsolationLevel(level sql.IsolationLevel) TxOption {
	return func(c *TxConfig) {
		c.Options.Isolation = level
	}
}

// WithReadOnly sets the read-only access mode for the transaction.
func WithReadOnly(readOnly bool) TxOption {
	return func(c *TxConfig) {
		c.Options.ReadOnly = readOnly
	}
}

// WithTimeout sets a custom timeout duration for the transaction context.
func WithTimeout(timeout time.Duration) TxOption {
	return func(c *TxConfig) {
		c.Timeout = timeout
	}
}

// ExecInTx executes a function within a database transaction with automatic rollback/commit handling.
func (m *sqlxTxManager) ExecInTx(
	ctx context.Context,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
	opts ...TxOption,
) error {
	cfg := TxConfig{
		Options: sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		},
		Timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	txCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	tx, err := m.db.BeginTxx(txCtx, &cfg.Options)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txCtx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil && !errors.Is(rbErr, sql.ErrTxDone) {
			return fmt.Errorf("transaction error: %v, rollback error: %w", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExecInTxWithRetry executes a transaction and retries on MySQL Deadlock (1213) or Lock Wait Timeout (1205).
func (m *sqlxTxManager) ExecInTxWithRetry(
	ctx context.Context,
	maxRetries int,
	fn func(ctx context.Context, tx *sqlx.Tx) error,
	opts ...TxOption,
) error {
	var err error
	backoff := 50 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err = m.ExecInTx(ctx, fn, opts...)
		if err == nil {
			return nil
		}

		// Retry on MySQL Deadlock (1213) or Lock Wait Timeout Exceeded (1205)
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && (mysqlErr.Number == 1213 || mysqlErr.Number == 1205) {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
				backoff *= 2
				continue
			}
		}

		return err
	}

	return fmt.Errorf("transaction aborted after %d retries: %w", maxRetries, err)
}

// ExecSavepoint executes a function within a named savepoint inside an active transaction.
func (m *sqlxTxManager) ExecSavepoint(
	ctx context.Context,
	tx *sqlx.Tx,
	name string,
	fn func(ctx context.Context) error,
) error {
	if _, err := tx.ExecContext(ctx, "SAVEPOINT "+name); err != nil {
		return fmt.Errorf("failed to establish savepoint %s: %w", name, err)
	}

	if err := fn(ctx); err != nil {
		_, _ = tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT "+name)
		return err
	}

	if _, err := tx.ExecContext(ctx, "RELEASE SAVEPOINT "+name); err != nil {
		return fmt.Errorf("failed to release savepoint %s: %w", name, err)
	}

	return nil
}
