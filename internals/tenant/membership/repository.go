package membership

import (
	"context"
	"time"

	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/google/uuid"
)

//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package=mocks

type Reader interface {
	GetMembership(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (tenantmodel.Membership, error)
	ListMemberships(ctx context.Context, tenantID uuid.UUID) ([]tenantmodel.Membership, error)
}
type Writer interface {
	CreateMembership(ctx context.Context, membership *tenantmodel.Membership) error
	UpdateMembershipRole(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, role tenantmodel.Role, updatedAt time.Time) error
	UpdateMembershipStatus(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, status tenantmodel.MembershipStatus, updatedAt time.Time) error
}

type Repository interface {
	Reader
	Writer
}
