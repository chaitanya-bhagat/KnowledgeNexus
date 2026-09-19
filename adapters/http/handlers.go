package httpadapter

import (
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	identityhandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/identity"
	knowledgebasehandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/base"
	documenthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/document"
	documentversionhandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/documentversion"
	tenanthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant"
)

type Handlers struct {
	Health          *health.HealthHandler
	Tenant          *tenanthandler.TenantHandler
	Membership      *tenanthandler.MembershipHandler
	Identity        *identityhandler.IdentityHandler
	KnowledgeBase   *knowledgebasehandler.KnowledgeBaseHandler
	Document        *documenthandler.DocumentHandler
	DocumentVersion *documentversionhandler.DocumentVersionHandler
}
