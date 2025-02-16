package events

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

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

// MonitoringEvent 모니터링 시스템의 이벤트입니다.
type MonitoringEvent struct {
	BaseEvent
}

// NewMonitoringEvent 새로운 모니터링 이벤트를 생성합니다.
func NewMonitoringEvent(eventType string, source string, payload interface{}, metadata map[string]string) domain.Event {
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

// SimplePublisher 기본적인 이벤트 발행자 구현체입니다.
type SimplePublisher struct {
	mu       sync.RWMutex
	handlers []domain.Handler
}

// NewSimplePublisher 새로운 SimplePublisher를 생성합니다.
func NewSimplePublisher() *SimplePublisher {
	return &SimplePublisher{
		handlers: make([]domain.Handler, 0),
	}
}

// Publish 이벤트를 발행합니다.
func (p *SimplePublisher) Publish(ctx context.Context, event domain.Event) error {
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
func (p *SimplePublisher) Subscribe(handler domain.Handler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = append(p.handlers, handler)
	return nil
}

// Unsubscribe 이벤트 핸들러를 제거합니다.
func (p *SimplePublisher) Unsubscribe(handler domain.Handler) error {
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
