package httpadapter

import (
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	identityhandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/identity"
	tenanthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant"
)

type Handlers struct {
	Health     *health.HealthHandler
	Tenant     *tenanthandler.TenantHandler
	Membership *tenanthandler.MembershipHandler
	Identity   *identityhandler.IdentityHandler
}
