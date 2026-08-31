package tenant

import "errors"

var (
	ErrNotFound     = errors.New("tenant not found")
	ErrInvalidName  = errors.New("tenant name is required")
	ErrInvalidSlug  = errors.New("tenant slug is required")
	ErrSlugConflict = errors.New("tenant slug already exists")
)
