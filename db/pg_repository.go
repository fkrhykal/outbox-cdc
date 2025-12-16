package db

import (
	"context"
	"database/sql"

	"github.com/fkrhykal/outbox-cdc/data"
)

type pgRepository struct{ db *sql.DB }

func (p *pgRepository) Executor(ctx context.Context) data.SqlExecutor {
	if txCtx, ok := ctx.(data.SqlTxContext); ok {
		return txCtx.Executor()
	}
	return p.db
}
