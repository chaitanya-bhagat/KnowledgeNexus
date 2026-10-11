package adapteroutbox

import (
	"context"
	"fmt"
	"uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/outbox"
	"github.com/jackc/pgx/v5/pgxpool"
)

type outboxRepository struct {
	db *pgxpool.Pool
}

func NewOutboxRepository(db *pgxpool.Pool) *outboxRepository {
	return &outboxRepository{
		db: db,
	}
}

func (or *outboxRepository) GetUnpublishedEvents(ctx context.Context, limit int) ([]outbox.Event, error) {
	if limit <= 0 {
		return nil, fmt.Errorf("limit must be greater than 0")
	}

	query := `
	SELECT id, tenant_id, event_type, aggregate_id, payload, created_at, published_at
	FROM table_outbox_events
	WHERE published_at IS NULL
	ORDER BY created_at, id
	LIMIT $1
	`
	rows, err := or.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query unpublished events: %w", err)
	}
	defer rows.Close()

	var events []outbox.Event

	for rows.Next() {
		var event outbox.Event
		err := rows.Scan(&event.ID, &event.TenantID, &event.EventType, &event.AggregateID, &event.Payload, &event.CreatedAt, &event.PublishedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event row: %w", err)
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over event rows: %w", err)
	}
	return events, nil
}

func (or *outboxRepository) MarkEventAsPublished(ctx context.Context, eventID uuid.UUID) error {
	query := `
	UPDATE table_outbox_events
	SET published_at = NOW()
	WHERE id = $1 AND published_at IS NULL
	`
	_, err := or.db.Exec(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("failed to mark event as published: %w", err)
	}
	return nil
}
