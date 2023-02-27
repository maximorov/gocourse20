package transport

import (
	"context"
	"solid/l/sender"
)

func NewRabbitMQCreator() sender.Transport {
	return &rabbitMQCreator{}
}

type rabbitMQCreator struct {
}

func (s *rabbitMQCreator) Name() string {
	return `rabbitmq`
}

func (s rabbitMQCreator) Create(_ ol.Notification) (sender.Target, error) {
	return &rabbitMQTransport{}, nil
}

type rabbitMQTransport struct {
}

func (s *rabbitMQTransport) Send(_ context.Context, _ ol.Notification, _ ol.ItemToProcess) error {
	panic(`Not implemented`)
}

func (s *rabbitMQTransport) Die() error {
	return nil
}
