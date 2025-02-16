package events

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewExponentialBackoff(t *testing.T) {
	backoff := NewExponentialBackoff(
		100*time.Millisecond,
		10*time.Second,
		5,
		2.0,
	)

	assert.NotNil(t, backoff)
	assert.Equal(t, 100*time.Millisecond, backoff.initialInterval)
	assert.Equal(t, 10*time.Second, backoff.maxInterval)
	assert.Equal(t, 5, backoff.maxAttempts)
	assert.Equal(t, 2.0, backoff.multiplier)
}

func TestExponentialBackoff_ShouldRetry(t *testing.T) {
	backoff := NewExponentialBackoff(
		100*time.Millisecond,
		10*time.Second,
		5,
		2.0,
	)

	assert.True(t, backoff.ShouldRetry(errors.New("test error")))
}

func TestExponentialBackoff_NextBackoff(t *testing.T) {
	backoff := NewExponentialBackoff(
		100*time.Millisecond,
		10*time.Second,
		5,
		2.0,
	)

	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(0))
	assert.Equal(t, 100*time.Millisecond, backoff.NextBackoff(1))
	assert.Equal(t, 200*time.Millisecond, backoff.NextBackoff(2))
	assert.Equal(t, 400*time.Millisecond, backoff.NextBackoff(3))
	assert.Equal(t, 800*time.Millisecond, backoff.NextBackoff(4))
	assert.Equal(t, 1600*time.Millisecond, backoff.NextBackoff(5))
}

func TestExponentialBackoff_MaxAttempts(t *testing.T) {
	backoff := NewExponentialBackoff(
		100*time.Millisecond,
		10*time.Second,
		5,
		2.0,
	)

	assert.Equal(t, 5, backoff.MaxAttempts())
}

func TestRetryableEventHandler(t *testing.T) {
	expectedError := errors.New("test error")
	handler := &mockHandler{
		eventType: "test.event",
		err:       expectedError,
	}

	policy := NewExponentialBackoff(
		10*time.Millisecond,
		100*time.Millisecond,
		3,
		2.0,
	)

	retryHandler := NewRetryableEventHandler(handler, policy)
	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)

	err := retryHandler.HandleEvent(context.Background(), event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "max retry attempts (3) reached")
}

func TestRetryableEventHandler_Success(t *testing.T) {
	handler := &mockHandler{
		eventType: "test.event",
		err:       nil,
	}

	policy := NewExponentialBackoff(
		10*time.Millisecond,
		100*time.Millisecond,
		3,
		2.0,
	)

	retryHandler := NewRetryableEventHandler(handler, policy)
	event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)

	err := retryHandler.HandleEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.True(t, handler.handled)
}

func TestRetryableEventHandler_HandlerType(t *testing.T) {
	handler := &mockHandler{eventType: "test.event"}
	policy := DefaultRetryPolicy()
	retryHandler := NewRetryableEventHandler(handler, policy)

	assert.Equal(t, "test.event", retryHandler.HandlerType())
}

func TestDefaultRetryPolicy(t *testing.T) {
	policy := DefaultRetryPolicy()

	assert.NotNil(t, policy)
	assert.Equal(t, 5, policy.MaxAttempts())
	assert.Equal(t, 100*time.Millisecond, policy.NextBackoff(1))
}
