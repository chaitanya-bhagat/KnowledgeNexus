package document

import (
	"context"
	"time"

	"github.com/google/uuid"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
)

type Repository interface {
	Create(ctx context.Context, document kbmodel.Document) error
	GetByID(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error)
	GetList(ctx context.Context, tenantID uuid.UUID) ([]kbmodel.Document, error)
	Update(ctx context.Context, document kbmodel.Document) error
	UpdateStatus(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID, status kbmodel.DocumentStatus, updatedAt time.Time) error
}
