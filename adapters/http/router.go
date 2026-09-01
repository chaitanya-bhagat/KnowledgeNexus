package httpadapter

import (
	"github.com/go-chi/chi/v5"
)

func LoadRoutes(deps Handlers) *chi.Mux {

	router := chi.NewRouter()

	// healthHandler := health.NewHealthHandler()

	router.Get("/health/live", deps.Health.Live)
	router.Get("/health/ready", deps.Health.Ready)

	router.Post("/tenant", deps.Tenant.Create)
	router.Get("/tenant/{id}", deps.Tenant.GetByID)
	router.Patch("/tenants/{id}", deps.Tenant.Update)
	router.Post("/tenants/{id}/disable", deps.Tenant.Disable)

	router.Post("/tenants/members", deps.Membership.Create)
	router.Get("/tenants/{tenantID}/members", deps.Membership.List)
	router.Post("/tenant-memberships/get", deps.Membership.Get)
	router.Post("/tenant-memberships/update-role", deps.Membership.UpdateRole)
	router.Post("/tenant-memberships/disable", deps.Membership.Disable)
	router.Post("/tenant-memberships/enable", deps.Membership.Enable)

	return router

}
