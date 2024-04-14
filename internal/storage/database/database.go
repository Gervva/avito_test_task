package database

import (
	"context"
	"database/sql"
)

type Database struct {
	db *sql.DB
}

func NewDatabase(db *sql.DB) Database {
	return Database{
		db: db,
	}
}

type ctxKey string

const txKey = ctxKey("transaction")

type TransactionWrapper func(ctx context.Context, f func(ctx context.Context) error) error

func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if tx := extractTx(ctx); tx != nil {
		_, err := tx.ExecContext(ctx, query, args...)
		return nil, err
	}

	affectedRows, err := db.db.ExecContext(ctx, query, args...)
	return affectedRows, err
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryContext(ctx, query, args...)
	}

	return db.db.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if tx := extractTx(ctx); tx != nil {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return db.db.QueryRowContext(ctx, query, args...)
}

func (db *Database) WithTransaction(ctx context.Context, f func(ctx context.Context) error) error {
	tx, err := db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() { _ = tx.Rollback() }()

	err = f(context.WithValue(ctx, txKey, tx))
	if err != nil {
		return err
	}

	return tx.Commit()
}

// extractTx checks if there is transaction in context. If transaction in context, it returns transaction struct to
// perform operation in transactions
func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey).(*sql.Tx); ok {
		return tx
	}
	return nil
}