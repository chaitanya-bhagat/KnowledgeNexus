package httpmodel

import (
	"time"

	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
)

type CreateTenantRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateTenantRequest struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type TenantResponse struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Slug      string             `json:"slug"`
	Status    tenantmodel.Status `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

func NewTenantResponse(t *tenantmodel.Tenant) *TenantResponse {
	return &TenantResponse{
		ID:        t.ID.String(),
		Name:      t.Name,
		Slug:      t.Slug,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
