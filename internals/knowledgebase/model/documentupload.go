package kbmodel

import (
	"time"

	"github.com/google/uuid"
)

type DocumentUploadStatus string

var (
	DocumentUploadPending   DocumentUploadStatus = "pending"
	DocumentUploadCompleted DocumentUploadStatus = "completed"
	DocumentUploadExpired   DocumentUploadStatus = "expired"
)

type DocumentUpload struct {
	ID uuid.UUID

	TenantID   uuid.UUID
	DocumentID uuid.UUID

	ObjectKey         string
	OriginalFileName  string
	ContentType       string
	ExpectedSizeBytes int64

	Status    DocumentUploadStatus
	CreatedBy uuid.UUID

	CreatedAt   time.Time
	ExpiredAt   time.Time
	CompletedAt *time.Time
}
