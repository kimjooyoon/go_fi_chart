package events

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSimpleSchemaRegistry(t *testing.T) {
	registry := NewSimpleSchemaRegistry()

	assert.NotNil(t, registry)
	assert.NotNil(t, registry.schemas)
	assert.Empty(t, registry.schemas)
}

func TestSimpleSchemaRegistry_RegisterSchema(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	schema := EventSchema{
		Version:     1,
		Name:        "TestEvent",
		Description: "Test event schema",
		Fields: map[string]FieldSchema{
			"field1": {
				Type:        "string",
				Required:    true,
				Description: "Test field",
			},
		},
	}

	err := registry.RegisterSchema("test.event", schema)
	assert.NoError(t, err)
	assert.Len(t, registry.schemas["test.event"], 1)
	assert.Equal(t, schema, registry.schemas["test.event"][1])

	err = registry.RegisterSchema("test.event", schema)
	assert.Error(t, err)
}

func TestSimpleSchemaRegistry_GetSchema(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	schema := EventSchema{
		Version:     1,
		Name:        "TestEvent",
		Description: "Test event schema",
		Fields: map[string]FieldSchema{
			"field1": {
				Type:        "string",
				Required:    true,
				Description: "Test field",
			},
		},
	}

	registry.RegisterSchema("test.event", schema)

	found, err := registry.GetSchema("test.event", 1)
	assert.NoError(t, err)
	assert.Equal(t, schema, found)

	_, err = registry.GetSchema("test.event", 2)
	assert.Error(t, err)
}

func TestSimpleSchemaRegistry_ValidateEvent(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	schema := EventSchema{
		Version:     1,
		Name:        "TestEvent",
		Description: "Test event schema",
		Fields: map[string]FieldSchema{
			"field1": {
				Type:        "string",
				Required:    true,
				Description: "Test field",
			},
		},
	}

	registry.RegisterSchema("test.event", schema)

	validPayload := map[string]interface{}{
		"field1": "test value",
	}
	validEvent := NewEvent("test.event", uuid.New(), "test.aggregate", 1, validPayload, nil)

	err := registry.ValidateEvent(validEvent)
	assert.NoError(t, err)

	invalidPayload := map[string]interface{}{}
	invalidEvent := NewEvent("test.event", uuid.New(), "test.aggregate", 1, invalidPayload, nil)

	err = registry.ValidateEvent(invalidEvent)
	assert.Error(t, err)
}

func TestSimpleEventUpgrader(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	upgrader := NewSimpleEventUpgrader(registry)

	assert.NotNil(t, upgrader)
	assert.NotNil(t, upgrader.upgrades)
	assert.Empty(t, upgrader.upgrades)
}

func TestSimpleEventUpgrader_RegisterUpgrade(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	upgrader := NewSimpleEventUpgrader(registry)

	upgrade := func(event Event) (Event, error) {
		return event, nil
	}

	err := upgrader.RegisterUpgrade("test.event", 1, upgrade)
	assert.NoError(t, err)
	assert.Len(t, upgrader.upgrades["test.event"], 1)
	assert.NotNil(t, upgrader.upgrades["test.event"][1])
}

func TestSimpleEventUpgrader_UpgradeEvent(t *testing.T) {
	registry := NewSimpleSchemaRegistry()
	upgrader := NewSimpleEventUpgrader(registry)

	v1Event := NewEvent("test.event", uuid.New(), "test.aggregate", 1, nil, nil)
	v2Payload := map[string]interface{}{"version": 2}

	upgrader.RegisterUpgrade("test.event", 1, func(event Event) (Event, error) {
		return NewEvent(event.EventType(), event.AggregateID(), event.AggregateType(), 2, v2Payload, event.Metadata()), nil
	})

	upgraded, err := upgrader.UpgradeEvent(v1Event)
	assert.NoError(t, err)
	assert.Equal(t, uint(2), upgraded.Version())
	assert.Equal(t, v2Payload, upgraded.Payload())
}
