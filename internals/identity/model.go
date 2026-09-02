package identity

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
)

type User struct {
	ID          uuid.UUID
	DisplayName string
	Email       string
	Status      Status
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreatedUser struct {
	Name  string
	Email string
}

type UpdateUser struct {
	Name   string
	Status Status
	Email  string
}
