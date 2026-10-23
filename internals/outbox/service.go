package outbox

import (
	"context"
	"fmt"
)

type OutboxService struct {
	repo      Repository
	publisher Publisher
}

func NewOutboxService(repo Repository, publisher Publisher) *OutboxService {
	return &OutboxService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *OutboxService) ProcessUnpublishedEvents(ctx context.Context, limit int) error {
	if limit <= 0 {
		return fmt.Errorf("limit must be greater than 0")
	}
	events, err := s.repo.GetUnpublishedEvents(ctx, limit)
	if err != nil {
		return fmt.Errorf("failed to get unpublished events: %w", err)
	}

	for _, event := range events {
		err := s.publisher.Publish(ctx, event)
		if err != nil {
			return fmt.Errorf("failed to publish event %s: %w", event.ID, err)
		}
		err = s.repo.MarkEventAsPublished(ctx, event.ID)
		if err != nil {
			return fmt.Errorf("failed to mark event %s as published: %w", event.ID, err)
		}
	}
	return nil
}
