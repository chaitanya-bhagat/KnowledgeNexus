package documentupload

import (
	"context"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

type Reader interface {
	GetByID(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) (kbmodel.DocumentUpload, error)
}
type Writer interface {
	Create(ctx context.Context, upload kbmodel.DocumentUpload) error
	// MarkCompleted(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error
	// MarkExpired(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error
}
type Repository interface {
	Reader
	Writer
}
