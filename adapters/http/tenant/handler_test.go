package tenanthandler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	tenanthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/mocks"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

func newTestHandler(t *testing.T) (*tenanthandler.TenantHandler, *mocks.MockRepository) {
	t.Helper()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepository(ctrl)
	service := tenant.NewTenantService(repo)

	handler := tenanthandler.NewTenantHandler(
		service,
		zap.NewNop(),
	)

	return handler, repo
}

func newTestRouter(handler *tenanthandler.TenantHandler) http.Handler {
	r := chi.NewRouter()

	r.Post("/tenants", handler.Create)
	r.Get("/tenants/{id}", handler.GetByID)
	r.Patch("/tenants/{id}", handler.Update)
	r.Post("/tenants/{id}/disable", handler.Disable)

	return r
}

func TestHandler_Create(t *testing.T) {
	repositoryErr := errors.New("database error")

	tests := []struct {
		name       string
		body       string
		setupMock  func(*mocks.MockRepository)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			body: `{
				"name": "Acme Legal",
				"slug": "acme-legal"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantStatus: http.StatusCreated,
			wantBody:   `"status":"active"`,
		},
		{
			name:       "invalid json",
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid request body"`,
		},
		{
			name: "unknown field",
			body: `{
				"name": "Acme",
				"slug": "acme",
				"unknown": "value"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid request body"`,
		},
		{
			name: "invalid name",
			body: `{
				"name": "",
				"slug": "acme"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"tenant name is required"`,
		},
		{
			name: "invalid slug",
			body: `{
				"name": "Acme",
				"slug": ""
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"tenant slug is required"`,
		},
		{
			name: "slug conflict",
			body: `{
				"name": "Acme",
				"slug": "acme"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(tenant.ErrSlugConflict).
					Times(1)
			},
			wantStatus: http.StatusConflict,
			wantBody:   `"error":"tenant slug already exists"`,
		},
		{
			name: "internal error",
			body: `{
				"name": "Acme",
				"slug": "acme"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(repositoryErr).
					Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"error":"internal server error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, repo := newTestHandler(t)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			req := httptest.NewRequest(
				http.MethodPost,
				"/tenants",
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			handler.Create(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
			require.Contains(
				t,
				rec.Header().Get("Content-Type"),
				"application/json",
			)
		})
	}
}

