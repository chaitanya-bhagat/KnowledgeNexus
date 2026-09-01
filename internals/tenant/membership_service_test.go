package tenant_test

import (
	"context"
	"errors"
	"testing"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/mocks"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestMembershipService_Create(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()
	userID := uuid.New()

	activeTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusActive,
	}

	disabledTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusDisabled,
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name       string
		tenantID   uuid.UUID
		input      tenantmodel.CreateMembershipInput
		setupMocks func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository)
		wantRole   tenantmodel.Role
		wantErr    error
	}{
		{
			name:     "creates member successfully",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, m *tenantmodel.Membership) error {
					require.NotEqual(t, uuid.Nil, m.ID)
					require.Equal(t, tenantID, m.TenantID)
					require.Equal(t, userID, m.UserID)
					require.Equal(t, tenantmodel.RoleMember, m.Role)
					require.Equal(t, tenantmodel.MembershipStatusActive, m.Status)
					require.False(t, m.CreatedAt.IsZero())
					require.False(t, m.UpdatedAt.IsZero())

					return nil
				})
			},
			wantRole: tenantmodel.RoleMember,
		},
		{
			name:     "creates admin successfully",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantRole: tenantmodel.RoleAdmin,
		},
		{
			name:     "defaults empty role to member",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).DoAndReturn(func(_ context.Context, m *tenantmodel.Membership) error {
					require.Equal(t, tenantmodel.RoleMember, m.Role)
					return nil
				})
			},
			wantRole: tenantmodel.RoleMember,
		},
		{
			name:     "rejects nil tenant id",
			tenantID: uuid.Nil,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:     "rejects nil user id",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: uuid.Nil,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidUserID,
		},
		{
			name:     "rejects invalid role",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.Role("super-admin"),
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidRole,
		},
		{
			name:     "rejects owner role",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleOwner,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrOwnerRoleManagedSeparately,
		},
		{
			name:     "returns tenant not found",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, tenant.ErrNotFound)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name:     "rejects disabled tenant",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(disabledTenant, nil)
			},
			wantErr: tenant.ErrTenantDisabled,
		},
		{
			name:     "returns duplicate membership error",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).Return(tenant.ErrMembershipExists)
			},
			wantErr: tenant.ErrMembershipExists,
		},
		{
			name:     "returns user not found error",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).Return(tenant.ErrUserNotFound)
			},
			wantErr: tenant.ErrUserNotFound,
		},
		{
			name:     "returns repository error",
			tenantID: tenantID,
			input: tenantmodel.CreateMembershipInput{
				UserID: userID,
				Role:   tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().CreateMembership(gomock.Any(), gomock.Any()).Return(repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.Create(ctx, tt.tenantID, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tenantmodel.Membership{}, got)
				return
			}

			require.NoError(t, err)

			require.NotEqual(t, uuid.Nil, got.ID)
			require.Equal(t, tt.tenantID, got.TenantID)
			require.Equal(t, userID, got.UserID)
			require.Equal(t, tt.wantRole, got.Role)

			require.Equal(t, tenantmodel.MembershipStatusActive, got.Status)

			require.False(t, got.CreatedAt.IsZero())
			require.False(t, got.UpdatedAt.IsZero())
			require.Equal(t, got.CreatedAt, got.UpdatedAt)
		})
	}
}

