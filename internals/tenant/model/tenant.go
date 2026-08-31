package tenantmodel

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusDisabled Status = "disabled"
)

type Tenant struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateInput struct {
	Name string
	Slug string
}

type UpdateInput struct {
	Name string
	Slug string
}
