package document

import (
	"context"
	"time"

	"github.com/google/uuid"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

type Reader interface {
	GetByID(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error)
	GetList(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID) ([]kbmodel.Document, error)
}
type Writer interface {
	Create(ctx context.Context, document kbmodel.Document) error
	Update(ctx context.Context, document kbmodel.Document) error
	UpdateStatus(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID, status kbmodel.DocumentStatus, updatedAt time.Time) error
}
type Repository interface {
	Reader
	Writer
}
