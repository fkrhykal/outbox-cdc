package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fkrhykal/outbox-cdc/internal/messaging"
	"github.com/fkrhykal/outbox-cdc/internal/outbox"
)

var _ outbox.OutboxRepository = (*PgOutboxRepository)(nil)

//	TABLE public.outbox (
//
//	id UUID PRIMARY KEY,
//	aggregate_id VARCHAR(255) NOT NULL,
//	aggregate_type VARCHAR(255) NOT NULL,
//	type VARCHAR(255) NOT NULL,
//	payload JSONB NOT NULL,
//	timestamp TIMESTAMP WITHOUT TIME ZONE DEFAULT now() NOT NULL
//
// );
// CREATE INDEX outbox_aggregate_id_idx ON public.outbox (aggregate_id);
type PgOutboxRepository struct {
	pgRepository
}

// SaveEvent implements [outbox.OutboxRepository].
func (p *PgOutboxRepository) SaveEvent(ctx context.Context, event messaging.Event) error {
	outboxRecord, err := outbox.Event(event)
	if err != nil {
		return fmt.Errorf("failed to map event to outbox record: %w", err)
	}
	return p.Save(ctx, outboxRecord)
}

// Save implements outbox.OutboxPersistence.
func (p *PgOutboxRepository) Save(ctx context.Context, outbox *outbox.Outbox) error {
	query := `
		INSERT INTO outbox(id, type, aggregateid, aggregatetype, payload)
		VALUES($1, $2, $3, $4, $5)
	`
	_, err := p.Executor(ctx).
		ExecContext(ctx, query, outbox.EventID, outbox.EventType, outbox.AggregateID, outbox.AggregateType, outbox.Payload)
	if err != nil {
		return fmt.Errorf("failed to insert outbox record: %w", err)
	}
	return nil
}

func NewPgOutboxRepository(db *sql.DB) *PgOutboxRepository {
	return &PgOutboxRepository{
		pgRepository: pgRepository{
			db: db,
		},
	}
}
