package httpadapter

import "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"

type Handlers struct {
	Health *health.HealthHandler
}
