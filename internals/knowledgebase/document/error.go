package document

import "errors"

var (
	ErrInvalidDocumentID      = errors.New("invalid document id")
	ErrInvalidTenantID        = errors.New("invalid tenant id")
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrInvalidKnowledgeBaseID = errors.New("invalid knowledge base id")
	ErrInvalidTitle           = errors.New("invalid title")

	ErrDocumentNotFound      = errors.New("document not found")
	ErrKnowledgeBaseNotFound = errors.New("knowledge base not found")
	ErrKnowledgeBaseArchived = errors.New("knowledge base archived")

	ErrTenantNotFound = errors.New("tenant not found")
	ErrTenantDisabled = errors.New("tenant disabled")
	ErrUserNotFound   = errors.New("user not found")
)
