package outbox

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type Worker struct {
	outboxService *OutboxService
	interval      time.Duration
	batchSize     int
	logger        *zap.Logger
}

func NewWorker(outboxService *OutboxService, interval time.Duration, batchSize int, logger *zap.Logger) (*Worker, error) {
	if outboxService == nil {
		return nil, fmt.Errorf("outboxService cannot be nil")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be greater than 0")
	}
	if batchSize <= 0 {
		return nil, fmt.Errorf("batchSize must be greater than 0")
	}

	return &Worker{
		outboxService: outboxService,
		interval:      interval,
		batchSize:     batchSize,
		logger:        logger,
	}, nil
}

func (w *Worker) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := w.outboxService.ProcessUnpublishedEvents(ctx, w.batchSize)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			w.logger.Error("failed to process unpublished events", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
