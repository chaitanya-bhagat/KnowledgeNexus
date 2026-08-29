CREATE TABLE IF NOT EXISTS table_documents(
    id UUID PRIMARY KEY,

    tenant_id UUID NOT NULL,
    knowledge_base_id UUID NOT NULL,

    title TEXT NOT NULL,
    document_type TEXT,

    status TEXT NOT NULL DEFAULT 'active',

    created_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT documents_tenant_fk FOREIGN KEY (tenant_id) REFERENCES table_tenants(id),

    CONSTRAINT documents_knowledge_base_tenant_fk FOREIGN KEY (tenant_id, knowledge_base_id) REFERENCES table_knowledge_bases (tenant_id, id),

    CONSTRAINT documents_created_by_fk FOREIGN KEY (created_by) REFERENCES table_users(id),

    CONSTRAINT documents_tenant_id_id_unique UNIQUE (tenant_id, id),

    CONSTRAINT documents_status_check CHECK (status IN ('active', 'archived'))
);