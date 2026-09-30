package httpmodel

import (
	"time"
)

type InitiateRequest struct {
	TenantID   string `json:"tenantID"`
	CreatedBy  string `json:"createdBy"`
	DocumentID string `json:"documentID,omitempty"`

	Upload UploadRequest `json:"upload"`

	Document *DocumentRequest `json:"document,omitempty"`
}

type UploadRequest struct {
	OriginalFilename string `json:"originalFilename"`
	ContentType      string `json:"contentType"`
	SizeBytes        int64  `json:"sizeBytes"`
}

type DocumentRequest struct {
	KnowledgeBaseID string `json:"knowledgeBaseID"`
	Title           string `json:"title"`
	DocumentType    string `json:"documentType"`
}

type InitiateResponse struct {
	DocumentID string    `json:"documentID"`
	UploadID   string    `json:"uploadID"`
	URL        string    `json:"url"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

type CompleteRequest struct {
	TenantID string `json:"tenantID"`
	UploadID string `json:"uploadID"`
}

type CompleteResponse struct {
	DocumentID        string `json:"documentID"`
	DocumentVersionID string `json:"documentVersionID"`
	VersionNumber     int    `json:"versionNumber"`
}
