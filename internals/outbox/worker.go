package outbox

import (
	"context"
	"fmt"
	"time"
)

type Worker struct {
	outboxService *OutboxService
	interval      time.Duration
	batchSize     int
}

func NewWorker(outboxService *OutboxService, interval time.Duration, batchSize int) (*Worker, error) {
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
			return fmt.Errorf("failed to process unpublished events: %w", err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
