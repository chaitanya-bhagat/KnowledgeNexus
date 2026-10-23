package nats

import (
	"fmt"
	"time"

	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatsClient struct {
	Conn      *natsgo.Conn
	JetStream jetstream.JetStream
}

func NewNatsClient(url string) (*NatsClient, error) {
	conn, err := natsgo.Connect(
		url,
		natsgo.Name("knowledge-nexus-outbox-publisher"),
		natsgo.Timeout(5*time.Second),
		natsgo.MaxReconnects(-1),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS server: %w", err)
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}
	return &NatsClient{
		Conn:      conn,
		JetStream: js,
	}, nil

}

func (c *NatsClient) Close() {
	c.Conn.Close()
}
