package adapterutils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/membership"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant"
	"go.uber.org/zap"
)

func WriteJson(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func HandleMembershipError(w http.ResponseWriter, err error, logger *zap.Logger) {
	switch {
	case errors.Is(err, membership.ErrInvalidTenantID):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, membership.ErrInvalidUserID):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, membership.ErrInvalidRole):
		WriteJson(
			w,
			http.StatusBadRequest,
			err.Error(),
		)

	case errors.Is(err, membership.ErrOwnerRoleManagedSeparately):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, membership.ErrMembershipExists):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, membership.ErrTenantDisabled):
		WriteJson(
			w,
			http.StatusConflict,
			err.Error(),
		)

	case errors.Is(err, tenant.ErrNotFound), errors.Is(err, membership.ErrUserNotFound), errors.Is(err, membership.ErrMembershipNotFound):
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
