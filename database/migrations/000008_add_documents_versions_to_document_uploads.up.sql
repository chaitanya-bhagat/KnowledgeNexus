ALTER TABLE table_document_uploads ADD COLUMN document_version_id UUID;

ALTER TABLE table_document_uploads ADD CONSTRAINT document_upload_document_version_fk
FOREIGN KEY (document_version_id) REFERENCES table_document_versions(id);