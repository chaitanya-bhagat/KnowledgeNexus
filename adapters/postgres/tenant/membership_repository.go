package adaptertenant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	tenantmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipRepository struct {
	db *pgxpool.Pool
}

func NewMembershipRepository(db *pgxpool.Pool) *MembershipRepository {
	return &MembershipRepository{
		db: db,
	}
}

// var _ tenantdomain.MembershipRepository = (*MembershipRepository)(nil)

func (mr *MembershipRepository) CreateMembership(ctx context.Context, m *tenantmodel.Membership) error {
	fmt.Println("mem", m.Status)
	const query = `
		INSERT INTO table_tenant_memberships (
			id,
			tenant_id,
			user_id,
			role,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := mr.db.Exec(
		ctx,
		query,
		m.ID,
		m.TenantID,
		m.UserID,
		m.Role,
		m.Status,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError

	fmt.Println("err 61", err)
	if !errors.As(err, &pgErr) {
		return err
	}

	switch {
	case pgErr.Code == "23505" &&
		pgErr.ConstraintName == "tenant_memberships_unique":

		return tenant.ErrMembershipExists

	case pgErr.Code == "23503" &&
		pgErr.ConstraintName == "tenant_memberships_tenant_fk":

		return tenant.ErrNotFound

	case pgErr.Code == "23503" &&
		pgErr.ConstraintName == "tenant_memberships_user_fk":

		return tenant.ErrUserNotFound

	default:
		return err
	}
}

func (mr *MembershipRepository) GetMembership(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID) (tenantmodel.Membership, error) {
	const query = `
		SELECT
			id,
			tenant_id,
			user_id,
			role,
			status,
			created_at,
			updated_at
		FROM table_tenant_memberships
		WHERE tenant_id = $1
		  AND user_id = $2
	`

	var m tenantmodel.Membership

	err := mr.db.QueryRow(
		ctx,
		query,
		tenantID,
		userID,
	).Scan(
		&m.ID,
		&m.TenantID,
		&m.UserID,
		&m.Role,
		&m.Status,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return tenantmodel.Membership{}, tenant.ErrMembershipNotFound
	}

	if err != nil {
		return tenantmodel.Membership{}, err
	}

	return m, nil
}

func (r *MembershipRepository) ListMemberships(ctx context.Context, tenantID uuid.UUID) ([]tenantmodel.Membership, error) {
	const query = `
		SELECT
			id,
			tenant_id,
			user_id,
			role,
			status,
			created_at,
			updated_at
		FROM table_tenant_memberships
		WHERE tenant_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		tenantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	memberships := make([]tenantmodel.Membership, 0)

	for rows.Next() {
		var m tenantmodel.Membership

		if err := rows.Scan(
			&m.ID,
			&m.TenantID,
			&m.UserID,
			&m.Role,
			&m.Status,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}

		memberships = append(memberships, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return memberships, nil
}

func (r *MembershipRepository) UpdateMembershipRole(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, role tenantmodel.Role, updatedAt time.Time) error {
	const query = `
		UPDATE table_tenant_memberships
		SET
			role = $3,
			updated_at = $4
		WHERE tenant_id = $1
		  AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		tenantID,
		userID,
		role,
		updatedAt,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return tenant.ErrMembershipNotFound
	}

	return nil
}

func (r *MembershipRepository) UpdateMembershipStatus(ctx context.Context, tenantID uuid.UUID, userID uuid.UUID, status tenantmodel.MembershipStatus, updatedAt time.Time) error {
	const query = `
		UPDATE table_tenant_memberships
		SET
			status = $3,
			updated_at = $4
		WHERE tenant_id = $1
		  AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		tenantID,
		userID,
		status,
		updatedAt,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return tenant.ErrMembershipNotFound
	}

	return nil
}
