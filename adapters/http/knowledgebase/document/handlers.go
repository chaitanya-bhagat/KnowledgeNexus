package documenthandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"go.uber.org/zap"
)

type DocumentHandler struct {
	documentService *document.DocumentService
	logger          *zap.Logger
}

func NewDocumentHandler(documentService *document.DocumentService, logger *zap.Logger) *DocumentHandler {
	return &DocumentHandler{
		documentService: documentService,
		logger:          logger,
	}
}

func (dh *DocumentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.CreateDocumentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	kbID, err := uuid.Parse(input.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledge base id",
		})
		return
	}

	userID, err := uuid.Parse(input.CreatedBy)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	document, err := dh.documentService.Create(r.Context(), kbmodel.CreateDocument{
		TenantID:        tenantID,
		KnowledgeBaseID: kbID,
		Title:           input.Title,
		DocumentType:    input.DocumentType,
		CreatedBy:       userID,
	})

	if err != nil {
		dh.logger.Error("document creation failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"document creation failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusCreated, httpmodel.ToDocumentResponse(document))
}

func (dh *DocumentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.GetDocumentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	documentID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document id",
		})
		return
	}
	document, err := dh.documentService.GetByID(r.Context(), tenantID, documentID)
	if err != nil {
		dh.logger.Error("get document failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"get document failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentResponse(document))
}

func (dh *DocumentHandler) GetList(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.GetListRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	kbID, err := uuid.Parse(input.KbID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledge base id",
		})
		return
	}
	documents, err := dh.documentService.ListByKnowledgeBase(r.Context(), kbID, tenantID)
	if err != nil {
		dh.logger.Error("get list of documents failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"get list of documents failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentListResponse(documents))
}

func (dh *DocumentHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.UpdateDocumentRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	documentID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document id",
		})
		return
	}
	document, err := dh.documentService.Update(r.Context(), tenantID, documentID, kbmodel.UpdateDocumentInput{
		Title:        input.Title,
		DocumentType: input.DocumentType,
	})
	if err != nil {
		dh.logger.Error("failed to update document", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"led to update document": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentResponse(document))
}

func (dh *DocumentHandler) Archive(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.ChangeDocumentStatusRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	documentID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document id",
		})
		return
	}
	document, err := dh.documentService.Archive(r.Context(), tenantID, documentID)
	if err != nil {
		dh.logger.Error("failed to archive document", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"led to archive document": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentResponse(document))
}

func (dh *DocumentHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.ChangeDocumentStatusRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}
	tenantID, err := uuid.Parse(input.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	documentID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document id",
		})
		return
	}
	document, err := dh.documentService.Activate(r.Context(), tenantID, documentID)
	if err != nil {
		dh.logger.Error("failed to activate document", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"led to activate document": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentResponse(document))
}
