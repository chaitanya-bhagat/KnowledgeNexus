package httpmodel

import (
	"time"

	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
)

type CreateDocumentVersion struct {
	TenantID         string `json:"tenantID"`
	DocumentID       string `json:"documentID"`
	ObjectKey        string `json:"objectKey"`
	OriginalFilename string `json:"originalFilename"`
	ContentType      string `json:"contentType"`
	SizeBytes        int64  `json:"sizeBytes"`
	Checksum         string `json:"checksum"`
	CreatedBy        string `json:"createdBy"`
}

type GetDocumentVersion struct {
	TenantID  string `json:"tenantID"`
	VersionID string `json:"versionID"`
}

type ListDocumentVersions struct {
	TenantID   string `json:"tenantID"`
	DocumentID string `json:"documentID"`
}

//Response

type DocumentVersion struct {
	ID               string    `json:"id"`
	TenantID         string    `json:"tenantID"`
	DocumentID       string    `json:"documentID"`
	VersionNumber    int       `json:"versionNumber"`
	ObjectKey        string    `json:"objectKey"`
	OriginalFilename string    `json:"originalFilename"`
	ContentType      string    `json:"contentType"`
	SizeBytes        int64     `json:"sizeBytes"`
	Checksum         string    `json:"checksum"`
	Status           string    `json:"status"`
	FailureReason    string    `json:"failureReason"`
	CreatedBy        string    `json:"createdBy"`
	CreatedAt        time.Time `json:"createdAt"`
}

func ToDocumentVersionResponse(version kbmodel.DocumentVersion) DocumentVersion {
	return DocumentVersion{
		ID:               version.ID.String(),
		TenantID:         version.TenantID.String(),
		DocumentID:       version.DocumentID.String(),
		VersionNumber:    version.VersionNumber,
		ObjectKey:        version.ObjectKey,
		OriginalFilename: version.OriginalFileName,
		ContentType:      version.ContentType,
		SizeBytes:        version.SizeBytes,
		Checksum:         version.Checksum,
		Status:           string(version.Status),
		FailureReason:    version.FailureReason,
		CreatedBy:        version.CreatedBy.String(),
		CreatedAt:        version.CreatedAt,
	}
}

func ToDocumentVersionListResponse(versions []kbmodel.DocumentVersion) []DocumentVersion {
	response := make([]DocumentVersion, 0, len(versions))

	for _, version := range versions {
		response = append(
			response,
			ToDocumentVersionResponse(version),
		)
	}

	return response
}
