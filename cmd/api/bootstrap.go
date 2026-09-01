package main

import (
	"context"
	"fmt"
	"net/http"

	httpadapter "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	tenanthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres"
	adaptertenant "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant"
	"go.uber.org/zap"
)

func buildApp(ctx context.Context, cfg config.Config, logger *zap.Logger) (*App, error) {

	dbPool, err := postgres.NewPostgresPool(ctx, cfg.Postgresql.DNS())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db %w", err)
	}

	tenantRepo := adaptertenant.NewTenantRepository(dbPool)
	tenantService := tenant.NewTenantService(tenantRepo)
	tenantHandler := tenanthandler.NewTenantHandler(tenantService, logger)

	membershipRepo := adaptertenant.NewMembershipRepository(dbPool)
	membershipService := tenant.NewMembershipService(tenantRepo, membershipRepo)
	membershipHandler := tenanthandler.NewMembershipHandler(membershipService, logger)

	healthHandler := health.NewHealthHandler(dbPool)

	server := &http.Server{
		Addr: cfg.Server.Address(),
		Handler: httpadapter.LoadRoutes(httpadapter.Handlers{
			Health:     healthHandler,
			Tenant:     tenantHandler,
			Membership: membershipHandler,
		}),
	}
	return NewApp(server, logger, dbPool), nil

}
