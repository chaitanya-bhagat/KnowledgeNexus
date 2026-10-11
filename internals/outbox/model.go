package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	EventType   string
	AggregateID uuid.UUID
	Payload     json.RawMessage
	CreatedAt   time.Time
	PublishedAt *time.Time
}
