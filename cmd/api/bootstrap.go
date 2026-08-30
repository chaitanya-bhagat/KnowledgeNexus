package main

import (
	"context"
	"fmt"
	"net/http"

	httpadapter "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"go.uber.org/zap"
)

func buildApp(ctx context.Context, cfg config.Config, logger *zap.Logger) (*App, error) {

	dbPool, err := postgres.NewPostgresPool(ctx, cfg.Postgresql.DNS())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db %w", err)
	}

	healthHandler := health.NewHealthHandler(dbPool)

	server := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: httpadapter.LoadRoutes(httpadapter.Handlers{Health: healthHandler}),
	}
	return NewApp(server, logger, dbPool), nil

}
