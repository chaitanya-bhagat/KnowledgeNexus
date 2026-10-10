package events

import "github.com/google/uuid"

const DocumentUploadedV1 = "knowledge.document.uploaded.v1"

type DocumentUploaded struct {
	EventID           uuid.UUID `json:"eventID"`
	DocumentID        uuid.UUID `json:"documentID"`
	TenantID          uuid.UUID `json:"tenantID"`
	DocumentVersionID uuid.UUID `json:"documentVersionID"`
	ObjectKey         string    `json:"objectKey"`
}
