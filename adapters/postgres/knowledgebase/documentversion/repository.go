package adapterdocumentversion

import (
	"context"
	"errors"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentversion"
	kbmodel "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type documentVersionRepository struct {
	db *pgxpool.Pool
}

func NewDocumentVersion(db *pgxpool.Pool) *documentVersionRepository {
	return &documentVersionRepository{
		db: db,
	}
}

func (dvr *documentVersionRepository) Create(ctx context.Context, version kbmodel.DocumentVersion) error {

	tx, err := dvr.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	const lockQuery = `
	SELECT id 
	FROM table_documents 
	WHERE tenant_id = $1 AND id =$2
	FOR UPDATE
	`

	var documentID uuid.UUID

	err = tx.QueryRow(ctx, lockQuery, version.TenantID, version.DocumentID).Scan(&documentID)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			return documentversion.ErrDocumentNotFound
		}
		return err
	}

	const nextVersionQuery = `
	SELECT COALESCE(MAX(version_number), 0) + 1
	FROM table_document_versions
	WHERE document_id = $1
	`
	var nextVersion int
	err = tx.QueryRow(ctx, nextVersionQuery, version.ID).Scan(&nextVersion)
	if err != nil {
		return err
	}

	version.VersionNumber = nextVersion

	const insertQuery = `
	INSERT INTO table_document_versions (id, tenant_id, document_id, version_number, object_key, original_filename, content_type, size_bytes, checksum, status, failure_reason, created_by, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err = tx.Exec(ctx, insertQuery, version.ID, version.TenantID, version.DocumentID, version.VersionNumber, version.ObjectKey, version.OriginalFileName, version.ContentType, version.SizeBytes, version.Checksum, version.Status, version.FailureReason, version.CreatedBy, version.CreatedAt)
	if err != nil {
		return mapError(err)
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (dvr *documentVersionRepository) GetByID(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID) (kbmodel.DocumentVersion, error) {
	const query = `
	SELECT 	id, tenant_id, document_id, version_number, object_key, original_filename, COALESCE(content_type, ''), COALESCE(size_bytes, 0), COALESCE(checksum, ''), status, COALESCE(failure_reason,''), created_by, created_at
	FROM table_document_versions
	WHERE tenant_id = $1 AND id = $2
	`

	var version kbmodel.DocumentVersion

	err := dvr.db.QueryRow(ctx, query, tenantID, versionID).Scan(&version.ID, &version.TenantID, &version.DocumentID, &version.VersionNumber, &version.ObjectKey,
		&version.OriginalFileName, &version.ContentType, &version.SizeBytes, &version.Checksum, &version.Status, &version.FailureReason, &version.CreatedBy, &version.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return kbmodel.DocumentVersion{}, documentversion.ErrVersionNotFound
		}
		return kbmodel.DocumentVersion{}, err
	}
	return version, nil
}

func (dvr *documentVersionRepository) GetList(ctx context.Context, tenantID uuid.UUID, docID uuid.UUID) ([]kbmodel.DocumentVersion, error) {
	const query = `
	SELECT 	id, tenant_id, document_id, version_number, object_key, original_filename, COALESCE(content_type, ''), COALESCE(size_bytes, 0), COALESCE(checksum, ''), status, COALESCE(failure_reason,''), created_by, created_at
	FROM table_document_versions
	WHERE tenant_id = $1 AND document_id = $2
	ORDER BY version_number ASC
	`
	rows, err := dvr.db.Query(ctx, query, tenantID, docID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	versions := make([]kbmodel.DocumentVersion, 0)

	for rows.Next() {
		var version kbmodel.DocumentVersion
		err := rows.Scan(&version.ID, &version.TenantID, &version.DocumentID, &version.VersionNumber, &version.ObjectKey,
			&version.OriginalFileName, &version.ContentType, &version.SizeBytes, &version.Checksum, &version.Status, &version.FailureReason, &version.CreatedBy, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return versions, nil
}

func (dvr *documentVersionRepository) Update(ctx context.Context, tenantID uuid.UUID, versionID uuid.UUID, status kbmodel.DocumentVersionStatus, failureReason string) error {
	const query = `
	UPDATE table_document_versions 
	SET status = $1, failure_reason = $2
	WHERE tenant_id = $3 AND id = $4
	`
	result, err := dvr.db.Exec(ctx, query, status, nullableString(failureReason), tenantID, versionID)
	if err != nil {
		return mapError(err)
	}
	if result.RowsAffected() == 0 {
		return documentversion.ErrVersionNotFound
	}
	return nil
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}

	return value
}

func mapError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	switch pgErr.Code {
	case "23503":
		switch pgErr.ConstraintName {
		case "document_versions_document_tenant_fk":
			return document.ErrDocumentNotFound
		case "document_versions_tenant_fk":
			return document.ErrTenantNotFound
		case "document_versions_created_by_fk":
			return document.ErrUserNotFound
		}

	}
	return nil
}
