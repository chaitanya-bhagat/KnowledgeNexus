package identityhandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	"go.uber.org/zap"
)

type IdentityHandler struct {
	identityService *identity.IdentityService
	logger          *zap.Logger
}

func NewIdentityHandler(identityService *identity.IdentityService, logger *zap.Logger) *IdentityHandler {
	return &IdentityHandler{
		identityService: identityService,
		logger:          logger,
	}
}

func (ih *IdentityHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.CreateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	user, err := ih.identityService.CreateUser(r.Context(), &identity.User{
		DisplayName: req.DisplayName,
		Email:       req.Email,
	})
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}
	adapterutils.WriteJson(w, http.StatusCreated, httpmodel.ToUserResponse(*user))
}

func (ih *IdentityHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	user, err := ih.identityService.GetByID(r.Context(), userID)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToUserResponse(*user))
}

func (ih *IdentityHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "email query parameter is required",
		})
		return
	}

	user, err := ih.identityService.GetByEmail(r.Context(), email)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToUserResponse(*user))
}

func (ih *IdentityHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.UpdateUserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}
	user, err := ih.identityService.UpdatedUser(r.Context(), &identity.User{
		ID:          userID,
		DisplayName: req.DisplayName,
	})
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToUserResponse(*user))
}

func (ih *IdentityHandler) Disable(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.ChangeUserStatusRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}
	user, err := ih.identityService.Disable(r.Context(), userID, identity.StatusDisabled)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToUserResponse(*user))
}

func (ih *IdentityHandler) Enable(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.ChangeUserStatusRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	user, err := ih.identityService.Enable(r.Context(), userID, identity.StatusActive)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, ih.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToUserResponse(*user))
}
