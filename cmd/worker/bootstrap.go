package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/messaging/nats"
	"github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres"
	adapteroutbox "github.com/chaitanya-bhagat/knowledge-nexus/adapters/postgres/outbox"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/logger"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/outbox"
)

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	logger, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		log.Fatal(err)
	}
	dbPool, err := postgres.NewPostgresPool(ctx, cfg.Postgresql.DNS())
	if err != nil {
		return fmt.Errorf("failed to initialize db %w", err)
	}
	defer dbPool.Close()

	outboxRepo := adapteroutbox.NewOutboxRepository(dbPool)
	natsClient, err := nats.NewNatsClient(cfg.Nats.URL)
	if err != nil {
		return fmt.Errorf("failed to create nats client %w", err)
	}
	defer natsClient.Close()

	publisher := nats.NewJetStreamPublisher(natsClient.JetStream)

	outboxService := outbox.NewOutboxService(outboxRepo, publisher)
	outboxWorker, err := outbox.NewWorker(outboxService, 2*time.Second, 100, logger)
	if err != nil {
		return fmt.Errorf("failed to create outbox worker %w", err)
	}
	logger.Info("worker started")
	return outboxWorker.Run(ctx)

}
