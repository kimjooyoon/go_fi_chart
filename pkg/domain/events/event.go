package events

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Event는 도메인 이벤트의 기본 인터페이스입니다.
type Event interface {
	// EventID는 이벤트의 고유 식별자를 반환합니다.
	EventID() uuid.UUID

	// EventType은 이벤트의 타입을 반환합니다.
	EventType() string

	// AggregateID는 이벤트가 속한 애그리게잇의 식별자를 반환합니다.
	AggregateID() uuid.UUID

	// AggregateType은 이벤트가 속한 애그리게잇의 타입을 반환합니다.
	AggregateType() string

	// OccurredAt은 이벤트가 발생한 시간을 반환합니다.
	OccurredAt() time.Time

	// Version은 이벤트의 버전을 반환합니다.
	Version() uint

	// Metadata는 이벤트의 메타데이터를 반환합니다.
	Metadata() map[string]interface{}

	// Payload는 이벤트의 실제 데이터를 반환합니다.
	Payload() interface{}
}

// EventHandler는 이벤트를 처리하는 핸들러의 인터페이스입니다.
type EventHandler interface {
	// HandleEvent는 이벤트를 처리합니다.
	HandleEvent(ctx context.Context, event Event) error

	// HandlerType은 핸들러가 처리할 수 있는 이벤트 타입을 반환합니다.
	HandlerType() string
}

// EventBus는 이벤트를 발행하고 구독하는 버스의 인터페이스입니다.
type EventBus interface {
	// Publish는 이벤트를 발행합니다.
	Publish(ctx context.Context, event Event) error

	// Subscribe는 특정 타입의 이벤트를 구독합니다.
	Subscribe(eventType string, handler EventHandler) error

	// Unsubscribe는 특정 타입의 이벤트 구독을 취소합니다.
	Unsubscribe(eventType string, handler EventHandler) error

	// Close는 이벤트 버스를 종료합니다.
	Close() error
}

// EventStore는 이벤트를 저장하고 조회하는 저장소의 인터페이스입니다.
type EventStore interface {
	// Save는 이벤트를 저장합니다.
	Save(ctx context.Context, event Event) error

	// Load는 특정 애그리게잇의 모든 이벤트를 로드합니다.
	Load(ctx context.Context, aggregateID uuid.UUID) ([]Event, error)

	// LoadByType은 특정 타입의 모든 이벤트를 로드합니다.
	LoadByType(ctx context.Context, eventType string) ([]Event, error)
}

// BaseEvent는 Event 인터페이스의 기본 구현을 제공합니다.
type BaseEvent struct {
	ID            uuid.UUID
	Type          string
	AggrID        uuid.UUID
	AggrType      string
	Timestamp     time.Time
	EventVersion  uint
	EventMetadata map[string]interface{}
	EventPayload  interface{}
}

func (e *BaseEvent) EventID() uuid.UUID               { return e.ID }
func (e *BaseEvent) EventType() string                { return e.Type }
func (e *BaseEvent) AggregateID() uuid.UUID           { return e.AggrID }
func (e *BaseEvent) AggregateType() string            { return e.AggrType }
func (e *BaseEvent) OccurredAt() time.Time            { return e.Timestamp }
func (e *BaseEvent) Version() uint                    { return e.EventVersion }
func (e *BaseEvent) Metadata() map[string]interface{} { return e.EventMetadata }
func (e *BaseEvent) Payload() interface{}             { return e.EventPayload }

// NewEvent는 새로운 BaseEvent를 생성합니다.
func NewEvent(
	eventType string,
	aggregateID uuid.UUID,
	aggregateType string,
	version uint,
	payload interface{},
	metadata map[string]interface{},
) Event {
	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &BaseEvent{
		ID:            uuid.New(),
		Type:          eventType,
		AggrID:        aggregateID,
		AggrType:      aggregateType,
		Timestamp:     time.Now().UTC(),
		EventVersion:  version,
		EventMetadata: metadata,
		EventPayload:  payload,
	}
}
