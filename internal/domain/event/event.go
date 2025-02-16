package event

import (
	"context"
	"time"
)

// Type 이벤트의 타입을 나타냅니다.
type Type string

// 시스템에서 사용되는 이벤트 타입 상수들
const (
	TypeAssetCreated        Type = "asset.created"
	TypeAssetUpdated        Type = "asset.updated"
	TypeAssetDeleted        Type = "asset.deleted"
	TypeTransactionRecorded Type = "transaction.recorded"
	TypePortfolioRebalanced Type = "portfolio.rebalanced"
	TypeMetricCollected     Type = "metric.collected"
	TypeAlertTriggered      Type = "alert.triggered"
)

// Event 도메인 이벤트 인터페이스
type Event interface {
	// EventType 이벤트의 타입을 반환합니다.
	EventType() Type
	// AggregateID 이벤트가 발생한 애그리게잇의 ID를 반환합니다.
	AggregateID() string
	// AggregateType 이벤트가 발생한 애그리게잇의 타입을 반환합니다.
	AggregateType() string
	// Payload 이벤트의 페이로드를 반환합니다.
	Payload() interface{}
	// Metadata 이벤트의 메타데이터를 반환합니다.
	Metadata() map[string]string
	// Timestamp 이벤트가 발생한 시간을 반환합니다.
	Timestamp() time.Time
	// Version 이벤트의 버전을 반환합니다.
	Version() int
}

// Handler 이벤트를 처리하는 핸들러 인터페이스
type Handler interface {
	// HandleEvent 이벤트를 처리합니다.
	HandleEvent(ctx context.Context, event Event) error
	// HandlerName 핸들러의 이름을 반환합니다.
	HandlerName() string
}

// Bus 이벤트 버스 인터페이스입니다.
type Bus interface {
	// Publish 이벤트를 발행합니다.
	Publish(ctx context.Context, event Event) error
	// Subscribe 이벤트 핸들러를 등록합니다.
	Subscribe(handler Handler) error
	// Unsubscribe 이벤트 핸들러를 제거합니다.
	Unsubscribe(handler Handler) error
}

// Store 이벤트를 저장하는 저장소 인터페이스
type Store interface {
	// Save 이벤트를 저장합니다.
	Save(ctx context.Context, events ...Event) error
	// Load 특정 애그리게잇의 이벤트들을 로드합니다.
	Load(ctx context.Context, aggregateID string) ([]Event, error)
}

// BaseEvent 기본 이벤트 구현체
type BaseEvent struct {
	eventType     Type
	aggregateID   string
	aggregateType string
	payload       interface{}
	metadata      map[string]string
	timestamp     time.Time
	version       int
}

// NewEvent 새로운 이벤트를 생성합니다.
func NewEvent(
	eventType Type,
	aggregateID string,
	aggregateType string,
	payload interface{},
	metadata map[string]string,
	version int,
) Event {
	return &BaseEvent{
		eventType:     eventType,
		aggregateID:   aggregateID,
		aggregateType: aggregateType,
		payload:       payload,
		metadata:      metadata,
		timestamp:     time.Now(),
		version:       version,
	}
}

func (e *BaseEvent) EventType() Type {
	return e.eventType
}

func (e *BaseEvent) AggregateID() string {
	return e.aggregateID
}

func (e *BaseEvent) AggregateType() string {
	return e.aggregateType
}

func (e *BaseEvent) Payload() interface{} {
	return e.payload
}

func (e *BaseEvent) Metadata() map[string]string {
	return e.metadata
}

func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BaseEvent) Version() int {
	return e.version
}

// SimpleBus 기본 이벤트 버스 구현체
type SimpleBus struct {
	handlers []Handler
}

// NewSimpleBus 새로운 SimpleBus를 생성합니다.
func NewSimpleBus() *SimpleBus {
	return &SimpleBus{
		handlers: make([]Handler, 0),
	}
}
