package knowledgebase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateKnowledgeBase(ctx context.Context, kb KnowledgeBase) error
	GetKnowledgeBaseByID(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID) (KnowledgeBase, error)
	ListKnowledgeBasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, kb *KnowledgeBase) error
	UpdateKnowledgeBaseStatus(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID, status Status, updatedAt time.Time) error
}
