package tenant

import (
	"context"
	"time"

	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/google/uuid"
)

type MembershipService struct {
	tenantRepo     Repository
	membershipRepo MembershipRepository
}

func NewMembershipService(tenantRepo Repository, membershipRepo MembershipRepository) *MembershipService {
	return &MembershipService{
		tenantRepo:     tenantRepo,
		membershipRepo: membershipRepo,
	}
}

func (s *MembershipService) Create(ctx context.Context, tenantID uuid.UUID, input tenantmodel.CreateMembershipInput) (tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidTenantID
	}

	if input.UserID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidUserID
	}

	role := input.Role

	if role == "" {
		role = tenantmodel.RoleMember
	}

	if !isValidMembershipRole(role) {
		return tenantmodel.Membership{}, ErrInvalidRole
	}

	// Owner assignment is deliberately not handled by the generic API.
	if role == tenantmodel.RoleOwner {
		return tenantmodel.Membership{}, ErrOwnerRoleManagedSeparately
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if tenant.Status == tenantmodel.StatusDisabled {
		return tenantmodel.Membership{}, ErrTenantDisabled
	}

	now := time.Now().UTC()

	membership := tenantmodel.Membership{
		ID:        uuid.New(),
		TenantID:  tenantID,
		UserID:    input.UserID,
		Role:      role,
		Status:    tenantmodel.MembershipStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.membershipRepo.CreateMembership(
		ctx,
		&membership,
	); err != nil {
		return tenantmodel.Membership{}, err
	}

	return membership, nil
}

func (s *MembershipService) Get(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidTenantID
	}

	if userID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidUserID
	}

	return s.membershipRepo.GetMembership(
		ctx,
		tenantID,
		userID,
	)
}

func (s *MembershipService) List(ctx context.Context, tenantID uuid.UUID) ([]tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return nil, ErrInvalidTenantID
	}

	if _, err := s.tenantRepo.GetByID(ctx, tenantID); err != nil {
		return nil, err
	}

	return s.membershipRepo.ListMemberships(
		ctx,
		tenantID,
	)
}

func (s *MembershipService) UpdateRole(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, input tenantmodel.UpdateMembershipRoleInput) (tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidTenantID
	}

	if userID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidUserID
	}

	if !isValidMembershipRole(input.Role) {
		return tenantmodel.Membership{}, ErrInvalidRole
	}

	if input.Role == tenantmodel.RoleOwner {
		return tenantmodel.Membership{}, ErrOwnerRoleManagedSeparately
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if tenant.Status == tenantmodel.StatusDisabled {
		return tenantmodel.Membership{}, ErrTenantDisabled
	}

	membership, err := s.membershipRepo.GetMembership(
		ctx,
		tenantID,
		userID,
	)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if membership.Role == tenantmodel.RoleOwner {
		return tenantmodel.Membership{}, ErrOwnerRoleManagedSeparately
	}

	if membership.Role == input.Role {
		return membership, nil
	}

	now := time.Now().UTC()

	if err := s.membershipRepo.UpdateMembershipRole(
		ctx,
		tenantID,
		userID,
		input.Role,
		now,
	); err != nil {
		return tenantmodel.Membership{}, err
	}

	membership.Role = input.Role
	membership.UpdatedAt = now

	return membership, nil
}

func (s *MembershipService) Disable(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidTenantID
	}

	if userID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidUserID
	}

	membership, err := s.membershipRepo.GetMembership(ctx, tenantID, userID)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if membership.Role == tenantmodel.RoleOwner {
		return tenantmodel.Membership{}, ErrOwnerRoleManagedSeparately
	}

	if membership.Status == tenantmodel.MembershipStatusDisabled {
		return membership, nil
	}

	now := time.Now().UTC()

	if err := s.membershipRepo.UpdateMembershipStatus(
		ctx,
		tenantID,
		userID,
		tenantmodel.MembershipStatusDisabled,
		now,
	); err != nil {
		return tenantmodel.Membership{}, err
	}

	membership.Status = tenantmodel.MembershipStatusDisabled
	membership.UpdatedAt = now

	return membership, nil
}

func (s *MembershipService) Enable(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (tenantmodel.Membership, error) {
	if tenantID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidTenantID
	}

	if userID == uuid.Nil {
		return tenantmodel.Membership{}, ErrInvalidUserID
	}

	tenant, err := s.tenantRepo.GetByID(ctx, tenantID)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if tenant.Status == tenantmodel.StatusDisabled {
		return tenantmodel.Membership{}, ErrTenantDisabled
	}

	membership, err := s.membershipRepo.GetMembership(
		ctx,
		tenantID,
		userID,
	)
	if err != nil {
		return tenantmodel.Membership{}, err
	}

	if membership.Role == tenantmodel.RoleOwner {
		return tenantmodel.Membership{}, ErrOwnerRoleManagedSeparately
	}

	if membership.Status == tenantmodel.MembershipStatusActive {
		return membership, nil
	}

	now := time.Now().UTC()

	if err := s.membershipRepo.UpdateMembershipStatus(
		ctx,
		tenantID,
		userID,
		tenantmodel.MembershipStatusActive,
		now,
	); err != nil {
		return tenantmodel.Membership{}, err
	}

	membership.Status = tenantmodel.MembershipStatusActive
	membership.UpdatedAt = now

	return membership, nil
}

func isValidMembershipRole(role tenantmodel.Role) bool {
	switch role {
	case tenantmodel.RoleOwner, tenantmodel.RoleAdmin, tenantmodel.RoleMember:
		return true
	default:
		return false
	}
}
