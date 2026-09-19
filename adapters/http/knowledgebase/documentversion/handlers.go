package documentversionhandler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	httpmodel "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/model"
	adapterutils "github.com/chaitanya-bhagat/knowledge-nexus/adapters/utils"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentversion"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"go.uber.org/zap"
)

type DocumentVersionHandler struct {
	documentVersionService *documentversion.DocVersionService
	logger                 *zap.Logger
}

func NewDocumentVersionHandler(documentVersionService *documentversion.DocVersionService, logger *zap.Logger) *DocumentVersionHandler {
	return &DocumentVersionHandler{
		documentVersionService: documentVersionService,
		logger:                 logger,
	}
}

func (dvh *DocumentVersionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.CreateDocumentVersion

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

	docID, err := uuid.Parse(input.DocumentID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document id",
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
	docVersion, err := dvh.documentVersionService.Create(r.Context(), kbmodel.DocumentVersionInput{
		TenantID:         tenantID,
		DocumentID:       docID,
		ObjectKey:        input.ObjectKey,
		OriginalFileName: input.OriginalFilename,
		ContentType:      input.ContentType,
		Checksum:         input.Checksum,
		SizeBytes:        input.SizeBytes,
		CreatedBy:        userID,
	})
	if err != nil {
		dvh.logger.Error("document version creation failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"document version creation failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusCreated, httpmodel.ToDocumentVersionResponse(docVersion))
}

func (dvh *DocumentVersionHandler) Get(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.GetDocumentVersion

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

	versionID, err := uuid.Parse(input.VersionID)
	if err != nil {
		adapterutils.WriteJson(w, http.StatusBadRequest, map[string]string{
			"error": "invalid document version id",
		})
		return
	}
	documentVersion, err := dvh.documentVersionService.GetByID(r.Context(), tenantID, versionID)
	if err != nil {
		dvh.logger.Error("get document version failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"get document version failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentVersionResponse(documentVersion))

}

func (dvh *DocumentVersionHandler) List(w http.ResponseWriter, r *http.Request) {
	var input httpmodel.ListDocumentVersions

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
			"error": "invalid knowledge base id",
		})
		return
	}
	documentVersions, err := dvh.documentVersionService.GetListByDocumentID(r.Context(), tenantID, documentID)
	if err != nil {
		dvh.logger.Error("get list of document version failed", zap.Error(err))
		adapterutils.WriteJson(w, http.StatusInternalServerError, map[string]string{"get list of documents failed": err.Error()})
		return
	}
	adapterutils.WriteJson(w, http.StatusOK, httpmodel.ToDocumentVersionListResponse(documentVersions))
}
