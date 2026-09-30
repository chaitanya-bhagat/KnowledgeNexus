package adapterdocumentupload

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentupload"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentversion"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type documentUploadRepository struct {
	db *pgxpool.Pool
}

func NewDocumentUploadRepository(db *pgxpool.Pool) *documentUploadRepository {
	return &documentUploadRepository{
		db: db,
	}
}

var _ documentupload.Repository = (*documentUploadRepository)(nil)
var _ documentupload.CompletionRepository = (*documentUploadRepository)(nil)

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
	SELECT id, tenant_id, document_id, document_version_id, object_key, original_filename, COALESCE(content_type, ''), COALESCE(expected_size_bytes, 0), status, created_by, expire_at, created_at, completed_at
	FROM table_document_uploads
	WHERE tenant_id = $1 AND id = $2
	`
	var upload kbmodel.DocumentUpload
	err := dvr.db.QueryRow(ctx, query, tenantID, uploadID).Scan(&upload.ID, &upload.TenantID, &upload.DocumentID, &upload.DocumentVersionID, &upload.ObjectKey,
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

func (dvr *documentUploadRepository) MarkCompleted(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) (kbmodel.DocumentVersion, error) {
	tx, err := dvr.db.Begin(ctx)
	if err != nil {
		return kbmodel.DocumentVersion{}, fmt.Errorf("being upload completion transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	upload, err := getUploadForUpdate(ctx, tx, tenantID, uploadID)
	if err != nil {
		return kbmodel.DocumentVersion{}, err
	}
	if upload.Status == kbmodel.DocumentUploadCompleted {
		if upload.DocumentVersionID == nil {
			return kbmodel.DocumentVersion{}, fmt.Errorf("completed upload %s has no document version", uploadID)
		}
		version, err := getDocumentVersion(ctx, tx, tenantID, *upload.DocumentVersionID)
		if err != nil {
			return kbmodel.DocumentVersion{}, err
		}
		if err := tx.Commit(ctx); err != nil {
			return kbmodel.DocumentVersion{}, fmt.Errorf("failed to commit upload completion transaction: %w", err)
		}
		return version, nil
	}
	if upload.Status != kbmodel.DocumentUploadPending {
		return kbmodel.DocumentVersion{}, documentupload.ErrInvalidStatus
	}

	if err := lockParentDocument(ctx, tx, tenantID, upload.DocumentID); err != nil {
		return kbmodel.DocumentVersion{}, err
	}

	versionNumber, err := nextVersionNumber(ctx, tx, tenantID, upload.DocumentID)
	if err != nil {
		return kbmodel.DocumentVersion{}, err
	}

	createVersion, err := createDocumentVersion(ctx, tx, upload, versionNumber)
	if err != nil {
		return kbmodel.DocumentVersion{}, err
	}

	if err := completeUpload(ctx, tx, tenantID, uploadID, createVersion.ID); err != nil {
		return kbmodel.DocumentVersion{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return kbmodel.DocumentVersion{}, fmt.Errorf("failed to commit upload completion transaction: %w", err)
	}
	return createVersion, nil
}

// func (dvr *documentUploadRepository) MarkExpired(ctx context.Context, tenantID uuid.UUID, uploadID uuid.UUID) error {
// 	const query = `
// 	UPDATE table_document_uploads
// 	SET status = 'expired'
// 	WHERE tenant_id=$1 AND id = $2
// 	`
// 	result, err := dvr.db.Exec(ctx, query, tenantID, uploadID)
// 	if err != nil {
// 		return fmt.Errorf("mark document upload completed failed: %w", err)
// 	}
// 	if result.RowsAffected() == 0 {
// 		return documentupload.ErrDocumentUploadNotFound
// 	}
// 	return nil
// }

