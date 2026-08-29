CREATE TABLE IF NOT EXISTS table_knowledge_bases(
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    domain_type TEXT NOT NULL,

    status TEXT NOT NULL DEFAULT 'active',

    created_by UUID NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT knowledge_bases_tenant_fk FOREIGN KEY (tenant_id) REFERENCES table_tenants(id),
    CONSTRAINT knowledge_bases_created_by_fk FOREIGN KEY (created_by) REFERENCES table_users(id),

    CONSTRAINT knowledge_bases_name_unique UNIQUE (tenant_id, name),
    CONSTRAINT knowledge_bases_tenant_id_id_unique UNIQUE (tenant_id, id),

    CONSTRAINT knowledge_bases_status_check CHECK (status IN ('active', 'archived'))
);