package tenanthandler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type TenantHandler struct {
	service *tenant.TenantService
	logger  *zap.Logger
}

func NewTenantHandler(service *tenant.TenantService, logger *zap.Logger) *TenantHandler {
	return &TenantHandler{
		service: service,
		logger:  logger,
	}
}

func (th *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.CreateTenantRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	input := tenantmodel.CreateInput{
		Name: req.Name,
		Slug: req.Slug,
	}
	t, err := th.service.Create(r.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, tenant.ErrInvalidName):
			adapterutils.WriteJson(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
		case errors.Is(err, tenant.ErrInvalidSlug):
			adapterutils.WriteJson(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)

		case errors.Is(err, tenant.ErrSlugConflict):
			adapterutils.WriteJson(
				w,
				http.StatusConflict,
				map[string]string{
					"error": err.Error(),
				},
			)
		default:
			th.logger.Error("failed to create tenant", zap.Error(err))
			adapterutils.WriteJson(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": "internal server error",
				},
			)

		}
		return
	}
	adapterutils.WriteJson(w, http.StatusCreated, httpmodel.NewTenantResponse(t))
}

func (th *TenantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid tenant id",
			},
		)
		return
	}

	t, err := th.service.GetByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, tenant.ErrNotFound):
			adapterutils.WriteJson(
				w,
				http.StatusNotFound,
				map[string]string{
					"error": err.Error(),
				},
			)

		default:
			th.logger.Error(
				"failed to get tenant",
				zap.String("tenant_id", id.String()),
				zap.Error(err),
			)

			adapterutils.WriteJson(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": "internal server error",
				},
			)
		}

		return
	}

	adapterutils.WriteJson(
		w,
		http.StatusOK,
		httpmodel.NewTenantResponse(t),
	)
}
func (th *TenantHandler) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid tenant id",
			},
		)
		return
	}

	var req httpmodel.UpdateTenantRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request body",
			},
		)
		return
	}

	input := tenantmodel.UpdateInput{
		Name: req.Name,
		Slug: req.Slug,
	}

	t, err := th.service.Update(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, tenant.ErrInvalidName),
			errors.Is(err, tenant.ErrInvalidSlug):

			adapterutils.WriteJson(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)

		case errors.Is(err, tenant.ErrNotFound):

			adapterutils.WriteJson(
				w,
				http.StatusNotFound,
				map[string]string{
					"error": err.Error(),
				},
			)

		case errors.Is(err, tenant.ErrSlugConflict):

			adapterutils.WriteJson(
				w,
				http.StatusConflict,
				map[string]string{
					"error": err.Error(),
				},
			)

		default:
			th.logger.Error(
				"failed to update tenant",
				zap.String("tenant_id", id.String()),
				zap.Error(err),
			)

			adapterutils.WriteJson(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": "internal server error",
				},
			)
		}

		return
	}

	adapterutils.WriteJson(
		w,
		http.StatusOK,
		httpmodel.NewTenantResponse(t),
	)
}

func (th *TenantHandler) Disable(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid tenant id",
			},
		)
		return
	}

	t, err := th.service.Disable(
		r.Context(),
		id,
	)
	if err != nil {
		switch {
		case errors.Is(err, tenant.ErrNotFound):
			adapterutils.WriteJson(
				w,
				http.StatusNotFound,
				map[string]string{
					"error": err.Error(),
				},
			)

		default:
			th.logger.Error(
				"failed to disable tenant",
				zap.String("tenant_id", id.String()),
				zap.Error(err),
			)

			adapterutils.WriteJson(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": "internal server error",
				},
			)
		}

		return
	}

	adapterutils.WriteJson(
		w,
		http.StatusOK,
		httpmodel.NewTenantResponse(t),
	)
}
