package adapterdocument

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type documentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) *documentRepository {
	return &documentRepository{
		db: db,
	}
}

func (dr *documentRepository) Create(ctx context.Context, document kbmodel.Document) error {
	const query = `
	INSERT INTO table_document(id, tenant_id, knowledge_base_id, title, document_type, status, created_by, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := dr.db.Exec(ctx, query, document.ID, document.TenantID, document.KbID, document.Title, document.DocumentType, document.Status, document.CreatedBy, document.CreatedAt, document.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (dr *documentRepository) GetByID(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) (kbmodel.Document, error) {

	const query = `
	SELECT id, tenant_id, knowledge_base_id, title,  COALESCE(document_type, ''), status, created_by, created_at, updated_at
	FROM table_document
	WHERE tenant_id = $1 AND id = $2
	`

	var doc kbmodel.Document
	err := dr.db.QueryRow(ctx, query, tenantID, docID).Scan(&doc.ID, &doc.TenantID, &doc.KbID, &doc.Title, &doc.DocumentType, &doc.Status, &doc.CreatedBy, &doc.CreatedAt, &doc.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kbmodel.Document{}, document.ErrDocumentNotFound
		}
		return kbmodel.Document{}, err
	}
	return doc, nil
}

func (dr *documentRepository) GetList(ctx context.Context, tenantID uuid.UUID, kbID uuid.UUID) ([]kbmodel.Document, error) {

	const query = `
	SELECT id, tenant_id, knowledge_base_id, title,  COALESCE(document_type, ''), status, created_by, created_at, updated_at
	FROM table_document
	WHERE tenant_id = $1 AND knowledge_base_id = $2
	ORDER BY created_at ASC
	`

	var docs []kbmodel.Document
	rows, err := dr.db.Query(ctx, query, tenantID, kbID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var doc kbmodel.Document
		err := rows.Scan(&doc.ID, &doc.TenantID, &doc.KbID, &doc.Title, &doc.DocumentType, &doc.Status, &doc.CreatedBy, &doc.CreatedAt, &doc.UpdatedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return docs, nil
}

func (dr *documentRepository) Update(ctx context.Context, doc kbmodel.Document) error {
	const query = `
	 UPDATE table_document SET title = $1, document_type = $2, updated_at = $3 
	 WHERE tenant_id = $4 AND id = $5 `

	result, err := dr.db.Exec(ctx, query, doc.Title, doc.DocumentType, doc.UpdatedAt, doc.TenantID, doc.ID)
	if err != nil {
		return mapError(err)
	}
	if result.RowsAffected() == 0 {
		return document.ErrDocumentNotFound
	}
	return nil
}

func (dr *documentRepository) UpdateStatus(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID, status kbmodel.DocumentStatus, updatedAt time.Time) error {
	const query = ` UPDATE table_document SET status = $1, updated_at = $2 
	WHERE tenant_id = $3 AND id = $4 `
	result, err := dr.db.Exec(ctx, query, status, updatedAt, tenantID, docID)
	if err != nil {
		return mapError(err)
	}
	if result.RowsAffected() == 0 {
		return document.ErrDocumentNotFound
	}
	return nil
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23503":
		switch pgErr.ConstraintName {
		case "documents_knowledge_base_tenant_fk":
			return document.ErrKnowledgeBaseNotFound
		case "documents_tenant_fk":
			return document.ErrTenantNotFound
		case "documents_created_by_fk":
			return document.ErrUserNotFound
		}

	}
	return nil
}
