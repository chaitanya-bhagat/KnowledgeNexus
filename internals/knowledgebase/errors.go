package knowledgebase

import "errors"

var (
	ErrKnowledgeBaseNotFound  = errors.New("knowledge base not found")
	ErrInvalidKnowledgeBase   = errors.New("invalid knowledge base name")
	ErrInvalidKnowledgeBaseID = errors.New("invalid knowledge base ID")
	ErrKnowledgeBaseConflict  = errors.New("knowledge base name already exists")
	ErrKnowledgeBaseDisabled  = errors.New("knowledge base is disabled")

	ErrInvalidTenantID = errors.New("invalid tenant ID")
	ErrInvalidUserID   = errors.New("invalid user ID")
	ErrTenantDisabled  = errors.New("tenant is disabled")
	ErrUserNotFound    = errors.New("user not found")
	ErrTenantNotFound  = errors.New("tenant not found")
)
