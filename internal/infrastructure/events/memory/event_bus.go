package memory

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/internal/domain/event"
)

// EventBus 인메모리 이벤트 버스 구현체
type EventBus struct {
	handlers map[string][]event.Handler
	mu       sync.RWMutex
	closed   bool
}

// NewEventBus 새로운 인메모리 이벤트 버스를 생성합니다.
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]event.Handler),
	}
}

// Publish 이벤트를 발행합니다.
func (b *EventBus) Publish(ctx context.Context, evt event.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if b.closed {
		return event.ErrEventBusClosed
	}

	// 모든 핸들러에게 이벤트 전달
	for _, handlers := range b.handlers {
		for _, handler := range handlers {
			if err := handler.HandleEvent(ctx, evt); err != nil {
				// 에러가 발생해도 다른 핸들러는 계속 실행
				continue
			}
		}
	}

	return nil
}

// Subscribe 이벤트 핸들러를 등록합니다.
func (b *EventBus) Subscribe(handler event.Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return event.ErrEventBusClosed
	}

	handlerName := handler.HandlerName()
	b.handlers[handlerName] = append(b.handlers[handlerName], handler)
	return nil
}

// Unsubscribe 이벤트 핸들러를 제거합니다.
func (b *EventBus) Unsubscribe(handler event.Handler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return event.ErrEventBusClosed
	}

	handlerName := handler.HandlerName()
	handlers := b.handlers[handlerName]
	for i, h := range handlers {
		if h == handler {
			b.handlers[handlerName] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
	return nil
}

// Close 이벤트 버스를 종료합니다.
func (b *EventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closed = true
	b.handlers = nil
	return nil
}