func TestMembershipService_Get(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()
	userID := uuid.New()

	expectedMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusActive,
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name string

		tenantID   uuid.UUID
		userID     uuid.UUID
		setupMocks func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository)
		want       tenantmodel.Membership
		wantErr    error
	}{
		{
			name:     "gets membership successfully",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(expectedMembership, nil)
			},
			want: expectedMembership,
		},
		{
			name:       "rejects nil tenant id",
			tenantID:   uuid.Nil,
			userID:     userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:       "rejects nil user id",
			tenantID:   tenantID,
			userID:     uuid.Nil,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidUserID,
		},
		{
			name:     "returns membership not found",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, tenant.ErrMembershipNotFound)
			},
			wantErr: tenant.ErrMembershipNotFound,
		},
		{
			name:     "returns repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.Get(ctx, tt.tenantID, tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tenantmodel.Membership{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMembershipService_List(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()

	activeTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusActive,
	}

	expectedMemberships := []tenantmodel.Membership{
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   uuid.New(),
			Role:     tenantmodel.RoleOwner,
			Status:   tenantmodel.MembershipStatusActive,
		},
		{
			ID:       uuid.New(),
			TenantID: tenantID,
			UserID:   uuid.New(),
			Role:     tenantmodel.RoleMember,
			Status:   tenantmodel.MembershipStatusActive,
		},
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name       string
		tenantID   uuid.UUID
		setupMocks func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository)
		want       []tenantmodel.Membership
		wantErr    error
	}{
		{
			name:     "lists memberships successfully",
			tenantID: tenantID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().ListMemberships(gomock.Any(), tenantID).Return(expectedMemberships, nil)
			},
			want: expectedMemberships,
		},
		{
			name:     "returns empty memberships",
			tenantID: tenantID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)

				membershipRepo.EXPECT().ListMemberships(gomock.Any(), tenantID).Return([]tenantmodel.Membership{}, nil)
			},
			want: []tenantmodel.Membership{},
		},
		{
			name:       "rejects nil tenant id",
			tenantID:   uuid.Nil,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:     "returns tenant not found",
			tenantID: tenantID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, tenant.ErrNotFound)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name:     "returns tenant repository error",
			tenantID: tenantID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, repoErr)
			},
			wantErr: repoErr,
		},
		{
			name:     "returns membership repository error",
			tenantID: tenantID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().ListMemberships(gomock.Any(), tenantID).Return(nil, repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.List(ctx, tt.tenantID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestMembershipService_UpdateRole(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()
	userID := uuid.New()

	activeTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusActive,
	}

	disabledTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusDisabled,
	}

	memberMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusActive,
	}

	adminMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleAdmin,
		Status:   tenantmodel.MembershipStatusActive,
	}

	ownerMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleOwner,
		Status:   tenantmodel.MembershipStatusActive,
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name string

		tenantID uuid.UUID
		userID   uuid.UUID
		input    tenantmodel.UpdateMembershipRoleInput

		setupMocks func(
			tenantRepo *mocks.MockRepository,
			membershipRepo *mocks.MockMembershipRepository,
		)

		wantRole tenantmodel.Role
		wantErr  error
	}{
		{
			name:     "updates member to admin successfully",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(memberMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipRole(gomock.Any(), tenantID, userID, tenantmodel.RoleAdmin, gomock.Any()).Return(nil)
			},
			wantRole: tenantmodel.RoleAdmin,
		},
		{
			name:     "updates admin to member successfully",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(adminMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipRole(gomock.Any(), tenantID, userID, tenantmodel.RoleMember, gomock.Any()).Return(nil)
			},
			wantRole: tenantmodel.RoleMember,
		},
		{
			name:     "same role is idempotent",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleMember,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(memberMembership, nil)
			},
			wantRole: tenantmodel.RoleMember,
		},
		{
			name:     "rejects nil tenant id",
			tenantID: uuid.Nil,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:     "rejects nil user id",
			tenantID: tenantID,
			userID:   uuid.Nil,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
			},
			wantErr: tenant.ErrInvalidUserID,
		},
		{
			name:     "rejects invalid role",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.Role("super-admin"),
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidRole,
		},
		{
			name:     "rejects changing role to owner",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleOwner,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrOwnerRoleManagedSeparately,
		},
		{
			name:     "returns tenant not found",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, tenant.ErrNotFound)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name:     "rejects disabled tenant",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(disabledTenant, nil)
			},
			wantErr: tenant.ErrTenantDisabled,
		},
		{
			name:     "returns membership not found",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, tenant.ErrMembershipNotFound)
			},
			wantErr: tenant.ErrMembershipNotFound,
		},
		{
			name:     "rejects modifying owner membership",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(ownerMembership, nil)
			},
			wantErr: tenant.ErrOwnerRoleManagedSeparately,
		},
		{
			name:     "returns update repository error",
			tenantID: tenantID,
			userID:   userID,
			input: tenantmodel.UpdateMembershipRoleInput{
				Role: tenantmodel.RoleAdmin,
			},
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(memberMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipRole(gomock.Any(), tenantID, userID, tenantmodel.RoleAdmin, gomock.Any()).Return(repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.UpdateRole(ctx, tt.tenantID, tt.userID, tt.input)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tenantmodel.Membership{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.tenantID, got.TenantID)
			require.Equal(t, tt.userID, got.UserID)
			require.Equal(t, tt.wantRole, got.Role)
		})
	}
}

func TestMembershipService_Disable(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()
	userID := uuid.New()

	activeMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusActive,
	}

	disabledMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusDisabled,
	}

	ownerMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleOwner,
		Status:   tenantmodel.MembershipStatusActive,
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name string

		tenantID uuid.UUID
		userID   uuid.UUID

		setupMocks func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository)

		wantStatus tenantmodel.MembershipStatus
		wantErr    error
	}{
		{
			name:     "disables active membership successfully",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(activeMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipStatus(gomock.Any(), tenantID, userID, tenantmodel.MembershipStatusDisabled, gomock.Any()).Return(nil)
			},
			wantStatus: tenantmodel.MembershipStatusDisabled,
		},
		{
			name:     "already disabled membership is idempotent",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(disabledMembership, nil)
			},
			wantStatus: tenantmodel.MembershipStatusDisabled,
		},
		{
			name:       "rejects nil tenant id",
			tenantID:   uuid.Nil,
			userID:     userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:       "rejects nil user id",
			tenantID:   tenantID,
			userID:     uuid.Nil,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidUserID,
		},
		{
			name:     "returns membership not found",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.
					EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, tenant.ErrMembershipNotFound)
			},
			wantErr: tenant.ErrMembershipNotFound,
		},
		{
			name:     "rejects disabling owner membership",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(ownerMembership, nil)
			},
			wantErr: tenant.ErrOwnerRoleManagedSeparately,
		},
		{
			name:     "returns get repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, repoErr)
			},
			wantErr: repoErr,
		},
		{
			name:     "returns update status repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(activeMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipStatus(gomock.Any(), tenantID, userID, tenantmodel.MembershipStatusDisabled, gomock.Any()).Return(repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.Disable(ctx, tt.tenantID, tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tenantmodel.Membership{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.tenantID, got.TenantID)
			require.Equal(t, tt.userID, got.UserID)
			require.Equal(t, tt.wantStatus, got.Status)
		})
	}
}

