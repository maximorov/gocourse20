package transport

import (
	"context"
	"solid/notification"
	"solid/sender"
)

func NewHttpCreator() sender.Transport {
	return &httpCreator{}
}

type httpCreator struct {
}

func (s *httpCreator) Name() string {
	return `http`
}

func (s httpCreator) Create(_ notification.Notification) (sender.Target, error) {
	return &httpTransport{}, nil
}

type httpTransport struct {
}

func (s *httpTransport) Send(_ context.Context, _ notification.Notification, _ notification.ItemToProcess) error {
	panic(`Not implemented`)
}

func (s *httpTransport) Die() error {
	return nil
}
