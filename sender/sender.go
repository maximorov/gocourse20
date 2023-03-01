package sender

import (
	"context"
	"fmt"
	"solid/notification"
	"sync"
	"time"
)

const targetTTL = 5 * time.Minute

type Target interface {
	Send(context.Context, notification.Notification, notification.ItemToProcess) error
	Die() error
}

// Dependency Inversion Principle
type Transport interface {
	Name() string
	Create(notification notification.Notification) (Target, error)
}

// Open-Closed Principle
func NewSender(ts ...Transport) notification.Sender {
	res := &sender{
		transports: ts,
		mu:         sync.RWMutex{},
		_targets:   make(map[int]Target),
	}

	go res.initTTL()

	return res
}

type sender struct {
	transports []Transport
	mu         sync.RWMutex
	_targets   map[int]Target
}

func (s *sender) Send(ctx context.Context, n notification.Notification, i notification.ItemToProcess) error {
	for _, transport := range s.transports {
		if transport.Name() == n.Transport() {
			if t, err := s.initTarget(transport, n); err != nil {
				return err
			} else {
				return t.Send(ctx, n, i)
			}
		}
	}

	return fmt.Errorf(`transport %s not defined`, n.Transport())
}

func (s *sender) initTarget(tp Transport, n notification.Notification) (Target, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if tt, ok := s._targets[n.ID()]; ok {
		return tt, nil
	}

	tt, err := tp.Create(n)
	if err != nil {
		return nil, err
	}

	s._targets[n.ID()] = tt

	return tt, nil
}

func (s *sender) initTTL() {
	c := time.NewTicker(targetTTL)
	for {
		<-c.C
		for i := range s._targets {
			_ = s._targets[i].Die() // add real TTL
		}
	}
}