func TestMembershipService_Enable(t *testing.T) {
	ctx := context.Background()

	tenantID := uuid.New()
	userID := uuid.New()

	activeTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusActive,
	}

	disabledTenant := &tenantmodel.Tenant{
		ID:     tenantID,
		Name:   "Acme",
		Slug:   "acme",
		Status: tenantmodel.StatusDisabled,
	}

	disabledMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusDisabled,
	}

	activeMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleMember,
		Status:   tenantmodel.MembershipStatusActive,
	}

	ownerMembership := tenantmodel.Membership{
		ID:       uuid.New(),
		TenantID: tenantID,
		UserID:   userID,
		Role:     tenantmodel.RoleOwner,
		Status:   tenantmodel.MembershipStatusDisabled,
	}

	repoErr := errors.New("repository error")

	tests := []struct {
		name string

		tenantID uuid.UUID
		userID   uuid.UUID

		setupMocks func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository)

		wantStatus tenantmodel.MembershipStatus
		wantErr    error
	}{
		{
			name:     "enables disabled membership successfully",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(disabledMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipStatus(gomock.Any(), tenantID, userID, tenantmodel.MembershipStatusActive, gomock.Any()).Return(nil)
			},
			wantStatus: tenantmodel.MembershipStatusActive,
		},
		{
			name:     "already active membership is idempotent",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(activeMembership, nil)
			},
			wantStatus: tenantmodel.MembershipStatusActive,
		},
		{
			name:       "rejects nil tenant id",
			tenantID:   uuid.Nil,
			userID:     userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidTenantID,
		},
		{
			name:       "rejects nil user id",
			tenantID:   tenantID,
			userID:     uuid.Nil,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {},
			wantErr:    tenant.ErrInvalidUserID,
		},
		{
			name:     "returns tenant not found",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, tenant.ErrNotFound)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name:     "rejects disabled tenant",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(tenantRepo *mocks.MockRepository, membershipRepo *mocks.MockMembershipRepository) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(disabledTenant, nil)
			},
			wantErr: tenant.ErrTenantDisabled,
		},
		{
			name:     "returns membership not found",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository,
				membershipRepo *mocks.MockMembershipRepository,
			) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, tenant.ErrMembershipNotFound)
			},
			wantErr: tenant.ErrMembershipNotFound,
		},
		{
			name:     "rejects enabling owner membership",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository,
				membershipRepo *mocks.MockMembershipRepository,
			) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(ownerMembership, nil)
			},
			wantErr: tenant.ErrOwnerRoleManagedSeparately,
		},
		{
			name:     "returns tenant repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository,
				membershipRepo *mocks.MockMembershipRepository,
			) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(&tenantmodel.Tenant{}, repoErr)
			},
			wantErr: repoErr,
		},
		{
			name:     "returns get membership repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository,
				membershipRepo *mocks.MockMembershipRepository,
			) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(tenantmodel.Membership{}, repoErr)
			},
			wantErr: repoErr,
		},
		{
			name:     "returns update status repository error",
			tenantID: tenantID,
			userID:   userID,
			setupMocks: func(
				tenantRepo *mocks.MockRepository,
				membershipRepo *mocks.MockMembershipRepository,
			) {
				tenantRepo.EXPECT().GetByID(gomock.Any(), tenantID).Return(activeTenant, nil)
				membershipRepo.EXPECT().GetMembership(gomock.Any(), tenantID, userID).Return(disabledMembership, nil)
				membershipRepo.EXPECT().UpdateMembershipStatus(gomock.Any(), tenantID, userID, tenantmodel.MembershipStatusActive, gomock.Any()).Return(repoErr)
			},
			wantErr: repoErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			tenantRepo := mocks.NewMockRepository(ctrl)
			membershipRepo := mocks.NewMockMembershipRepository(ctrl)

			tt.setupMocks(tenantRepo, membershipRepo)

			service := tenant.NewMembershipService(tenantRepo, membershipRepo)

			got, err := service.Enable(ctx, tt.tenantID, tt.userID)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, tenantmodel.Membership{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.tenantID, got.TenantID)
			require.Equal(t, tt.userID, got.UserID)
			require.Equal(t, tt.wantStatus, got.Status)
		})
	}
}