func completeUpload(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID, uploadID uuid.UUID, versionID uuid.UUID) error {
	query := `
		UPDATE table_document_uploads
		SET status = 'completed', document_version_id = $3, completed_at = NOW()
		WHERE tenant_id = $1 AND id = $2 AND status = 'pending'
		`
	tag, err := tx.Exec(ctx, query, tenantID, uploadID, versionID)
	if err != nil {
		return fmt.Errorf("failed to complete document upload: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return documentupload.ErrInvalidStatus
	}

	return nil
}

func createDocumentVersion(ctx context.Context, tx pgx.Tx, upload kbmodel.DocumentUpload, versionNumber int) (kbmodel.DocumentVersion, error) {
	versionID := uuid.New()

	var version kbmodel.DocumentVersion

	err := tx.QueryRow(
		ctx,
		`
		INSERT INTO table_document_versions (id, tenant_id, document_id, version_number, object_key, original_filename, status, created_by)
		VALUES ($1, $2, $3, $4, $5, $6, 'uploaded', $7)
		RETURNING id, tenant_id, document_id, version_number, object_key, status, created_by, created_at
		`,
		versionID, upload.TenantID, upload.DocumentID, versionNumber, upload.ObjectKey, upload.OriginalFileName, upload.CreatedBy).Scan(
		&version.ID, &version.TenantID, &version.DocumentID, &version.VersionNumber, &version.ObjectKey, &version.Status,
		&version.CreatedBy, &version.CreatedAt)

	if err != nil {
		return kbmodel.DocumentVersion{}, fmt.Errorf("create document version: %w", err)
	}

	return version, nil
}

func nextVersionNumber(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID, documentID uuid.UUID) (int, error) {
	var versionNumber int
	query := `
		SELECT COALESCE(MAX(version_number), 0) + 1
		FROM table_document_versions
		WHERE tenant_id = $1 AND document_id = $2
		`
	err := tx.QueryRow(ctx, query, tenantID, documentID).Scan(&versionNumber)
	if err != nil {
		return 0, fmt.Errorf("get next document version number: %w", err)
	}

	return versionNumber, nil
}
func lockParentDocument(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID, docID uuid.UUID) error {
	query := `
	SELECT id FROM table_documents WHERE tenant_id=$1 AND id=$2 FOR UPDATE	
	`
	var id uuid.UUID
	err := tx.QueryRow(ctx, query, tenantID, docID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return documentupload.ErrDocumentUploadNotFound
		}
		return fmt.Errorf("failed to lock the document %w", err)
	}
	return nil
}

func getUploadForUpdate(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID, uploadID uuid.UUID) (kbmodel.DocumentUpload, error) {
	var upload kbmodel.DocumentUpload
	query := `
	SELECT id, tenant_id, document_id, document_version_id, object_key, original_filename, content_type, expected_size_bytes, status, created_by, expire_at, created_at, completed_at
	FROM table_document_uploads 
	WHERE tenant_id=$1 AND id=$2
	FOR UPDATE`
	err := tx.QueryRow(ctx, query, tenantID, uploadID).Scan(&upload.ID, &upload.TenantID, &upload.DocumentID, &upload.DocumentVersionID, &upload.ObjectKey,
		&upload.OriginalFileName, &upload.ContentType, &upload.ExpectedSizeBytes, &upload.Status, &upload.CreatedBy, &upload.ExpiredAt,
		&upload.CreatedAt, &upload.CompletedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kbmodel.DocumentUpload{}, documentupload.ErrDocumentUploadNotFound
		}
		return kbmodel.DocumentUpload{}, err
	}

	return upload, nil
}

func getDocumentVersion(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID, versionID uuid.UUID) (kbmodel.DocumentVersion, error) {
	const query = `
	SELECT 	id, tenant_id, document_id, version_number, object_key,  status, created_by, created_at
	FROM table_document_versions
	WHERE tenant_id = $1 AND id = $2
	`

	var version kbmodel.DocumentVersion

	err := tx.QueryRow(ctx, query, tenantID, versionID).Scan(&version.ID, &version.TenantID, &version.DocumentID, &version.VersionNumber, &version.ObjectKey,
		&version.Status, &version.CreatedBy, &version.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kbmodel.DocumentVersion{}, documentversion.ErrVersionNotFound
		}
		return kbmodel.DocumentVersion{}, err
	}
	return version, nil

}
