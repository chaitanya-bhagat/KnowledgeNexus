package knowledgebase

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Repository interface {
	CreateKnowledgeBase(ctx context.Context, kb KnowledgeBase) error
	GetKnowledgeBaseByID(ctx context.Context, kbID uuid.UUID) (KnowledgeBase, error)
	GetKnowledgeBasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]KnowledgeBase, error)
	UpdateKnowledgeBase(ctx context.Context, kb *KnowledgeBase) error
	UpdateKnowledgeBaseStatus(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID, status Status, UpdatedAt time.Time) error
}
