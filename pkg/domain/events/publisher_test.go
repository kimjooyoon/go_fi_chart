package events

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	eventType string
	handled   bool
	err       error
}

func (h *mockHandler) HandleEvent(_ context.Context, _ Event) error {
	h.handled = true
	return h.err
}

func (h *mockHandler) HandlerType() string {
	return h.eventType
}

func TestNewSimplePublisher(t *testing.T) {
	publisher := NewSimplePublisher()

	assert.NotNil(t, publisher)
	assert.NotNil(t, publisher.handlers)
	assert.Empty(t, publisher.handlers)
}

func TestSimplePublisher_RegisterHandler(t *testing.T) {
	publisher := NewSimplePublisher()
	handler := &mockHandler{eventType: "test.event"}

	err := publisher.RegisterHandler(handler)
	assert.NoError(t, err)
	assert.Len(t, publisher.handlers[handler.eventType], 1)
	assert.Equal(t, handler, publisher.handlers[handler.eventType][0])
}

func TestSimplePublisher_UnregisterHandler(t *testing.T) {
	publisher := NewSimplePublisher()
	handler := &mockHandler{eventType: "test.event"}

	publisher.RegisterHandler(handler)
	err := publisher.UnregisterHandler(handler)
	assert.NoError(t, err)
	assert.Empty(t, publisher.handlers[handler.eventType])
}

func TestSimplePublisher_Publish(t *testing.T) {
	publisher := NewSimplePublisher()
	handler1 := &mockHandler{eventType: "test.event"}
	handler2 := &mockHandler{eventType: "test.event"}

	publisher.RegisterHandler(handler1)
	publisher.RegisterHandler(handler2)

	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)
	err := publisher.Publish(context.Background(), event)

	assert.NoError(t, err)
	assert.True(t, handler1.handled)
	assert.True(t, handler2.handled)
}

func TestSimplePublisher_PublishWithError(t *testing.T) {
	publisher := NewSimplePublisher()
	expectedErr := errors.New("handler error")
	handler := &mockHandler{eventType: "test.event", err: expectedErr}

	publisher.RegisterHandler(handler)

	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)
	err := publisher.Publish(context.Background(), event)

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.True(t, handler.handled)
}

func TestSimplePublisher_PublishWithNoHandlers(t *testing.T) {
	publisher := NewSimplePublisher()
	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)

	err := publisher.Publish(context.Background(), event)
	assert.NoError(t, err)
}

type TestEvent struct {
	data string
}

func (e *TestEvent) Name() string {
	return "test_event"
}

func TestEventPublisher_SimplePublish(t *testing.T) {
	publisher := NewSimplePublisher()

	handler := &mockHandler{eventType: "test_event"}
	publisher.RegisterHandler(handler)

	event := NewEvent("test_event", uuid.New(), "test.aggregate", 1, nil, nil)
	err := publisher.Publish(context.Background(), event)

	assert.NoError(t, err)
	assert.True(t, handler.handled)
}
