package documentuploadhandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentupload"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"go.uber.org/zap"
)

type DocumentUploadHandler struct {
	documentUploadService *documentupload.DocumentUploadService
	logger                *zap.Logger
}

func NewDocumentUploadHandler(documentUploadService *documentupload.DocumentUploadService, logger *zap.Logger) *DocumentUploadHandler {
	return &DocumentUploadHandler{
		documentUploadService: documentUploadService,
		logger:                logger,
	}
}

func (duh *DocumentUploadHandler) Complete(w http.ResponseWriter, r *http.Request) {
	var request httpmodel.CompleteRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	tenantID, err := uuid.Parse(request.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}

	uploadID, err := uuid.Parse(request.UploadID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid upload id",
		})
		return
	}
	result, err := duh.documentUploadService.Complete(r.Context(), kbmodel.CompleteInput{
		TenantID: tenantID,
		UploadID: uploadID,
	})
	if err != nil {
		duh.logger.Error("document upload complete failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"document upload complete failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusCreated, result)

}

func (duh *DocumentUploadHandler) Initiate(w http.ResponseWriter, r *http.Request) {
	var request httpmodel.InitiateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
		return
	}

	tenantID, err := uuid.Parse(request.TenantID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid tenant id",
		})
		return
	}
	var documentID *uuid.UUID
	if request.DocumentID != "" {
		parsedDocumentID, parseErr := uuid.Parse(request.DocumentID)
		err = parseErr
		if err != nil {
			adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
				"error": "invalid document id",
			})
			return
		}
		documentID = &parsedDocumentID
	}

	userID, err := uuid.Parse(request.CreatedBy)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid user id",
		})
		return
	}

	kbID, err := uuid.Parse(request.Document.KnowledgeBaseID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid knowledge base id",
		})
		return
	}

	upload, err := duh.documentUploadService.Initiate(r.Context(), kbmodel.InitiateInput{
		TenantID:   tenantID,
		DocumentID: documentID,
		CreatedBy:  userID,
		Upload:     kbmodel.UploadInput(request.Upload),
		Document: &kbmodel.DocumentInput{
			KnowledgeBaseID: kbID,
			Title:           request.Document.Title,
			DocumentType:    request.Document.DocumentType,
		},
	})
	if err != nil {
		duh.logger.Error("failed to initiate upload", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"failed to initiate upload": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, upload)

}
