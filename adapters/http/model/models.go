package httpmodel

import (
	"time"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
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
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
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

// type CreateMembershipRequest struct {
// 	UserID string `json:"userId"`
// 	Role   string `json:"role"`
// }

//	type UpdateMembershipRoleRequest struct {
//		Role string `json:"role"`
//	}
type MembershipResponse struct {
	ID        string                       `json:"id"`
	TenantID  string                       `json:"tenantID"`
	UserID    string                       `json:"userID"`
	Role      tenantmodel.Role             `json:"role"`
	Status    tenantmodel.MembershipStatus `json:"status"`
	CreatedAt time.Time                    `json:"createdAt"`
	UpdatedAt time.Time                    `json:"updatedAt"`
}

func NewMembershipResponse(m tenantmodel.Membership) MembershipResponse {
	return MembershipResponse{
		ID:        m.ID.String(),
		TenantID:  m.TenantID.String(),
		UserID:    m.UserID.String(),
		Role:      m.Role,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

type CreateMembershipRequest struct {
	TenantID string `json:"tenantID"`
	UserID   string `json:"userID"`
	Role     string `json:"role"`
}

type GetMembershipRequest struct {
	TenantID string `json:"tenantID"`
	UserID   string `json:"userID"`
}

type ListMembershipsRequest struct {
	TenantID string `json:"tenantID"`
}

type UpdateMembershipRoleRequest struct {
	TenantID string `json:"tenantID"`
	UserID   string `json:"userID"`
	Role     string `json:"role"`
}

type ChangeMembershipStatusRequest struct {
	TenantID string `json:"tenantID"`
	UserID   string `json:"userID"`
}

// Identity model

type CreateUserRequest struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type GetUserByIDRequest struct {
	UserID string `json:"userID"`
}

type GetUserByEmailRequest struct {
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	UserID      string `json:"userID"`
	DisplayName string `json:"displayName"`
}

type ChangeUserStatusRequest struct {
	UserID string `json:"userID"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"displayName"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ToUserResponse(user identity.User) UserResponse {
	return UserResponse{
		ID:          user.ID.String(),
		Email:       user.Email,
		DisplayName: user.DisplayName,
		Status:      string(user.Status),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
