package knowledgebasehandler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase"
	"go.uber.org/zap"
)

type KnowledgeBaseHandler struct {
	kbService *knowledgebase.KnowledgeBaseService
	logger    *zap.Logger
}

func NewKnowledgeBasehandler(kbService knowledgebase.KnowledgeBaseService, logger *zap.Logger) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{
		kbService: &kbService,
		logger:    logger,
	}
}

func (kbh *KnowledgeBaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var kb httpmodel.CreateKnowledgeBase

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&kb); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	tenantID, err := uuid.Parse(kb.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}
	userID, err := uuid.Parse(kb.CreatedBy)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}
	newKb, err := kbh.kbService.CreateKnowledgeBase(r.Context(), knowledgebase.KnowledgeBase{
		Name:        kb.Name,
		TenantID:    tenantID,
		CreatedBy:   userID,
		DomainType:  kb.DomainType,
		Description: kb.Description,
	})
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}

	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseResponse(newKb))
}

func (kbh *KnowledgeBaseHandler) GetKnowledgeBase(w http.ResponseWriter, r *http.Request) {
	var kb httpmodel.GetKnowledgeBase
	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	if err := decode.Decode(&kb); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	kbID, err := uuid.Parse(kb.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledgebase id",
		})
		return
	}
	tenantID, err := uuid.Parse(kb.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	kbDetail, err := kbh.kbService.GetKnowledgeBaseByID(r.Context(), kbID, tenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseResponse(kbDetail))

}

func (kbh *KnowledgeBaseHandler) GetList(w http.ResponseWriter, r *http.Request) {
	tenantID, err := uuid.Parse(chi.URLParam(r, "tenantID"))
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	kbs, err := kbh.kbService.ListKnowledgeBasesByTenantID(r.Context(), tenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseListResponse(kbs))
}

func (kbh *KnowledgeBaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var kb httpmodel.UpdateKnowledgeBase

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	if err := decode.Decode(&kb); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	kbID, err := uuid.Parse(kb.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledgebase id",
		})
		return
	}
	tenantID, err := uuid.Parse(kb.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	kbDetails, err := kbh.kbService.UpdateKnowledgeBase(r.Context(), kbID, tenantID, &knowledgebase.UpdateKnowledgeBase{
		Name:        kb.Name,
		Description: kb.Description,
	})
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseResponse(kbDetails))
}

func (kbh *KnowledgeBaseHandler) Archive(w http.ResponseWriter, r *http.Request) {
	var kb httpmodel.ChangeKnowledgeBaseStatus

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	if err := decode.Decode(&kb); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	kbID, err := uuid.Parse(kb.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledgebase id",
		})
		return
	}
	tenantID, err := uuid.Parse(kb.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}
	kbDetails, err := kbh.kbService.Archive(r.Context(), tenantID, kbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseResponse(kbDetails))

}

func (kbh *KnowledgeBaseHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var kb httpmodel.ChangeKnowledgeBaseStatus

	decode := json.NewDecoder(r.Body)
	decode.DisallowUnknownFields()
	if err := decode.Decode(&kb); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	kbID, err := uuid.Parse(kb.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledgebase id",
		})
		return
	}
	tenantID, err := uuid.Parse(kb.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}
	kbDetails, err := kbh.kbService.Activate(r.Context(), tenantID, kbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusInternalServerError, err)
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToKnowledgeBaseResponse(kbDetails))

}
