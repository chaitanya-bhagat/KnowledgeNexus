package tenant

import (
	"context"
	"time"

	"github.com/google/uuid"

	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

type Repository interface {
	Create(ctx context.Context, tenant *tenantmodel.Tenant) error
	GetByID(ctx context.Context, id uuid.UUID) (*tenantmodel.Tenant, error)
	Update(ctx context.Context, tenant *tenantmodel.Tenant) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status tenantmodel.Status, updatedAt time.Time) error
}
