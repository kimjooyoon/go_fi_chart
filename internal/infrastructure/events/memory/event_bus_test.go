package memory

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/internal/domain/event"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	name     string
	events   []event.Event
	mu       sync.Mutex
	handleFn func(event.Event) error
}

func newMockHandler(name string) *mockHandler {
	return &mockHandler{
		name:   name,
		events: make([]event.Event, 0),
	}
}

func (h *mockHandler) HandleEvent(_ context.Context, evt event.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.events = append(h.events, evt)
	if h.handleFn != nil {
		return h.handleFn(evt)
	}
	return nil
}

func (h *mockHandler) HandlerName() string {
	return h.name
}

func Test_NewEventBus_should_create_empty_bus(t *testing.T) {
	// When
	bus := NewEventBus()

	// Then
	assert.NotNil(t, bus)
	assert.Empty(t, bus.handlers)
	assert.False(t, bus.closed)
}

func Test_EventBus_should_publish_events_to_handlers(t *testing.T) {
	// Given
	bus := NewEventBus()
	handler1 := newMockHandler("handler1")
	handler2 := newMockHandler("handler2")

	err := bus.Subscribe(handler1)
	assert.NoError(t, err)
	err = bus.Subscribe(handler2)
	assert.NoError(t, err)

	evt := event.NewEvent(
		event.TypeAssetCreated,
		"test-asset-1",
		"asset",
		map[string]interface{}{"name": "Test Asset"},
		map[string]string{"user": "test-user"},
		1,
	)

	// When
	err = bus.Publish(context.Background(), evt)

	// Then
	assert.NoError(t, err)
	assert.Len(t, handler1.events, 1)
	assert.Len(t, handler2.events, 1)
	assert.Equal(t, evt, handler1.events[0])
	assert.Equal(t, evt, handler2.events[0])
}

func Test_EventBus_should_handle_handler_errors(t *testing.T) {
	// Given
	bus := NewEventBus()
	handler := newMockHandler("handler")
	handler.handleFn = func(event.Event) error {
		return assert.AnError
	}

	err := bus.Subscribe(handler)
	assert.NoError(t, err)

	evt := event.NewEvent(
		event.TypeAssetCreated,
		"test-asset-1",
		"asset",
		map[string]interface{}{"name": "Test Asset"},
		map[string]string{"user": "test-user"},
		1,
	)

	// When
	err = bus.Publish(context.Background(), evt)

	// Then
	assert.NoError(t, err) // 핸들러 에러가 발생해도 발행은 성공
	assert.Len(t, handler.events, 1)
}

func Test_EventBus_should_unsubscribe_handler(t *testing.T) {
	// Given
	bus := NewEventBus()
	handler := newMockHandler("handler")

	err := bus.Subscribe(handler)
	assert.NoError(t, err)

	// When
	err = bus.Unsubscribe(handler)

	// Then
	assert.NoError(t, err)
	assert.Empty(t, bus.handlers[handler.HandlerName()])
}

func Test_EventBus_should_close(t *testing.T) {
	// Given
	bus := NewEventBus()
	handler := newMockHandler("handler")

	err := bus.Subscribe(handler)
	assert.NoError(t, err)

	// When
	err = bus.Close()

	// Then
	assert.NoError(t, err)
	assert.True(t, bus.closed)
	assert.Nil(t, bus.handlers)

	// 닫힌 후에는 작업이 실패해야 함
	err = bus.Subscribe(handler)
	assert.Equal(t, event.ErrEventBusClosed, err)

	err = bus.Unsubscribe(handler)
	assert.Equal(t, event.ErrEventBusClosed, err)

	err = bus.Publish(context.Background(), event.NewEvent(
		event.TypeAssetCreated,
		"test-asset-1",
		"asset",
		nil,
		nil,
		1,
	))
	assert.Equal(t, event.ErrEventBusClosed, err)
}

func Test_EventBus_should_be_thread_safe(t *testing.T) {
	// Given
	bus := NewEventBus()
	handler := newMockHandler("handler")
	err := bus.Subscribe(handler)
	assert.NoError(t, err)

	iterations := 1000
	wg := sync.WaitGroup{}
	wg.Add(3)

	// When
	go func() {
		defer wg.Done()
		for i := 0; i < iterations; i++ {
			_ = bus.Publish(context.Background(), event.NewEvent(
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
		for i := 0; i < iterations/2; i++ {
			_ = bus.Subscribe(newMockHandler("handler" + string(rune(i))))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < iterations/2; i++ {
			_ = bus.Unsubscribe(handler)
			_ = bus.Subscribe(handler)
		}
	}()

	// Then
	wg.Wait()
}
