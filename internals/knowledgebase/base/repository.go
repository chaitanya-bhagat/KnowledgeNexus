package knowledgebase

import (
	"context"
	"time"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/google/uuid"
)

type Repository interface {
	CreateKnowledgeBase(ctx context.Context, kb kbmodel.KnowledgeBase) error
	GetKnowledgeBaseByID(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID) (kbmodel.KnowledgeBase, error)
	ListKnowledgeBasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]kbmodel.KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, kb *kbmodel.KnowledgeBase) error
	UpdateKnowledgeBaseStatus(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID, status kbmodel.Status, updatedAt time.Time) error
}
