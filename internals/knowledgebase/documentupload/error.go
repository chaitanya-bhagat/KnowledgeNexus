package documentupload

import "errors"

var (
	ErrInvalidDocumentID = errors.New("invalid document id")
	ErrDocumentArchived  = errors.New("cannot upload to archived document")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
	ErrInvalidUserID     = errors.New("invalid user id")

	ErrInvalidOriginalFileName = errors.New("invalid original file name")
	ErrInvalidSize             = errors.New("invalid upload size")

	ErrExpired                 = errors.New("document upload expired")
	ErrDocumentUploadNotFound  = errors.New("document upload not found")
	ErrDocumentUploadCompleted = errors.New("document upload already completed")
)
