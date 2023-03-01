package notification

import (
	"context"
	"time"
)

type PayloadType string

const (
	PayloadTypeDiff     PayloadType = `diff`
	PayloadTypeSnapshot PayloadType = `snapshot`
)

// Interface Segregation Principle
type (
	Queryer func(ctx context.Context, query string) ([]byte, error)

	Notification interface {
		AllTargets() []Target
		GetTarget(string) Target
		Description() string
		GraphqlQuery(context.Context) ([]byte, error)
		Transport() string
		ID() int
		IsItPossible() (bool, error)
		TargetSettings() string
	}

	Target interface {
		Entity() string
		Intersects([]string) bool
	}

	Holder interface {
		FindByEntity(ItemToProcess) []Notification
		FetchById(context.Context, int) (Notification, error)
		WindUp(context.Context) error
		WarmUp(context.Context) error
	}

	ItemToProcess interface {
		EntityName() string
		EventTime() time.Time
		EventType() PayloadType
		BeforeAfterFieldsDiff() ([]string, error)
	}

	Sender interface {
		Send(context.Context, Notification, ItemToProcess) error
	}
)
