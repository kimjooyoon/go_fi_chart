package events

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleSubscriber(t *testing.T) {
	subscriber := NewSimpleSubscriber()

	assert.NotNil(t, subscriber)
	assert.NotNil(t, subscriber.handlers)
	assert.NotNil(t, subscriber.eventChan)
	assert.Empty(t, subscriber.handlers)
	assert.False(t, subscriber.isRunning)
}

func TestSimpleSubscriber_Subscribe(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	handler := &mockHandler{eventType: "test.event"}

	err := subscriber.Subscribe("test.event", handler)
	assert.NoError(t, err)
	assert.Len(t, subscriber.handlers["test.event"], 1)
	assert.Equal(t, handler, subscriber.handlers["test.event"][0])
}

func TestSimpleSubscriber_Unsubscribe(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	handler := &mockHandler{eventType: "test.event"}

	subscriber.Subscribe("test.event", handler)
	err := subscriber.Unsubscribe("test.event", handler)
	assert.NoError(t, err)
	assert.Empty(t, subscriber.handlers["test.event"])
}

func TestSimpleSubscriber_StartStop(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	ctx := context.Background()

	err := subscriber.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, subscriber.isRunning)

	err = subscriber.Stop()
	assert.NoError(t, err)
	assert.False(t, subscriber.isRunning)
}

func TestSimpleSubscriber_PublishEvent(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	handler := &mockHandler{eventType: "test.event"}
	ctx := context.Background()

	subscriber.Subscribe("test.event", handler)
	subscriber.Start(ctx)
	defer subscriber.Stop()

	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)
	err := subscriber.PublishEvent(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, handler.handled)
}

func TestSimpleSubscriber_PublishEventWhenStopped(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)

	err := subscriber.PublishEvent(event)
	assert.NoError(t, err)
}

func TestSimpleSubscriber_MultipleHandlers(t *testing.T) {
	subscriber := NewSimpleSubscriber()
	handler1 := &mockHandler{eventType: "test.event"}
	handler2 := &mockHandler{eventType: "test.event"}
	ctx := context.Background()

	subscriber.Subscribe("test.event", handler1)
	subscriber.Subscribe("test.event", handler2)
	subscriber.Start(ctx)
	defer subscriber.Stop()

	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)
	err := subscriber.PublishEvent(event)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.True(t, handler1.handled)
	assert.True(t, handler2.handled)
}
