package event

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewEvent_should_create_event_with_valid_data(t *testing.T) {
	// Given
	eventType := TypeAssetCreated
	aggregateID := "test-asset-1"
	aggregateType := "asset"
	payload := map[string]interface{}{
		"name":   "Test Asset",
		"amount": 1000.0,
	}
	metadata := map[string]string{
		"user": "test-user",
	}
	version := 1

	// When
	event := NewEvent(eventType, aggregateID, aggregateType, payload, metadata, version)

	// Then
	assert.NotNil(t, event)
	assert.Equal(t, eventType, event.EventType())
	assert.Equal(t, aggregateID, event.AggregateID())
	assert.Equal(t, aggregateType, event.AggregateType())
	assert.Equal(t, payload, event.Payload())
	assert.Equal(t, metadata, event.Metadata())
	assert.Equal(t, version, event.Version())
	assert.NotZero(t, event.Timestamp())
	assert.True(t, event.Timestamp().Before(time.Now()))
}

func Test_BaseEvent_should_implement_Event_interface(t *testing.T) {
	// Given
	var event Event = &BaseEvent{}

	// Then
	assert.NotNil(t, event)
	assert.Implements(t, (*Event)(nil), event)
}

func Test_BaseEvent_methods_should_return_correct_values(t *testing.T) {
	// Given
	now := time.Now()
	event := &BaseEvent{
		eventType:     TypeAssetCreated,
		aggregateID:   "test-asset-1",
		aggregateType: "asset",
		payload: map[string]interface{}{
			"name":   "Test Asset",
			"amount": 1000.0,
		},
		metadata: map[string]string{
			"user": "test-user",
		},
		timestamp: now,
		version:   1,
	}

	// Then
	assert.Equal(t, TypeAssetCreated, event.EventType())
	assert.Equal(t, "test-asset-1", event.AggregateID())
	assert.Equal(t, "asset", event.AggregateType())
	assert.Equal(t, map[string]interface{}{
		"name":   "Test Asset",
		"amount": 1000.0,
	}, event.Payload())
	assert.Equal(t, map[string]string{
		"user": "test-user",
	}, event.Metadata())
	assert.Equal(t, now, event.Timestamp())
	assert.Equal(t, 1, event.Version())
}
