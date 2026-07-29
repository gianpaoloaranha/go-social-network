package pubsub

import "context"

type Publisher interface {
	Publish(ctx context.Context, topic string, payload []byte) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string) (<-chan []byte, func(), error)
}