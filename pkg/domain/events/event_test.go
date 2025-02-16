package events

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	eventType := "test.event"
	aggregateID := uuid.New()
	aggregateType := "test.aggregate"
	version := uint(1)
	payload := map[string]interface{}{"key": "value"}
	metadata := map[string]interface{}{"meta": "data"}

	event := NewEvent(eventType, aggregateID, aggregateType, version, payload, metadata)

	assert.NotNil(t, event)
	assert.Equal(t, eventType, event.EventType())
	assert.Equal(t, aggregateID, event.AggregateID())
	assert.Equal(t, aggregateType, event.AggregateType())
	assert.Equal(t, version, event.Version())
	assert.Equal(t, payload, event.Payload())
	assert.Equal(t, metadata, event.Metadata())
	assert.NotZero(t, event.EventID())
	assert.NotZero(t, event.OccurredAt())
}

func TestNewEvent_WithNilMetadata(t *testing.T) {
	eventType := "test.event"
	aggregateID := uuid.New()
	aggregateType := "test.aggregate"
	version := uint(1)
	payload := map[string]interface{}{"key": "value"}

	event := NewEvent(eventType, aggregateID, aggregateType, version, payload, nil)

	assert.NotNil(t, event)
	assert.NotNil(t, event.Metadata())
	assert.Empty(t, event.Metadata())
}

func TestBaseEvent_Methods(t *testing.T) {
	id := uuid.New()
	eventType := "test.event"
	aggregateID := uuid.New()
	aggregateType := "test.aggregate"
	timestamp := time.Now()
	version := uint(1)
	metadata := map[string]interface{}{"meta": "data"}
	payload := map[string]interface{}{"key": "value"}

	event := &BaseEvent{
		ID:            id,
		Type:          eventType,
		AggrID:        aggregateID,
		AggrType:      aggregateType,
		Timestamp:     timestamp,
		EventVersion:  version,
		EventMetadata: metadata,
		EventPayload:  payload,
	}

	assert.Equal(t, id, event.EventID())
	assert.Equal(t, eventType, event.EventType())
	assert.Equal(t, aggregateID, event.AggregateID())
	assert.Equal(t, aggregateType, event.AggregateType())
	assert.Equal(t, timestamp, event.OccurredAt())
	assert.Equal(t, version, event.Version())
	assert.Equal(t, metadata, event.Metadata())
	assert.Equal(t, payload, event.Payload())
}
