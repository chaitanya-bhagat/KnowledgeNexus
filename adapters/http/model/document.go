package httpmodel

import (
	"time"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
)

type CreateDocumentRequest struct {
	TenantID     string `json:"tenantID"`
	KbID         string `json:"kbID"`
	Title        string `json:"title"`
	DocumentType string `json:"documentType"`
	CreatedBy    string `json:"createdBy"`
}

type GetDocumentRequest struct {
	TenantID   string `json:"tenantID"`
	DocumentID string `json:"documentID"`
}

type GetListRequest struct {
	TenantID string `json:"tenantID"`
	KbID     string `json:"kbID"`
}

type UpdateDocumentRequest struct {
	TenantID     string `json:"tenantID"`
	DocumentID   string `json:"documentID"`
	Title        string `json:"title"`
	DocumentType string `json:"documentType"`
}

type ChangeDocumentStatusRequest struct {
	TenantID   string `json:"tenantID"`
	DocumentID string `json:"documentID"`
}

// Response
type DocumentResponse struct {
	ID              string    `json:"id"`
	TenantID        string    `json:"tenantID"`
	KnowledgeBaseID string    `json:"knowledge_baseID"`
	Title           string    `json:"title"`
	DocumentType    string    `json:"documentType"`
	Status          string    `json:"status"`
	CreatedBy       string    `json:"createdBy"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func ToDocumentResponse(doc kbmodel.Document) DocumentResponse {
	return DocumentResponse{
		ID:              doc.ID.String(),
		TenantID:        doc.TenantID.String(),
		KnowledgeBaseID: doc.KbID.String(),
		Title:           doc.Title,
		DocumentType:    doc.DocumentType,
		Status:          string(doc.Status),
		CreatedBy:       doc.CreatedBy.String(),
		CreatedAt:       doc.CreatedAt,
		UpdatedAt:       doc.UpdatedAt,
	}
}

func ToDocumentListResponse(documents []kbmodel.Document) []DocumentResponse {
	response := make([]DocumentResponse, 0, len(documents))

	for _, doc := range documents {
		response = append(response, ToDocumentResponse(doc))
	}

	return response
}
