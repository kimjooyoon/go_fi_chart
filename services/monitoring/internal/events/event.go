package events

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
)

// EventType 이벤트의 종류를 나타냅니다.
type EventType string

const (
	TypeMetricCollected EventType = "METRIC_COLLECTED"
	TypeAlertTriggered  EventType = "ALERT_TRIGGERED"
)

// Event 모니터링 시스템의 이벤트를 나타냅니다.
type Event struct {
	Type      EventType         `json:"type"`
	Source    string            `json:"source"`
	Timestamp time.Time         `json:"timestamp"`
	Payload   interface{}       `json:"payload"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// MetricPayload 메트릭 수집 이벤트의 페이로드입니다.
type MetricPayload struct {
	Metrics []metrics.Metric `json:"metrics"`
}

// Handler 이벤트를 처리하는 핸들러입니다.
type Handler interface {
	Handle(ctx context.Context, event Event) error
}

// Publisher 이벤트를 발행하는 인터페이스입니다.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
	Subscribe(handler Handler) error
	Unsubscribe(handler Handler) error
}

// SimplePublisher 기본적인 이벤트 발행자 구현체입니다.
type SimplePublisher struct {
	mu       sync.RWMutex
	handlers []Handler
}

// NewSimplePublisher 새로운 SimplePublisher를 생성합니다.
func NewSimplePublisher() *SimplePublisher {
	return &SimplePublisher{
		handlers: make([]Handler, 0),
	}
}

// Publish 이벤트를 발행합니다.
func (p *SimplePublisher) Publish(ctx context.Context, event Event) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, handler := range p.handlers {
		if err := handler.Handle(ctx, event); err != nil {
			// 에러가 발생해도 다른 핸들러는 계속 실행
			continue
		}
	}
	return nil
}

// Subscribe 이벤트 핸들러를 등록합니다.
func (p *SimplePublisher) Subscribe(handler Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, handler)
	return nil
}

// Unsubscribe 이벤트 핸들러를 제거합니다.
func (p *SimplePublisher) Unsubscribe(handler Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i, h := range p.handlers {
		if h == handler {
			p.handlers = append(p.handlers[:i], p.handlers[i+1:]...)
			break
		}
	}
	return nil
}
