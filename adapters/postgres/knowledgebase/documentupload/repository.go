package documentupload

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentupload"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type documentUploadRepository struct {
	db *pgxpool.Pool
}

func NewDocumentUpload(db *pgxpool.Pool) *documentUploadRepository {
	return &documentUploadRepository{
		db: db,
	}
}

var _ documentupload.Repository = (*documentUploadRepository)(nil)

func (dvr *documentUploadRepository) Create(ctx context.Context, upload kbmodel.DocumentUpload) error {
	const query = `
	INSERT INTO table_document_uploads (id, tenant_id, document_id, object_key, original_filename, content_type, expected_size_bytes, status, created_by, expire_at, created_at, completed_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := dvr.db.Exec(ctx, query, upload.ID, upload.TenantID, upload.DocumentID, upload.ObjectKey, upload.OriginalFileName, upload.ContentType, upload.ExpectedSizeBytes, upload.Status, upload.CreatedBy, upload.ExpiredAt, upload.CreatedAt, upload.CompletedAt)
	if err != nil {
		return fmt.Errorf("failed to upload document: %w", err)
	}
	return nil
}

func (dvr *documentUploadRepository) GetByID(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) (kbmodel.DocumentUpload, error) {
	const query = `
	SELECT id, tenant_id, document_id, object_key, original_filename, COALESCE(content_type, ''), COALESCE(expected_size_bytes, 0), status, created_by, expire_at, created_at, completed_at
	FROM table_document_uploads
	WHERE tenant_id = $1 AND id = $2
	`
	var upload kbmodel.DocumentUpload
	err := dvr.db.QueryRow(ctx, query, tenantID, uploadID).Scan(&upload.ID, &upload.TenantID, &upload.DocumentID, &upload.ObjectKey,
		&upload.OriginalFileName, &upload.ContentType, &upload.ExpectedSizeBytes, &upload.Status, &upload.CreatedBy, &upload.ExpiredAt,
		&upload.CreatedAt, &upload.CompletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kbmodel.DocumentUpload{}, documentupload.ErrDocumentUploadNotFound
		}
		return kbmodel.DocumentUpload{}, fmt.Errorf("failed to upload document: %w", err)
	}
	return upload, nil
}

func (dvr *documentUploadRepository) MarkCompleted(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error {
	const query = `
	UPDATE table_document_uploads
	SET status = 'completed', completed_at= NOW()
	WHERE tenant_id=$1 AND id = $2
	`
	result, err := dvr.db.Exec(ctx, query, tenantID, uploadID)
	if err != nil {
		return fmt.Errorf("mark document upload completed failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return documentupload.ErrDocumentUploadNotFound
	}
	return nil
}

func (dvr *documentUploadRepository) MarkExpired(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error {
	const query = `
	UPDATE table_document_uploads
	SET status = 'expired'
	WHERE tenant_id=$1 AND id = $2
	`
	result, err := dvr.db.Exec(ctx, query, tenantID, uploadID)
	if err != nil {
		return fmt.Errorf("mark document upload completed failed: %w", err)
	}
	if result.RowsAffected() == 0 {
		return documentupload.ErrDocumentUploadNotFound
	}
	return nil
}
