package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage"
)

var _ storage.EventRepo = (*Storage)(nil)

type Storage struct {
	mu     sync.RWMutex
	events []storage.Event
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) GetEvent(_ context.Context, id string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, e := range s.events {
		if e.ID == id {
			return e, nil
		}
	}

	return storage.Event{}, nil
}

func (s *Storage) GetEventsInPeriod(
	_ context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0)
	for _, e := range s.events {
		if e.StartAt == startDate ||
			e.EndAt == endDate ||
			e.StartAt.After(startDate) && e.StartAt.Before(endDate) ||
			e.EndAt.After(startDate) && e.EndAt.Before(endDate) {
			result = append(result, e)
		}
	}

	return result, nil
}

func (s *Storage) GetEvents(_ context.Context) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]storage.Event, len(s.events))
	copy(list, s.events)

	return list, nil
}

func (s *Storage) CreateEvent(_ context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, e)

	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, e storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < len(s.events); i++ {
		if s.events[i].ID == e.ID {
			s.events[i] = e
			return nil
		}
	}

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i, v := range s.events {
		if v.ID == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return nil
	}
	s.events = append(s.events[:idx], s.events[idx+1:]...)

	return nil
}

func (s *Storage) Close() error {
	return nil
}
