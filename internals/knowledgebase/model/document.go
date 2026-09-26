package kbmodel

import (
	"time"

	"github.com/google/uuid"
)

type DocumentStatus string

const (
	DocumentStatusActive   DocumentStatus = "active"
	DocumentStatusArchived DocumentStatus = "archived"
)

type Document struct {
	ID           uuid.UUID
	Title        string
	TenantID     uuid.UUID
	KbID         uuid.UUID
	CreatedBy    uuid.UUID
	DocumentType string
	Status       DocumentStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateDocument struct {
	TenantID        uuid.UUID
	KnowledgeBaseID uuid.UUID
	Title           string
	DocumentType    string
	CreatedBy       uuid.UUID
}

type DocumentInput struct {
	KnowledgeBaseID uuid.UUID
	Title           string
	DocumentType    string
}

type UpdateDocumentInput struct {
	Title        string
	DocumentType string
}
