package tenanthandler

import (
	"encoding/json"
	"net/http"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MembershipHandler struct {
	service *tenant.MembershipService
	logger  *zap.Logger
}

func NewMembershipHandler(service *tenant.MembershipService, logger *zap.Logger) *MembershipHandler {
	return &MembershipHandler{
		service: service,
		logger:  logger,
	}
}

func (mh *MembershipHandler) Create(w http.ResponseWriter, r *http.Request) {
	// tenantID, err := uuid.Parse(
	// 	chi.URLParam(r, "tenantID"),
	// )
	// if err != nil {
	// 	adapterutils.WriteJson(
	// 		w,
	// 		http.StatusBadRequest,
	// 		"invalid tenant id",
	// 	)
	// 	return
	// }

	var req httpmodel.CreateMembershipRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid user id")
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	membership, err := mh.service.Create(r.Context(), tenantID,
		tenantmodel.CreateMembershipInput{
			UserID: userID,
			Role:   tenantmodel.Role(req.Role),
		},
	)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusCreated, httpmodel.NewMembershipResponse(membership))
}

func (mh *MembershipHandler) Get(w http.ResponseWriter, r *http.Request) {
	var req httpmodel.GetMembershipRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			"invalid request body",
		)
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			"invalid tenant id",
		)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(
			w,
			http.StatusBadRequest,
			"invalid user id",
		)
		return
	}

	membership, err := mh.service.Get(r.Context(), tenantID, userID)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	adapterutils.WriteJson(
		w,
		http.StatusOK,
		httpmodel.NewMembershipResponse(membership),
	)
}

func (mh *MembershipHandler) List(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "tenantID"))
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid tenant id")
		return
	}

	memberships, err := mh.service.List(r.Context(), tenantID)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	response := make([]httpmodel.MembershipResponse, 0, len(memberships))

	for _, membership := range memberships {
		response = append(response, httpmodel.NewMembershipResponse(membership))
	}

	adapterutils.WriteJson(w, http.StatusOK, response)
}

func (mh *MembershipHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {

	var req httpmodel.UpdateMembershipRoleRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid tenant id")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	membership, err := mh.service.UpdateRole(
		r.Context(),
		tenantID,
		userID,
		tenantmodel.UpdateMembershipRoleInput{
			Role: tenantmodel.Role(req.Role),
		},
	)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.NewMembershipResponse(membership))
}

func (mh *MembershipHandler) Disable(w http.ResponseWriter, r *http.Request) {

	var req httpmodel.ChangeMembershipStatusRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid tenant id")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid user id")
		return
	}

	membership, err := mh.service.Disable(r.Context(), tenantID, userID)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.NewMembershipResponse(membership))
}

func (mh *MembershipHandler) Enable(w http.ResponseWriter, r *http.Request) {

	var req httpmodel.ChangeMembershipStatusRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid request body")
		return
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid tenant id")
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, "invalid user id")
		return
	}
	membership, err := mh.service.Enable(r.Context(), tenantID, userID)
	if err != nil {
		adapterutils.HandleMembershipError(w, err, mh.logger)
		return
	}

	adapterutils.WriteJson(
		w,
		http.StatusOK,
		httpmodel.NewMembershipResponse(membership),
	)
}
