package nats

import (
	"context"
	"fmt"

	"github.com/chaitanya-bhagat/knowledge-nexus/internals/outbox"
	"github.com/nats-io/nats.go/jetstream"
)

type Publisher struct {
	js jetstream.JetStream
}

var _ outbox.Publisher = (*Publisher)(nil)

func NewJetStreamPublisher(js jetstream.JetStream) *Publisher {
	return &Publisher{
		js: js,
	}
}

func (p *Publisher) Publish(ctx context.Context, event outbox.Event) error {
	_, err := p.js.Publish(ctx, event.EventType, event.Payload)
	if err != nil {
		return fmt.Errorf("failed to publish event to NATS JetStream: %w", err)
	}
	return nil
}
