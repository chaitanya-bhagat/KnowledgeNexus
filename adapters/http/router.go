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

	router.Post("/users/create", deps.Identity.CreateUser)
	router.Get("/users/{userID}", deps.Identity.GetUserByID)
	router.Get("/users/getbyemail", deps.Identity.GetUserByEmail)
	router.Post("/users/disable", deps.Identity.Disable)
	router.Post("/users/enable", deps.Identity.Enable)

	router.Post("/knowledge-bases", deps.KnowledgeBase.Create)
	router.Post("/knowledge-bases/get", deps.KnowledgeBase.GetKnowledgeBase)
	router.Get("/knowledge-bases/{tenantID}/list", deps.KnowledgeBase.GetList)
	router.Post("/knowledge-bases/update", deps.KnowledgeBase.Update)
	router.Post("/knowledge-bases/archive", deps.KnowledgeBase.Archive)
	router.Post("/knowledge-bases/activate", deps.KnowledgeBase.Activate)

	router.Post("/documents", deps.Document.Create)
	router.Post("/documents/get", deps.Document.GetByID)
	router.Post("/documents/list", deps.Document.GetList)
	router.Post("/documents/update", deps.Document.Update)
	router.Post("/documents/archive", deps.Document.Archive)
	router.Post("/documents/active", deps.Document.Activate)

	return router

}
