package main

import (
	"net/http"

	httpadapter "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"go.uber.org/zap"
)

func buildApp(cfg config.Config, logger *zap.Logger) (*App, error) {

	server := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: httpadapter.LoadRoutes(),
	}
	return NewApp(server, logger), nil

}
