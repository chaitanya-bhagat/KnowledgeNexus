package identity

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid userID")
	ErrNotFound      = errors.New("user not found")
	ErrInvalidName   = errors.New("invalid name")
	ErrInvalidEmail  = errors.New("invalid email")
	ErrEmailConflict = errors.New("email already exists")
)
