package data

import (
	"context"
	"database/sql"
)

type TxContext[T any] interface {
	context.Context
	Rollback() error
	Commit() error
	Executor() T
}

type TxOption interface {
	apply(opt TxOption)
}

type TxManager[T any] interface {
	Begin(ctx context.Context, opt ...TxOption) (TxContext[T], error)
}

var _ SqlExecutor = (*sql.DB)(nil)
var _ SqlExecutor = (*sql.Tx)(nil)

type SqlExecutor interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

var _ TxContext[SqlExecutor] = (*SqlTxContext)(nil)

type SqlTxContext struct {
	context.Context
	tx *sql.Tx
}

// Commit implements TxContext.
func (s *SqlTxContext) Commit() error {
	return s.tx.Commit()
}

// Executor implements TxContext.
func (s *SqlTxContext) Executor() SqlExecutor {
	return s.tx
}

// Rollback implements TxContext.
func (s *SqlTxContext) Rollback() error {
	return s.tx.Rollback()
}

var _ TxManager[SqlExecutor] = (*SqlTxManager)(nil)

type SqlTxManager struct {
	db *sql.DB
}

func NewSqlTxManager(db *sql.DB) *SqlTxManager {
	return &SqlTxManager{db: db}
}

// Begin implements TxManager.
func (s *SqlTxManager) Begin(ctx context.Context, opt ...TxOption) (TxContext[SqlExecutor], error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}
	return &SqlTxContext{Context: ctx, tx: tx}, nil
}
