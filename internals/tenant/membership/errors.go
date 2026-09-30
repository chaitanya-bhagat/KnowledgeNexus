package membership

import "errors"

// Membership errors
var (
	ErrInvalidTenantID = errors.New("invalid tenant id")
	ErrInvalidUserID   = errors.New("invalid user id")

	ErrInvalidRole = errors.New("invalid membership role")

	ErrMembershipNotFound = errors.New("membership not found")
	ErrMembershipExists   = errors.New("membership already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrTenantDisabled = errors.New("tenant is disabled")

	ErrOwnerRoleManagedSeparately = errors.New("owner membership must be managed through the owner workflow")
)
