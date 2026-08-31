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

	return router

}
