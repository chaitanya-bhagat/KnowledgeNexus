package httpadapter

import (
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	"github.com/go-chi/chi/v5"
)

func LoadRoutes() *chi.Mux {

	router := chi.NewRouter()

	healthHandler := health.NewHealthHandler()

	router.Get("/health/live", healthHandler.Live)
	router.Get("/health/ready", healthHandler.Ready)

	return router

}
