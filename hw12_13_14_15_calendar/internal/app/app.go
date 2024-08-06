package app

import (
	"context"
	"errors"
	"time"

	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage"
)

var ErrDateBusy = errors.New("date is already busy by another event")

type App struct {
	logger  *logger.Logger
	storage Storage
}

type Storage interface {
	storage.EventRepo
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) GetEventsInPeriod(ctx context.Context, startDate time.Time, endDate time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsInPeriod(ctx, startDate, endDate)
}

func (a *App) GetEvents(ctx context.Context) ([]storage.Event, error) {
	return a.storage.GetEvents(ctx)
}

func (a *App) CreateEvent(ctx context.Context, event storage.Event) error {
	if err := isDateBusy(ctx, a.storage, event); err != nil {
		return err
	}

	return a.storage.CreateEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	if err := isDateBusy(ctx, a.storage, event); err != nil {
		return err
	}

	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, id string) error {
	return a.storage.DeleteEvent(ctx, id)
}

func isDateBusy(ctx context.Context, s Storage, event storage.Event) error {
	events, err := s.GetEventsInPeriod(ctx, event.StartAt, event.EndAt)
	if err != nil {
		return err
	}
	if len(events) > 0 {
		return ErrDateBusy
	}

	return nil
}
