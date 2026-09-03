package knowledgebase

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	KnowledgeBaseStatusActive   Status = "active"
	KnowledgeBaseStatusArchived Status = "archived"
)

type KnowledgeBase struct {
	ID          string
	Name        string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
	TenantID    uuid.UUID
	CreatedBy   uuid.UUID
	DomainType  string
	Description string
}

type CreateKnowledgeBase struct {
	Name        string
	Status      Status
	TenantID    uuid.UUID
	CreatedBy   uuid.UUID
	DomainType  string
	Description string
}

type UpdateKnowledgeBase struct {
	Name        string
	Description string
}
