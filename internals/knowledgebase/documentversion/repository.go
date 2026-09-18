package documentversion

import (
	"context"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, version kbmodel.DocumentVersion) error
	GetByID(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) (kbmodel.DocumentVersion, error)
	GetList(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) ([]kbmodel.DocumentVersion, error)
	Update(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID, status kbmodel.DocumentVersionStatus, failureReason string) error
}
