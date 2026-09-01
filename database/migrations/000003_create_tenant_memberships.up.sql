CREATE TABLE IF NOT EXISTS table_tenant_memberships(
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    user_id UUID NOT NULL,

    role TEXT NOT NULL DEFAULT 'member',
    status TEXT NOT NULL DEFAULT 'active',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tenant_memberships_tenant_fk FOREIGN KEY(tenant_id) REFERENCES table_tenants(id),
    CONSTRAINT tenant_memberships_user_fk FOREIGN KEY(user_id) REFERENCES table_users(id),

    CONSTRAINT tenant_memberships_unique UNIQUE(tenant_id, user_id),

    CONSTRAINT tenant_memberships_role_check CHECK(role IN('owner', 'admin', 'member')),
    CONSTRAINT tenant_memberships_status_check CHECK(status IN('active', 'disabled'))
);