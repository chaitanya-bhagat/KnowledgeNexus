package documentupload

import (
	"context"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, upload kbmodel.DocumentUpload) error
	GetByID(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) (kbmodel.DocumentUpload, error)
	MarkCompleted(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error
	MarkExpired(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error
}
