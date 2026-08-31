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

func TestService_Create(t *testing.T) {
	repositoryErr := errors.New("database error")

	tests := []struct {
		name      string
		input     tenantmodel.CreateInput
		setupMock func(*mocks.MockRepository)
		wantErr   error
		assert    func(*testing.T, *tenantmodel.Tenant)
	}{
		{
			name: "success",
			input: tenantmodel.CreateInput{
				Name: "Acme Legal",
				Slug: "acme-legal",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(t, "Acme Legal", got.Name)
				require.Equal(t, "acme-legal", got.Slug)
				require.Equal(t, tenantmodel.StatusActive, got.Status)
				require.NotEqual(t, uuid.Nil, got.ID)
				require.False(t, got.CreatedAt.IsZero())
				require.False(t, got.UpdatedAt.IsZero())
			},
		},
		{
			name: "trims name and normalizes slug",
			input: tenantmodel.CreateInput{
				Name: "  Acme Legal  ",
				Slug: "  ACME-LEGAL  ",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						tenant *tenantmodel.Tenant,
					) error {
						require.Equal(t, "Acme Legal", tenant.Name)
						require.Equal(t, "acme-legal", tenant.Slug)

						return nil
					}).
					Times(1)
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(t, "Acme Legal", got.Name)
				require.Equal(t, "acme-legal", got.Slug)
			},
		},
		{
			name: "empty name",
			input: tenantmodel.CreateInput{
				Name: "   ",
				Slug: "acme",
			},
			wantErr: tenant.ErrInvalidName,
		},
		{
			name: "empty slug",
			input: tenantmodel.CreateInput{
				Name: "Acme",
				Slug: "   ",
			},
			wantErr: tenant.ErrInvalidSlug,
		},
		{
			name: "repository error",
			input: tenantmodel.CreateInput{
				Name: "Acme",
				Slug: "acme",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(repositoryErr).
					Times(1)
			},
			wantErr: repositoryErr,
		},
		{
			name: "slug conflict",
			input: tenantmodel.CreateInput{
				Name: "Acme",
				Slug: "acme",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tenant.ErrSlugConflict).
					Times(1)
			},
			wantErr: tenant.ErrSlugConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mocks.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			service := tenant.NewTenantService(repo)

			got, err := service.Create(
				context.Background(),
				tt.input,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, &tenantmodel.Tenant{}, got)
				return
			}

			require.NoError(t, err)

			if tt.assert != nil {
				tt.assert(t, got)
			}
		})
	}
}

func TestService_GetByID(t *testing.T) {
	id := uuid.New()
	repositoryErr := errors.New("database error")

	expectedTenant := &tenantmodel.Tenant{
		ID:     id,
		Name:   "Acme Legal",
		Slug:   "acme-legal",
		Status: tenantmodel.StatusActive,
	}

	tests := []struct {
		name      string
		setupMock func(*mocks.MockRepository)
		want      *tenantmodel.Tenant
		wantErr   error
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(expectedTenant, nil).
					Times(1)
			},
			want: expectedTenant,
		},
		{
			name: "tenant not found",
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name: "repository error",
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(
						&tenantmodel.Tenant{},
						repositoryErr,
					).
					Times(1)
			},
			wantErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mocks.NewMockRepository(ctrl)

			tt.setupMock(repo)

			service := tenant.NewTenantService(repo)

			got, err := service.GetByID(
				context.Background(),
				id,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, &tenantmodel.Tenant{}, got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestService_Update(t *testing.T) {
	id := uuid.New()

	repositoryErr := errors.New("database error")

	existing := &tenantmodel.Tenant{
		ID:     id,
		Name:   "Old Name",
		Slug:   "old-name",
		Status: tenantmodel.StatusActive,
	}

	tests := []struct {
		name      string
		input     tenantmodel.UpdateInput
		setupMock func(*mocks.MockRepository)
		wantErr   error
		assert    func(*testing.T, *tenantmodel.Tenant)
	}{
		{
			name: "success",
			input: tenantmodel.UpdateInput{
				Name: "New Name",
				Slug: "new-name",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(func(
						ctx context.Context,
						tenant *tenantmodel.Tenant,
					) error {
						require.Equal(t, id, tenant.ID)
						require.Equal(t, "New Name", tenant.Name)
						require.Equal(t, "new-name", tenant.Slug)
						require.False(t, tenant.UpdatedAt.IsZero())

						return nil
					}).
					Times(1)
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(t, id, got.ID)
				require.Equal(t, "New Name", got.Name)
				require.Equal(t, "new-name", got.Slug)
			},
		},
		{
			name: "normalizes input",
			input: tenantmodel.UpdateInput{
				Name: "  New Name  ",
				Slug: "  NEW-NAME  ",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(t, "New Name", got.Name)
				require.Equal(t, "new-name", got.Slug)
			},
		},
		{
			name: "empty name",
			input: tenantmodel.UpdateInput{
				Name: "",
				Slug: "acme",
			},
			wantErr: tenant.ErrInvalidName,
		},
		{
			name: "empty slug",
			input: tenantmodel.UpdateInput{
				Name: "Acme",
				Slug: "",
			},
			wantErr: tenant.ErrInvalidSlug,
		},
		{
			name: "tenant not found",
			input: tenantmodel.UpdateInput{
				Name: "Acme",
				Slug: "acme",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name: "slug conflict",
			input: tenantmodel.UpdateInput{
				Name: "Acme",
				Slug: "existing-slug",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(tenant.ErrSlugConflict).
					Times(1)
			},
			wantErr: tenant.ErrSlugConflict,
		},
		{
			name: "repository update error",
			input: tenantmodel.UpdateInput{
				Name: "Acme",
				Slug: "acme",
			},
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(repositoryErr).
					Times(1)
			},
			wantErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mocks.NewMockRepository(ctrl)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			service := tenant.NewTenantService(repo)

			got, err := service.Update(
				context.Background(),
				id,
				tt.input,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, &tenantmodel.Tenant{}, got)
				return
			}

			require.NoError(t, err)

			if tt.assert != nil {
				tt.assert(t, got)
			}
		})
	}
}

