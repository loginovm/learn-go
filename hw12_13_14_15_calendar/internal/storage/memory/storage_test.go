package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/loginovm/learn-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	sut := New()
	ctx := context.Background()

	e1 := storage.Event{
		ID:      "1",
		StartAt: time.Date(2024, 8, 1, 10, 0, 0, 0, time.Local),
		EndAt:   time.Date(2024, 8, 1, 11, 0, 0, 0, time.Local),
	}
	err := sut.CreateEvent(ctx, e1)
	require.NoError(t, err)

	e2 := storage.Event{
		ID:      "2",
		StartAt: time.Date(2024, 8, 1, 15, 0, 0, 0, time.Local),
		EndAt:   time.Date(2024, 8, 1, 16, 0, 0, 0, time.Local),
	}
	err = sut.CreateEvent(ctx, e2)
	require.NoError(t, err)

	// Get all events that was added
	all, err := sut.GetEvents(ctx)
	require.NoError(t, err)
	expectedLen := 2
	require.Len(t, all, expectedLen)

	// Get event by time conditions
	selected, err := sut.GetEventsInPeriod(ctx,
		time.Date(2024, 8, 1, 15, 30, 0, 0, time.Local),
		time.Date(2024, 8, 1, 16, 0, 0, 0, time.Local))
	require.NoError(t, err)
	expectedLen = 1
	require.Len(t, selected, expectedLen)

	// Get event by ID
	event, err := sut.GetEvent(ctx, e1.ID)
	require.NoError(t, err)
	event.Description = "new"
	err = sut.UpdateEvent(ctx, event)
	require.NoError(t, err)

	// Verify that event was updated
	event, err = sut.GetEvent(ctx, e1.ID)
	require.NoError(t, err)
	require.Equal(t, "new", event.Description)

	// Delete event
	err = sut.DeleteEvent(ctx, e1.ID)
	require.NoError(t, err)

	// Get all events
	all, err = sut.GetEvents(ctx)
	require.NoError(t, err)
	expectedLen = 1
	require.Len(t, all, expectedLen)
}

func TestStorage_Concurrency(t *testing.T) {
	ctx := context.Background()
	sut := New()
	eventid := "1"
	e := storage.Event{ID: eventid}
	err := sut.CreateEvent(ctx, e)
	require.NoError(t, err)

	parallelizm := 10
	wg := sync.WaitGroup{}
	// Update same event from n threads in parallel
	// if storage is not thread safe test with -race flag will fail
	for i := 0; i < parallelizm; i++ {
		wg.Add(1)
		go func() {
			pid := i
			for j := 0; j < 5; j++ {
				event, err := sut.GetEvent(ctx, eventid)
				require.NoError(t, err)
				event.Description = fmt.Sprintf("pid%d_%d", pid, j)
				err = sut.UpdateEvent(ctx, event)
				require.NoError(t, err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
