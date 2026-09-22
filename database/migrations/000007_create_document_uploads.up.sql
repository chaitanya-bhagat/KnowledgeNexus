CREATE TABLE IF NOT EXISTS table_document_uploads(
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL,
    document_id UUID NOT NULL,

    object_key TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    content_type TEXT,
    expected_size_bytes BIGINT,
    
    status TEXT NOT NULL DEFAULT 'uploaded',
    
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    expire_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ,

    CONSTRAINT document_uploads_tenant_fk FOREIGN KEY (tenant_id) REFERENCES table_tenants(id),

    CONSTRAINT document_uploads_document_tenant_fk FOREIGN KEY (tenant_id, document_id) REFERENCES table_documents(tenant_id, id),

    CONSTRAINT document_uploads_created_by_fk FOREIGN KEY (created_by) REFERENCES table_users(id),


    CONSTRAINT document_uploads_size_check CHECK (expected_size_bytes IS NULL OR expected_size_bytes >= 0),

    CONSTRAINT document_uploads_status_check CHECK (status IN ('pending', 'completed', 'expired'))
);