func TestService_Disable(t *testing.T) {
	id := uuid.New()

	repositoryErr := errors.New("database error")

	newActiveTenant := func() *tenantmodel.Tenant {
		return &tenantmodel.Tenant{
			ID:     id,
			Name:   "Acme",
			Slug:   "acme",
			Status: tenantmodel.StatusActive,
		}
	}

	newDisabledTenant := func() *tenantmodel.Tenant {
		return &tenantmodel.Tenant{
			ID:     id,
			Name:   "Acme",
			Slug:   "acme",
			Status: tenantmodel.StatusDisabled,
		}
	}

	tests := []struct {
		name      string
		setupMock func(*mocks.MockRepository)
		wantErr   error
		assert    func(*testing.T, *tenantmodel.Tenant)
	}{
		{
			name: "success",
			setupMock: func(repo *mocks.MockRepository) {
				activeTenant := newActiveTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(activeTenant, nil).
					Times(1)

				repo.EXPECT().
					UpdateStatus(
						gomock.Any(),
						id,
						tenantmodel.StatusDisabled,
						gomock.Any(),
					).
					Return(nil).
					Times(1)
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(
					t,
					tenantmodel.StatusDisabled,
					got.Status,
				)

				require.False(t, got.UpdatedAt.IsZero())
			},
		},
		{
			name: "already disabled is idempotent",
			setupMock: func(repo *mocks.MockRepository) {
				disabledTenant := newDisabledTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(disabledTenant, nil).
					Times(1)

				// Intentionally no UpdateStatus expectation.
			},
			assert: func(t *testing.T, got *tenantmodel.Tenant) {
				require.Equal(
					t,
					tenantmodel.StatusDisabled,
					got.Status,
				)
			},
		},
		{
			name: "tenant not found",
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantErr: tenant.ErrNotFound,
		},
		{
			name: "update status fails",
			setupMock: func(repo *mocks.MockRepository) {
				activeTenant := newActiveTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), id).
					Return(activeTenant, nil).
					Times(1)

				repo.EXPECT().
					UpdateStatus(
						gomock.Any(),
						id,
						tenantmodel.StatusDisabled,
						gomock.Any(),
					).
					Return(repositoryErr).
					Times(1)
			},
			wantErr: repositoryErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mocks.NewMockRepository(ctrl)

			tt.setupMock(repo)

			service := tenant.NewTenantService(repo)

			got, err := service.Disable(
				context.Background(),
				id,
			)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, &tenantmodel.Tenant{}, got)
				return
			}

			require.NoError(t, err)

			if tt.assert != nil {
				tt.assert(t, got)
			}
		})
	}
}
