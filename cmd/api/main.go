package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/config"
	"github.com/chaitanya-bhagat/knowledge-nexus/internals/logger"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logger.NewLogger(config.Logger)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = logger.Sync()
	}()

	app, err := buildApp(config, logger)
	if err != nil {
		log.Fatal(err)

	}
	// defer app.server.Close() Write Close function in app.go

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	defer stop()

	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}

}
