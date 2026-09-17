package documentversion

import (
	"context"
	"uuid"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
)

type Repository interface {
	Create(ctx context.Context, version kbmodel.DocumentVersion) error
	GetByID(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) (kbmodel.DocumentVersion, error)
	GetList(ctx context.Context, docID uuid.UUID) ([]kbmodel.DocumentVersion, error)
	Update(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID, status kbmodel.DocumentVersionStatus, failureReason string) error
}
