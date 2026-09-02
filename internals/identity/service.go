package identity

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"

	"go.uber.org/zap"
)

type IdentityService struct {
	repo   Repository
	logger *zap.Logger
}

func NewIdentityService(repo Repository, Logger *zap.Logger) *IdentityService {
	return &IdentityService{
		repo:   repo,
		logger: Logger,
	}
}

func (is *IdentityService) CreateUser(ctx context.Context, user *User) (*User, error) {
	email := strings.ToLower(strings.TrimSpace(user.Email))
	name := strings.TrimSpace(user.DisplayName)
	if email == "" {
		return &User{}, ErrInvalidEmail
	}
	if name == "" {
		return &User{}, ErrInvalidName
	}

	time := time.Now().UTC()

	user = &User{
		ID:          uuid.New(),
		DisplayName: name,
		Email:       email,
		Status:      StatusActive,
		CreatedAt:   time,
		UpdatedAt:   time,
	}
	user, err := is.repo.Create(ctx, user)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (is *IdentityService) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	if id == uuid.Nil {
		return &User{}, ErrInvalidUserID
	}
	return is.repo.GetByID(ctx, id)
}

func (is *IdentityService) GetByEmail(ctx context.Context, email string) (*User, error) {
	if email == "" {
		return &User{}, ErrInvalidEmail
	}
	return is.repo.GetByEmail(ctx, email)
}

func (is *IdentityService) UpdatedUser(ctx context.Context, user *User) (*User, error) {
	if user.ID == uuid.Nil {
		return &User{}, ErrInvalidUserID
	}
	if user.DisplayName == "" {
		return &User{}, ErrInvalidName
	}

	user, err := is.repo.GetByID(ctx, user.ID)
	if err != nil {
		return &User{}, err
	}

	user.DisplayName = strings.TrimSpace(user.DisplayName)
	user.UpdatedAt = time.Now().UTC()
	err = is.repo.UpdateUser(ctx, user)
	if err != nil {
		return &User{}, err
	}
	return user, nil

}

func (is *IdentityService) Disable(ctx context.Context, id uuid.UUID, status Status) (*User, error) {
	if id == uuid.Nil {
		return &User{}, ErrInvalidUserID
	}
	user, err := is.repo.GetByID(ctx, id)
	if err != nil {
		return &User{}, err
	}
	if user.Status == StatusDisabled {
		return user, nil
	}
	user.Status = status
	user.UpdatedAt = time.Now().UTC()
	err = is.repo.UpdateStatus(ctx, id, status, user.UpdatedAt)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func (is *IdentityService) Enable(ctx context.Context, id uuid.UUID, status Status) (*User, error) {
	if id == uuid.Nil {
		return &User{}, ErrInvalidUserID
	}
	user, err := is.repo.GetByID(ctx, id)
	if err != nil {
		return &User{}, err
	}
	if user.Status == StatusActive {
		return user, nil
	}

	user.Status = status
	user.UpdatedAt = time.Now().UTC()
	err = is.repo.UpdateStatus(ctx, id, status, user.UpdatedAt)

	if err != nil {
		return &User{}, err
	}

	return user, nil
}
