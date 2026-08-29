CREATE TABLE IF NOT EXISTS table_tenants(
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT tenants_status_check
        CHECK (status IN ('active', 'disabled'))
);

CREATE UNIQUE INDEX tenants_slug_lower_unique
    ON tenants (LOWER(slug));