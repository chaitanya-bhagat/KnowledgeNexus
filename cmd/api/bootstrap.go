package main

import (
	"context"
	"fmt"
	"net/http"

	httpadapter "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/health"
	identityhandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/identity"
	knowledgebasehandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/base"
	documenthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/document"
	documentversionhandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/knowledgebase/documentversion"
	membershiphandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant/memebership"
	tenanthandler "github.com/chaitanya-bhagat/knowledge-nexus/adapters/http/tenant/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres"
	adapteridentity "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/identity"
	adapterknowledgebase "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/knowledgebase/base"
	adapterdocument "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/knowledgebase/document"
	adapterdocumentversion "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/knowledgebase/documentversion"
	adaptertenant "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/tenant"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/identity"
	knowledgebase "github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/base"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/document"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/knowledgebase/documentversion"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/membership"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/tenant/tenant"
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
	membershipService := membership.NewMembershipService(tenantRepo, membershipRepo)
	membershipHandler := membershiphandler.NewMembershipHandler(membershipService, logger)

	identityRepo := adapteridentity.NewIdentityRepository(dbPool)
	identityService := identity.NewIdentityService(identityRepo, logger)
	identityHandler := identityhandler.NewIdentityHandler(identityService, logger)

	kbRepo := adapterknowledgebase.NewKnowledgeBase(dbPool)
	kbService := knowledgebase.NewKnowledgeBaseService(kbRepo, tenantRepo, identityRepo)
	kbHandler := knowledgebasehandler.NewKnowledgeBaseHandler(*kbService, logger)

	documentRepo := adapterdocument.NewDocumentRepository(dbPool)
	documentService := document.NewDocumentService(documentRepo, tenantRepo, kbRepo, identityRepo, logger)
	documentHandler := documenthandler.NewDocumentHandler(documentService, logger)

	documentVersionRepo := adapterdocumentversion.NewDocumentVersion(dbPool)
	documentVersionService := documentversion.NewDocumentVersionService(documentVersionRepo, documentRepo, kbRepo, tenantRepo, identityRepo, logger)
	documentVersHandler := documentversionhandler.NewDocumentVersionHandler(documentVersionService, logger)

	healthHandler := health.NewHealthHandler(dbPool)

	server := &http.Server{
		Addr: cfg.Server.Address(),
		Handler: httpadapter.LoadRoutes(httpadapter.Handlers{
			Health:          healthHandler,
			Tenant:          tenantHandler,
			Membership:      membershipHandler,
			Identity:        identityHandler,
			KnowledgeBase:   kbHandler,
			Document:        documentHandler,
			DocumentVersion: documentVersHandler,
		}),
	}
	return NewApp(server, logger, dbPool), nil

}
