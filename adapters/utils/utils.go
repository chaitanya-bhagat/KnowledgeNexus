package adapterutils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	"go.uber.org/zap"
)

func WriteJson(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func HandleMembershipError(w http.ResponseWriter, err error, logger *zap.Logger) {
	switch {
	case errors.Is(err, tenant.ErrInvalidTenantID):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrInvalidUserID):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrInvalidRole):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrOwnerRoleManagedSeparately):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrMembershipExists):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrTenantDisabled):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrNotFound), errors.Is(err, tenant.ErrUserNotFound), errors.Is(err, tenant.ErrMembershipNotFound):
		WriteJson(
			w,
			http.StatusNotFound,
			err.Error(),
		)

	default:
		logger.Error(
			"membership request failed",
			zap.Error(err),
		)

		WriteJson(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
	}
}
