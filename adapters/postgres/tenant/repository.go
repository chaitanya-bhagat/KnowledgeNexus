package adaptertenant

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tenantRepository struct {
	db *pgxpool.Pool
}

func NewTenantRepository(db *pgxpool.Pool) *tenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (tr *tenantRepository) Create(ctx context.Context, t *tenantmodel.Tenant) error {
	const query = `
		INSERT INTO table_tenants (
			id,
			name,
			slug,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err := tr.db.Exec(ctx, query, t.ID, t.Name, t.Slug, t.Status, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" && pgErr.ConstraintName == "tenants_slug_lower_unique" {
				return tenant.ErrSlugConflict
			}
		}
		return err
	}
	return nil
}

func (tr *tenantRepository) GetByID(ctx context.Context, id uuid.UUID) (*tenantmodel.Tenant, error) {
	const query = `
		SELECT
			id,
			name,
			slug,
			status,
			created_at,
			updated_at
		FROM table_tenants
		WHERE id = $1
	`

	var t tenantmodel.Tenant

	err := tr.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&t.ID,
		&t.Name,
		&t.Slug,
		&t.Status,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &tenantmodel.Tenant{}, tenant.ErrNotFound
		}

		return &tenantmodel.Tenant{}, err
	}

	return &t, nil
}

func (tr *tenantRepository) Update(ctx context.Context, t *tenantmodel.Tenant) error {
	const query = `
		UPDATE table_tenants
		SET
			name = $2,
			slug = $3,
			updated_at = $4
		WHERE id = $1
	`

	result, err := tr.db.Exec(
		ctx,
		query,
		t.ID,
		t.Name,
		t.Slug,
		t.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == "23505" &&
			pgErr.ConstraintName == "tenants_slug_lower_unique" {
			return tenant.ErrSlugConflict
		}

		return err
	}

	if result.RowsAffected() == 0 {
		return tenant.ErrNotFound
	}

	return nil
}

func (tr *tenantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status tenantmodel.Status, updatedAt time.Time) error {
	const query = `
		UPDATE table_tenants
		SET
			status = $2,
			updated_at = $3
		WHERE id = $1
	`

	result, err := tr.db.Exec(
		ctx,
		query,
		id,
		status,
		updatedAt,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return tenant.ErrNotFound
	}

	return nil
}
