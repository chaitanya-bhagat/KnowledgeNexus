package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type App struct {
	server *http.Server
	logger *zap.Logger
}

func NewApp(server *http.Server, logger *zap.Logger) *App {
	return &App{
		server: server,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		a.logger.Info("api server starting", zap.String("address", a.server.Addr))
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
			return
		}
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		a.logger.Info("shutdown signal received")
		shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutDownCtx); err != nil {
			return err
		}
		a.logger.Info("api server stopped gracefully")
		return <-errChan
	}

}
