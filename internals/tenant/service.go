package tenant

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	// "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
)

type TenantService struct {
	repo Repository
}

func NewTenantService(repo Repository) *TenantService {
	return &TenantService{
		repo: repo,
	}
}

func (ts *TenantService) Create(ctx context.Context, input tenantmodel.CreateInput) (*tenantmodel.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &tenantmodel.Tenant{}, ErrInvalidName
	}
	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" {
		return &tenantmodel.Tenant{}, ErrInvalidSlug
	}

	now := time.Now().UTC()

	t := &tenantmodel.Tenant{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		Status:    tenantmodel.StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := ts.repo.Create(ctx, t); err != nil {
		return &tenantmodel.Tenant{}, err
	}
	return t, nil

}

func (ts *TenantService) GetByID(ctx context.Context, id uuid.UUID) (*tenantmodel.Tenant, error) {
	return ts.repo.GetByID(ctx, id)
}

func (ts *TenantService) Update(ctx context.Context, id uuid.UUID, input tenantmodel.UpdateInput) (*tenantmodel.Tenant, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return &tenantmodel.Tenant{}, ErrInvalidName
	}

	slug := strings.ToLower(strings.TrimSpace(input.Slug))
	if slug == "" {
		return &tenantmodel.Tenant{}, ErrInvalidSlug
	}

	t, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return &tenantmodel.Tenant{}, err
	}

	t.Name = name
	t.Slug = slug
	t.UpdatedAt = time.Now().UTC()

	if err := ts.repo.Update(ctx, t); err != nil {
		return &tenantmodel.Tenant{}, err
	}

	return t, nil
}

func (ts *TenantService) Disable(ctx context.Context, id uuid.UUID) (*tenantmodel.Tenant, error) {
	t, err := ts.repo.GetByID(ctx, id)
	if err != nil {
		return &tenantmodel.Tenant{}, err
	}

	if t.Status == tenantmodel.StatusDisabled {
		return t, nil
	}

	now := time.Now().UTC()

	if err := ts.repo.UpdateStatus(ctx, id, tenantmodel.StatusDisabled, now); err != nil {
		return &tenantmodel.Tenant{}, err
	}

	t.Status = tenantmodel.StatusDisabled
	t.UpdatedAt = now

	return t, nil
}
