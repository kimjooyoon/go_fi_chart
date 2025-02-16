package domain

import (
	"context"
	"time"
)

// Event 시스템에서 발생하는 이벤트를 나타냅니다.
type Event interface {
	// EventType 이벤트의 타입을 반환합니다.
	EventType() string
	// Timestamp 이벤트가 발생한 시간을 반환합니다.
	Timestamp() time.Time
	// Payload 이벤트의 페이로드를 반환합니다.
	Payload() interface{}
	// Source 이벤트의 발생 소스를 반환합니다.
	Source() string
	// Metadata 이벤트의 메타데이터를 반환합니다.
	Metadata() map[string]string
}

// Handler 이벤트를 처리하는 핸들러입니다.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Publisher 이벤트를 발행하는 발행자입니다.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(handler Handler) error
	Unsubscribe(handler Handler) error
}

// BaseEvent 기본 이벤트 구현체입니다.
type BaseEvent struct {
	Type string
	Time time.Time
	Data interface{}
	Src  string
	Meta map[string]string
}

func (e *BaseEvent) EventType() string {
	return e.Type
}

func (e *BaseEvent) Timestamp() time.Time {
	return e.Time
}

func (e *BaseEvent) Payload() interface{} {
	return e.Data
}

func (e *BaseEvent) Source() string {
	return e.Src
}

func (e *BaseEvent) Metadata() map[string]string {
	return e.Meta
}

// MetricType 메트릭의 타입을 나타냅니다.
type MetricType string

const (
	MetricTypeCounter   MetricType = "counter"
	MetricTypeGauge     MetricType = "gauge"
	MetricTypeHistogram MetricType = "histogram"
	MetricTypeSummary   MetricType = "summary"
)

// Metric 모니터링 시스템의 메트릭을 나타냅니다.
type Metric struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description"`
}

// MetricPayload 메트릭 수집 이벤트의 페이로드입니다.
type MetricPayload struct {
	Metrics []Metric `json:"metrics"`
}

// MonitoringEvent 모니터링 시스템의 이벤트입니다.
type MonitoringEvent struct {
	BaseEvent
}

// NewMonitoringEvent 새로운 모니터링 이벤트를 생성합니다.
func NewMonitoringEvent(eventType string, source string, payload interface{}, metadata map[string]string) Event {
	return &MonitoringEvent{
		BaseEvent: BaseEvent{
			Type: eventType,
			Time: time.Now(),
			Data: payload,
			Src:  source,
			Meta: metadata,
		},
	}
}

// 이벤트 타입 상수
const (
	TypeMetricCollected = "metric.collected"
	TypeAlertTriggered  = "alert.triggered"
)
