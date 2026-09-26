package documentupload

import "errors"

var (
	ErrInvalidDocumentID       = errors.New("invalid document id")
	ErrInvalidDocumentMetadata = errors.New("document metadata required")
	ErrDocumentArchived        = errors.New("cannot upload to archived document")
	ErrInvalidTenantID         = errors.New("invalid tenant id")
	ErrInvalidUserID           = errors.New("invalid user id")
	ErrInvalidKBID             = errors.New("invalid knowledge base id")
	ErrInvalidTitle            = errors.New("invalid title")
	ErrInvalidDocumentType     = errors.New("invalid document type")
	ErrInActiveTenant          = errors.New("tenant is deactivated")
	ErrKBArchived              = errors.New("knowledge base is arcchived")
	ErrUsersMembershipDisabled = errors.New("users membership is disabled")

	ErrInvalidOriginalFileName = errors.New("invalid original file name")
	ErrInvalidSize             = errors.New("invalid upload size")

	ErrExpired                 = errors.New("document upload expired")
	ErrDocumentUploadNotFound  = errors.New("document upload not found")
	ErrDocumentUploadCompleted = errors.New("document upload already completed")
)
