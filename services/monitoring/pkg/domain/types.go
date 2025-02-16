package domain

import (
	"context"
	"time"
)

// EventType 이벤트 타입입니다.
type EventType string

const (
	TypeMetricCollected EventType = "metric_collected"
	TypeAlertTriggered  EventType = "alert_triggered"
)

// Event 모니터링 이벤트입니다.
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// NewMonitoringEvent 새로운 모니터링 이벤트를 생성합니다.
func NewMonitoringEvent(eventType EventType, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// Handler 이벤트 핸들러 인터페이스입니다.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Publisher 이벤트 발행자 인터페이스입니다.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}
