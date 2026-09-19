package kbmodel

import (
	"time"

	"github.com/google/uuid"
)

type DocumentVersionStatus string

var (
	DocumentVersionStatusUploaded   DocumentVersionStatus = "uploaded"
	DocumentVersionStatusProcessing DocumentVersionStatus = "processing"
	DocumentVersionStatusReady      DocumentVersionStatus = "ready"
	DocumentVersionStatusFailed     DocumentVersionStatus = "failed"
)

type DocumentVersion struct {
	ID               uuid.UUID
	TenantID         uuid.UUID
	DocumentID       uuid.UUID
	VersionNumber    int
	ObjectKey        string
	OriginalFileName string
	ContentType      string
	SizeBytes        int64
	Checksum         string
	Status           DocumentVersionStatus
	FailureReason    string
	CreatedBy        uuid.UUID
	CreatedAt        time.Time
}

type DocumentVersionInput struct {
	TenantID         uuid.UUID
	DocumentID       uuid.UUID
	ObjectKey        string
	OriginalFileName string
	ContentType      string
	SizeBytes        int64
	Checksum         string
	CreatedBy        uuid.UUID
}