func TestHandler_GetByID(t *testing.T) {
	validID := uuid.New()
	repositoryErr := errors.New("database error")

	tests := []struct {
		name       string
		id         string
		setupMock  func(*mocks.MockRepository)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(
						&tenantmodel.Tenant{
							ID:     validID,
							Name:   "Acme",
							Slug:   "acme",
							Status: tenantmodel.StatusActive,
						},
						nil,
					).
					Times(1)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"name":"Acme"`,
		},
		{
			name:       "invalid uuid",
			id:         "not-a-uuid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid tenant id"`,
		},
		{
			name: "not found",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `"error":"tenant not found"`,
		},
		{
			name: "internal error",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(
						&tenantmodel.Tenant{},
						repositoryErr,
					).
					Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"error":"internal server error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, repo := newTestHandler(t)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			router := newTestRouter(handler)

			req := httptest.NewRequest(
				http.MethodGet,
				"/tenants/"+tt.id,
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}

func TestHandler_Update(t *testing.T) {
	validID := uuid.New()
	repositoryErr := errors.New("database error")

	existing := &tenantmodel.Tenant{
		ID:     validID,
		Name:   "Old",
		Slug:   "old",
		Status: tenantmodel.StatusActive,
	}

	tests := []struct {
		name       string
		id         string
		body       string
		setupMock  func(*mocks.MockRepository)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			id:   validID.String(),
			body: `{
				"name": "New Acme",
				"slug": "new-acme"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"name":"New Acme"`,
		},
		{
			name: "invalid uuid",
			id:   "invalid",
			body: `{
				"name": "Acme",
				"slug": "acme"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid tenant id"`,
		},
		{
			name:       "invalid json",
			id:         validID.String(),
			body:       `{"name":`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid request body"`,
		},
		{
			name: "invalid name",
			id:   validID.String(),
			body: `{
				"name": "",
				"slug": "acme"
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"tenant name is required"`,
		},
		{
			name: "invalid slug",
			id:   validID.String(),
			body: `{
				"name": "Acme",
				"slug": ""
			}`,
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"tenant slug is required"`,
		},
		{
			name: "not found",
			id:   validID.String(),
			body: `{
				"name": "Acme",
				"slug": "acme"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `"error":"tenant not found"`,
		},
		{
			name: "slug conflict",
			id:   validID.String(),
			body: `{
				"name": "Acme",
				"slug": "existing"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(tenant.ErrSlugConflict).
					Times(1)
			},
			wantStatus: http.StatusConflict,
			wantBody:   `"error":"tenant slug already exists"`,
		},
		{
			name: "internal error",
			id:   validID.String(),
			body: `{
				"name": "Acme",
				"slug": "acme"
			}`,
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(existing, nil).
					Times(1)

				repo.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(repositoryErr).
					Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"error":"internal server error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, repo := newTestHandler(t)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			router := newTestRouter(handler)

			req := httptest.NewRequest(
				http.MethodPatch,
				"/tenants/"+tt.id,
				bytes.NewBufferString(tt.body),
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}

func TestHandler_Disable(t *testing.T) {
	validID := uuid.New()
	repositoryErr := errors.New("database error")

	newActiveTenant := func() *tenantmodel.Tenant {
		return &tenantmodel.Tenant{
			ID:     validID,
			Name:   "Acme",
			Slug:   "acme",
			Status: tenantmodel.StatusActive,
		}
	}

	newDisabledTenant := func() *tenantmodel.Tenant {
		return &tenantmodel.Tenant{
			ID:     validID,
			Name:   "Acme",
			Slug:   "acme",
			Status: tenantmodel.StatusDisabled,
		}
	}

	tests := []struct {
		name       string
		id         string
		setupMock  func(*mocks.MockRepository)
		wantStatus int
		wantBody   string
	}{
		{
			name: "success",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				activeTenant := newActiveTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(activeTenant, nil).
					Times(1)

				repo.EXPECT().
					UpdateStatus(
						gomock.Any(),
						validID,
						tenantmodel.StatusDisabled,
						gomock.Any(),
					).
					Return(nil).
					Times(1)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":"disabled"`,
		},
		{
			name: "already disabled",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				disabledTenant := newDisabledTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(disabledTenant, nil).
					Times(1)
			},
			wantStatus: http.StatusOK,
			wantBody:   `"status":"disabled"`,
		},
		{
			name:       "invalid uuid",
			id:         "invalid",
			wantStatus: http.StatusBadRequest,
			wantBody:   `"error":"invalid tenant id"`,
		},
		{
			name: "not found",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(
						&tenantmodel.Tenant{},
						tenant.ErrNotFound,
					).
					Times(1)
			},
			wantStatus: http.StatusNotFound,
			wantBody:   `"error":"tenant not found"`,
		},
		{
			name: "internal error",
			id:   validID.String(),
			setupMock: func(repo *mocks.MockRepository) {
				activeTenant := newActiveTenant()
				repo.EXPECT().
					GetByID(gomock.Any(), validID).
					Return(activeTenant, nil).
					Times(1)

				repo.EXPECT().
					UpdateStatus(
						gomock.Any(),
						validID,
						tenantmodel.StatusDisabled,
						gomock.Any(),
					).
					Return(repositoryErr).
					Times(1)
			},
			wantStatus: http.StatusInternalServerError,
			wantBody:   `"error":"internal server error"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, repo := newTestHandler(t)

			if tt.setupMock != nil {
				tt.setupMock(repo)
			}

			router := newTestRouter(handler)

			req := httptest.NewRequest(
				http.MethodPost,
				"/tenants/"+tt.id+"/disable",
				nil,
			)

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, tt.wantStatus, rec.Code)
			require.Contains(t, rec.Body.String(), tt.wantBody)
		})
	}
}
