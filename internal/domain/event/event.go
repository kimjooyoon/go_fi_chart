package event

import (
	"time"
)

// Event 도메인 이벤트 인터페이스
type Event interface {
	GetEventType() string
	GetAggregateType() string
	GetPayload() map[string]interface{}
	GetTimestamp() time.Time
}

// BaseEvent 기본 이벤트 구현체
type BaseEvent struct {
	EventType     string
	AggregateType string
	Payload       map[string]interface{}
	Timestamp     time.Time
}

// NewEvent 새로운 이벤트를 생성합니다.
func NewEvent(eventType string, aggregateType string, payload map[string]interface{}) Event {
	return &BaseEvent{
		EventType:     eventType,
		AggregateType: aggregateType,
		Payload:       payload,
		Timestamp:     time.Now(),
	}
}

// GetEventType 이벤트 타입을 반환합니다.
func (e *BaseEvent) GetEventType() string {
	return e.EventType
}

// GetAggregateType 집계 타입을 반환합니다.
func (e *BaseEvent) GetAggregateType() string {
	return e.AggregateType
}

// GetPayload 이벤트 페이로드를 반환합니다.
func (e *BaseEvent) GetPayload() map[string]interface{} {
	return e.Payload
}

// GetTimestamp 이벤트 발생 시간을 반환합니다.
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}
