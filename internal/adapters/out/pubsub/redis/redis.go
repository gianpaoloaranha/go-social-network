package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Broker struct {
	client *redis.Client
}

func NewBroker(client *redis.Client) *Broker {
	return &Broker{client: client}
}

func (b *Broker) Publish(ctx context.Context, topic string, payload []byte) error {
	return b.client.Publish(ctx, topic, payload).Err()
}

func (b *Broker) Subscribe(ctx context.Context, topic string) (<-chan []byte, func(), error) {
	sub := b.client.Subscribe(ctx, topic)
	if _, err := sub.Receive(ctx); err != nil {
		_ = sub.Close()
		return nil, nil, err
	}

	redisCh := sub.Channel()

	out := make(chan []byte)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-redisCh:
				if !ok {
					return
				}
				out <- []byte(msg.Payload)
			}
		}
	}()

	cleanup := func() { _ = sub.Close() }
	return out, cleanup, nil
}
