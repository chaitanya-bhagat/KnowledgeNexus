package adapteridentity

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type identityRepository struct {
	db *pgxpool.Pool
}

func NewIdentityRepository(db *pgxpool.Pool) *identityRepository {
	return &identityRepository{
		db: db,
	}
}

func (ir *identityRepository) Create(ctx context.Context, u *identity.User) (*identity.User, error) {
	const query = `INSERT INTO table_users (id,	display_name, email, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, display_name, email, status, created_at, updated_at`

	_, err := ir.db.Exec(ctx, query, u.ID, u.DisplayName, u.Email, u.Status, u.CreatedAt, u.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_lower_unique" {
				return &identity.User{}, identity.ErrEmailConflict
			}
		}
		return &identity.User{}, err
	}
	return u, nil

}

func (ir *identityRepository) GetByID(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	const query = `SELECT id, display_name, email, status, created_at, updated_at FROM table_users WHERE id = $1`
	var u identity.User

	err := ir.db.QueryRow(ctx, query, id).Scan(&u.ID, &u.DisplayName, &u.Email, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &identity.User{}, identity.ErrNotFound
		}
		return &identity.User{}, err
	}
	return &u, nil
}

func (ir *identityRepository) GetByEmail(ctx context.Context, email string) (*identity.User, error) {
	const query = `SELECT id, display_name, email, status, created_at, updated_at FROM table_users	WHERE LOWER(email) = LOWER($1)`
	var u identity.User
	err := ir.db.QueryRow(ctx, query, email).Scan(&u.ID, &u.DisplayName, &u.Email, &u.Status, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &identity.User{}, identity.ErrNotFound
		}
		return &identity.User{}, err
	}
	return &u, nil
}

func (ir *identityRepository) UpdateUser(ctx context.Context, u *identity.User) error {
	const query = `UPDATE table_users SET display_name = $2, updated_at = $3 WHERE id = $1`

	result, err := ir.db.Exec(ctx, query, u.ID, u.DisplayName, u.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_lower_unique" {
				return identity.ErrEmailConflict
			}
		}
		return err
	}
	if result.RowsAffected() == 0 {
		return identity.ErrNotFound
	}
	return nil
}

func (ir *identityRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status identity.Status, updatedAt time.Time) error {
	const query = `UPDATE table_users SET status=$2, updated_at=$3 WHERE id= $1`
	result, err := ir.db.Exec(ctx, query, id, status, updatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return identity.ErrNotFound
	}
	return nil
}
