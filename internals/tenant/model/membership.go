package tenantmodel

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type MembershipStatus string

const (
	MembershipStatusActive   MembershipStatus = "active"
	MembershipStatusDisabled MembershipStatus = "disabled"
)

type Membership struct {
	ID        uuid.UUID
	TenantID  uuid.UUID
	UserID    uuid.UUID
	Role      Role
	Status    MembershipStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateMembershipInput struct {
	UserID uuid.UUID
	Role   Role
}

type UpdateMembershipRoleInput struct {
	Role Role
}
