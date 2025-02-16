package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/internal/domain/event"
	"github.com/stretchr/testify/assert"
)

func Test_NewEventStore_should_create_empty_store(t *testing.T) {
	// When
	store := NewEventStore()

	// Then
	assert.NotNil(t, store)
	assert.Empty(t, store.events)
}

func Test_EventStore_should_save_and_load_events(t *testing.T) {
	// Given
	store := NewEventStore()
	aggregateID := "test-asset-1"
	events := []event.Event{
		event.NewEvent(
			event.TypeAssetCreated,
			aggregateID,
			"asset",
			map[string]interface{}{"name": "Test Asset"},
			map[string]string{"user": "test-user"},
			1,
		),
		event.NewEvent(
			event.TypeAssetUpdated,
			aggregateID,
			"asset",
			map[string]interface{}{"amount": 1000.0},
			map[string]string{"user": "test-user"},
			2,
		),
	}

	// When
	err := store.Save(context.Background(), events...)
	assert.NoError(t, err)

	loaded, err := store.Load(context.Background(), aggregateID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, loaded, 2)
	assert.Equal(t, events[0], loaded[0])
	assert.Equal(t, events[1], loaded[1])
}

func Test_EventStore_should_return_empty_slice_for_unknown_aggregate(t *testing.T) {
	// Given
	store := NewEventStore()

	// When
	events, err := store.Load(context.Background(), "unknown")

	// Then
	assert.NoError(t, err)
	assert.Empty(t, events)
}

func Test_EventStore_should_clear_all_events(t *testing.T) {
	// Given
	store := NewEventStore()
	event := event.NewEvent(
		event.TypeAssetCreated,
		"test-asset-1",
		"asset",
		map[string]interface{}{"name": "Test Asset"},
		map[string]string{"user": "test-user"},
		1,
	)

	err := store.Save(context.Background(), event)
	assert.NoError(t, err)

	// When
	store.Clear()

	// Then
	assert.Empty(t, store.events)
}

func Test_EventStore_should_be_thread_safe(t *testing.T) {
	// Given
	store := NewEventStore()
	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(2)

	// When
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = store.Save(context.Background(), event.NewEvent(
				event.TypeAssetCreated,
				"test-asset-1",
				"asset",
				nil,
				nil,
				i,
			))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_, _ = store.Load(context.Background(), "test-asset-1")
		}
	}()

	// Then
	wg.Wait()

	events, err := store.Load(context.Background(), "test-asset-1")
	assert.NoError(t, err)
	assert.Len(t, events, iterations)
}
