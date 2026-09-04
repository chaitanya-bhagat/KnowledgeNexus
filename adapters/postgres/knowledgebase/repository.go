package adapterknowledgebase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type kbRepository struct {
	db *pgxpool.Pool
}

func NewKnowledgeBase(db *pgxpool.Pool) *kbRepository {
	return &kbRepository{
		db: db,
	}
}

func (kbr *kbRepository) CreateKnowledgeBase(ctx context.Context, kb knowledgebase.KnowledgeBase) error {
	const query = `
	INSERT INTO table_knowledge_bases (id, tenant_id, name, description, domain_type, status, created_by, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := kbr.db.Exec(ctx, query, kb.ID, kb.TenantID, kb.Name, kb.Description, kb.DomainType, kb.Status, kb.CreatedBy, kb.CreatedAt, kb.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (kbr *kbRepository) GetKnowledgeBaseByID(ctx context.Context, kbID uuid.UUID, tenantID uuid.UUID) (knowledgebase.KnowledgeBase, error) {
	const query = `
	SELECT id, tenant_id, name, description, domain_type, status, created_by, created_at, updated_at
	FROM table_knowledge_bases
	WHERE tenant_id = $1 AND id = $2
	`
	var kb knowledgebase.KnowledgeBase
	err := kbr.db.QueryRow(ctx, query, tenantID, kbID).Scan(&kb.ID, &kb.TenantID, &kb.Name, &kb.Description,
		&kb.DomainType, &kb.Status, &kb.CreatedBy,
		&kb.CreatedAt, &kb.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return knowledgebase.KnowledgeBase{}, knowledgebase.ErrKnowledgeBaseNotFound
		}
		return knowledgebase.KnowledgeBase{}, err
	}
	return kb, nil
}

func (kbr *kbRepository) ListKnowledgeBasesByTenantID(ctx context.Context, tenantID uuid.UUID) ([]knowledgebase.KnowledgeBase, error) {
	const query = `
	SELECT id, tenant_id, name, description, domain_type, status, created_by, created_at, updated_at
	FROM table_knowledge_bases
	WHERE tenant_id = $1 
	ORDER BY created_at ASC
	`
	rows, err := kbr.db.Query(ctx, query, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kbs := make([]knowledgebase.KnowledgeBase, 0)
	for rows.Next() {
		var kb knowledgebase.KnowledgeBase

		if err = rows.Scan(&kb.ID, &kb.TenantID, &kb.Name, &kb.Description, &kb.DomainType, &kb.Status,
			&kb.CreatedBy, &kb.CreatedAt, &kb.UpdatedAt); err != nil {
			return nil, err
		}
		kbs = append(kbs, kb)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return kbs, nil
}

func (kbr *kbRepository) UpdateKnowledgeBase(ctx context.Context, kb *knowledgebase.KnowledgeBase) error {
	const query = `
	UPDATE table_knowledge_bases SET name=$3, description=$4, updated_at=$5
	WHERE tenant_id=$1 AND id=$2
	`
	result, err := kbr.db.Exec(ctx, query, kb.TenantID, kb.ID, kb.Name, kb.Description, kb.UpdatedAt)
	if err != nil {
		return mapError(err)
	}
	if result.RowsAffected() == 0 {
		return knowledgebase.ErrKnowledgeBaseNotFound
	}

	return nil
}

func (kbr *kbRepository) UpdateKnowledgeBaseStatus(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID, status knowledgebase.Status, updatedAt time.Time) error {
	const query = `
	UPDATE table_knowledge_bases SET status=$3, updated_at=$4
	WHERE tenant_id=$1 AND id=$2
	`

	result, err := kbr.db.Exec(ctx, query, tenantID, kbID, status, updatedAt)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return knowledgebase.ErrKnowledgeBaseNotFound
	}
	return nil
}

func mapError(err error) error {
	var pgErr *pgconn.PgError

	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23505":
		if pgErr.ConstraintName == "knowledge_bases_name_unique" {
			return knowledgebase.ErrKnowledgeBaseConflict
		}

	case "23503":
		switch pgErr.ConstraintName {
		case "knowledge_bases_tenant_fk":
			return knowledgebase.ErrTenantNotFound

		case "knowledge_bases_created_by_fk":
			return knowledgebase.ErrUserNotFound
		}
	}

	return err
}
