package documentversion

import "errors"

var (
	ErrInvalidVersionID  = errors.New("invalid document version id")
	ErrInvalidDocumentID = errors.New("invalid document id")
	ErrInvalidTenantID   = errors.New("invalid tenant id")
	ErrInvalidUserID     = errors.New("invalid user id")

	ErrInvalidObjectKey        = errors.New("invalid object key")
	ErrInvalidOriginalFileName = errors.New("invalid original file name")
	ErrInvalidSize             = errors.New("invalid document size")
	ErrInvalidStatus           = errors.New("invalid document version status")

	ErrVersionNotFound  = errors.New("version not found")
	ErrDocumentNotFound = errors.New("document not found")
	ErrDocumentArchived = errors.New("document archived")

	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
	ErrKnowledgeBaseArchived = errors.New("knowledge base archived")

	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantDisabled = errors.New("tenant disabled")
	ErrUserNotFound   = errors.New("user not found")
)
