package events

import (
	"encoding/json"
	"fmt"
	"sync"
)

// EventSchema는 이벤트의 스키마를 정의합니다.
type EventSchema struct {
	Version     uint
	Name        string
	Description string
	Fields      map[string]FieldSchema
}

// FieldSchema는 이벤트 필드의 스키마를 정의합니다.
type FieldSchema struct {
	Type        string
	Required    bool
	Description string
}

// SchemaRegistry는 이벤트 스키마를 관리하는 레지스트리입니다.
type SchemaRegistry interface {
	// RegisterSchema는 새로운 이벤트 스키마를 등록합니다.
	RegisterSchema(eventType string, schema EventSchema) error

	// GetSchema는 특정 이벤트 타입과 버전의 스키마를 반환합니다.
	GetSchema(eventType string, version uint) (EventSchema, error)

	// ValidateEvent는 이벤트가 스키마에 맞는지 검증합니다.
	ValidateEvent(event Event) error
}

// SimpleSchemaRegistry는 SchemaRegistry의 기본 구현을 제공합니다.
type SimpleSchemaRegistry struct {
	schemas map[string]map[uint]EventSchema
	mu      sync.RWMutex
}

// NewSimpleSchemaRegistry는 새로운 SimpleSchemaRegistry를 생성합니다.
func NewSimpleSchemaRegistry() *SimpleSchemaRegistry {
	return &SimpleSchemaRegistry{
		schemas: make(map[string]map[uint]EventSchema),
	}
}

// RegisterSchema는 새로운 이벤트 스키마를 등록합니다.
func (r *SimpleSchemaRegistry) RegisterSchema(eventType string, schema EventSchema) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.schemas[eventType]; !exists {
		r.schemas[eventType] = make(map[uint]EventSchema)
	}

	if _, exists := r.schemas[eventType][schema.Version]; exists {
		return fmt.Errorf("schema already exists for event type %s version %d", eventType, schema.Version)
	}

	r.schemas[eventType][schema.Version] = schema
	return nil
}

// GetSchema는 특정 이벤트 타입과 버전의 스키마를 반환합니다.
func (r *SimpleSchemaRegistry) GetSchema(eventType string, version uint) (EventSchema, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if versions, exists := r.schemas[eventType]; exists {
		if schema, exists := versions[version]; exists {
			return schema, nil
		}
	}

	return EventSchema{}, fmt.Errorf("schema not found for event type %s version %d", eventType, version)
}

// ValidateEvent는 이벤트가 스키마에 맞는지 검증합니다.
func (r *SimpleSchemaRegistry) ValidateEvent(event Event) error {
	schema, err := r.GetSchema(event.EventType(), event.Version())
	if err != nil {
		return err
	}

	// 페이로드를 JSON으로 변환하여 검증
	payload := event.Payload()
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	var payloadMap map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// 필수 필드 검증
	for fieldName, fieldSchema := range schema.Fields {
		if fieldSchema.Required {
			if _, exists := payloadMap[fieldName]; !exists {
				return fmt.Errorf("required field %s is missing", fieldName)
			}
		}
	}

	return nil
}

// EventUpgrader는 이벤트 버전을 업그레이드하는 인터페이스입니다.
type EventUpgrader interface {
	// UpgradeEvent는 이벤트를 최신 버전으로 업그레이드합니다.
	UpgradeEvent(event Event) (Event, error)
}

// SimpleEventUpgrader는 EventUpgrader의 기본 구현을 제공합니다.
type SimpleEventUpgrader struct {
	registry SchemaRegistry
	upgrades map[string]map[uint]UpgradeFunc
	mu       sync.RWMutex
}

// UpgradeFunc는 이벤트를 다음 버전으로 업그레이드하는 함수입니다.
type UpgradeFunc func(event Event) (Event, error)

// NewSimpleEventUpgrader는 새로운 SimpleEventUpgrader를 생성합니다.
func NewSimpleEventUpgrader(registry SchemaRegistry) *SimpleEventUpgrader {
	return &SimpleEventUpgrader{
		registry: registry,
		upgrades: make(map[string]map[uint]UpgradeFunc),
	}
}

// RegisterUpgrade는 특정 버전의 업그레이드 함수를 등록합니다.
func (u *SimpleEventUpgrader) RegisterUpgrade(eventType string, fromVersion uint, upgrade UpgradeFunc) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, exists := u.upgrades[eventType]; !exists {
		u.upgrades[eventType] = make(map[uint]UpgradeFunc)
	}

	u.upgrades[eventType][fromVersion] = upgrade
	return nil
}

// UpgradeEvent는 이벤트를 최신 버전으로 업그레이드합니다.
func (u *SimpleEventUpgrader) UpgradeEvent(event Event) (Event, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	currentEvent := event
	eventType := event.EventType()
	currentVersion := event.Version()

	for {
		if upgrade, exists := u.upgrades[eventType][currentVersion]; exists {
			var err error
			currentEvent, err = upgrade(currentEvent)
			if err != nil {
				return nil, fmt.Errorf("failed to upgrade event from version %d: %w", currentVersion, err)
			}
			currentVersion = currentEvent.Version()
		} else {
			break
		}
	}

	return currentEvent, nil
}
