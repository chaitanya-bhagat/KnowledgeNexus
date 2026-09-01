package tenant

import "errors"

var (
	ErrNotFound     = errors.New("tenant not found")
	ErrInvalidName  = errors.New("tenant name is required")
	ErrInvalidSlug  = errors.New("tenant slug is required")
	ErrSlugConflict = errors.New("tenant slug already exists")
)

// Membership errors
var (
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidUserID   = errors.New("invalid user id")

	ErrInvalidRole = errors.New("invalid membership role")

	ErrMembershipNotFound = errors.New("membership not found")
	ErrMembershipExists   = errors.New("membership already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrTenantDisabled = errors.New("tenant is disabled")

	ErrOwnerRoleManagedSeparately = errors.New(
		"owner membership must be managed through the owner workflow",
	)
)
