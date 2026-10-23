package outbox

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	GetUnpublishedEvents(ctx context.Context, limit int) ([]Event, error)
	MarkEventAsPublished(ctx context.Context, eventID uuid.UUID) error
}
