package storage

import (
	"context"
	"time"
)

type EventRepo interface {
	GetEvent(ctx context.Context, id string) (Event, error)
	GetEventsInPeriod(ctx context.Context, startDate time.Time, endDate time.Time) ([]Event, error)
	GetEvents(ctx context.Context) ([]Event, error)
	CreateEvent(ctx context.Context, e Event) error
	UpdateEvent(ctx context.Context, e Event) error
	DeleteEvent(ctx context.Context, id string) error
}
