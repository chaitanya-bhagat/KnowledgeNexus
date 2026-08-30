CREATE TABLE IF NOT EXISTS table_document_versions(
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL,

    version_number INTEGER NOT NULL,

    object_key TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    content_type TEXT,
    size_bytes BIGINT,
    checksum TEXT,

    status TEXT NOT NULL DEFAULT 'uploaded',
    failure_reason TEXT,

    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT document_versions_tenant_fk FOREIGN KEY (tenant_id) REFERENCES table_tenants(id),

    CONSTRAINT document_versions_document_tenant_fk FOREIGN KEY (tenant_id, document_id) REFERENCES table_documents(tenant_id, id),

    CONSTRAINT document_versions_created_by_fk FOREIGN KEY (created_by) REFERENCES table_users(id),

    CONSTRAINT document_versions_number_unique UNIQUE (document_id, version_number),

    CONSTRAINT document_versions_version_number_check CHECK (version_number > 0),

    CONSTRAINT document_versions_size_check CHECK (size_bytes IS NULL OR size_bytes >= 0),

    CONSTRAINT document_versions_status_check CHECK (status IN ('uploaded', 'processing', 'ready', 'failed'))
);