package httpadapter

import (
	"github.com/go-chi/chi/v5"
)

func LoadRoutes(deps Handlers) *chi.Mux {

	router := chi.NewRouter()

	// healthHandler := health.NewHealthHandler()

	router.Get("/health/live", deps.Health.Live)
	router.Get("/health/ready", deps.Health.Ready)

	return router

}
