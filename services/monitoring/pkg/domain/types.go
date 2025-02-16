package domain

import (
	"context"
	"time"

	metricsdomain "github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
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
	Subscribe(handler Handler) error
	Unsubscribe(handler Handler) error
}

// MetricType 메트릭의 타입을 나타냅니다.
type MetricType = metricsdomain.Type

const (
	TypeCounter   = metricsdomain.TypeCounter
	TypeGauge     = metricsdomain.TypeGauge
	TypeHistogram = metricsdomain.TypeHistogram
	TypeSummary   = metricsdomain.TypeSummary
)

// MetricValue 메트릭의 값을 나타냅니다.
type MetricValue = metricsdomain.Value

// NewMetricValue 새로운 메트릭 값을 생성합니다.
func NewMetricValue(raw float64, labels map[string]string) MetricValue {
	return metricsdomain.NewValue(raw, labels)
}

// Metric 메트릭 인터페이스입니다.
type Metric = metricsdomain.Metric

// Collector 메트릭 수집기 인터페이스입니다.
type Collector = metricsdomain.Collector

// AlertLevel 알림의 심각도를 나타냅니다.
type AlertLevel string

const (
	LevelInfo     AlertLevel = "INFO"
	LevelWarning  AlertLevel = "WARNING"
	LevelError    AlertLevel = "ERROR"
	LevelCritical AlertLevel = "CRITICAL"
)

// Alert 알림을 나타냅니다.
type Alert struct {
	ID        string
	Level     AlertLevel
	Source    string
	Message   string
	Timestamp time.Time
	Metadata  map[string]string
}

// NewAlert 새로운 알림을 생성합니다.
func NewAlert(id, source, message string, level AlertLevel, metadata map[string]string) Alert {
	return Alert{
		ID:        id,
		Level:     level,
		Source:    source,
		Message:   message,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}
}

// Notifier 알림 처리자 인터페이스입니다.
type Notifier interface {
	Notify(ctx context.Context, alert Alert) error
}
