package events

import (
	"context"
	"sync"
)

// EventPublisher는 이벤트 발행을 담당하는 인터페이스입니다.
type EventPublisher interface {
	// Publish는 이벤트를 발행합니다.
	Publish(ctx context.Context, event Event) error

	// RegisterHandler는 이벤트 핸들러를 등록합니다.
	RegisterHandler(handler EventHandler) error

	// UnregisterHandler는 이벤트 핸들러를 제거합니다.
	UnregisterHandler(handler EventHandler) error
}

// SimplePublisher는 EventPublisher의 기본 구현을 제공합니다.
type SimplePublisher struct {
	handlers map[string][]EventHandler
	mu       sync.RWMutex
}

// NewSimplePublisher는 새로운 SimplePublisher를 생성합니다.
func NewSimplePublisher() *SimplePublisher {
	return &SimplePublisher{
		handlers: make(map[string][]EventHandler),
	}
}

// Publish는 등록된 모든 핸들러에게 이벤트를 전달합니다.
func (p *SimplePublisher) Publish(ctx context.Context, event Event) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// 이벤트 타입에 해당하는 핸들러들을 찾습니다
	handlers, exists := p.handlers[event.EventType()]
	if !exists {
		return nil
	}

	// 각 핸들러에게 이벤트를 전달합니다
	var wg sync.WaitGroup
	errCh := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			if err := h.HandleEvent(ctx, event); err != nil {
				errCh <- err
			}
		}(handler)
	}

	// 모든 핸들러의 처리가 완료될 때까지 대기합니다
	wg.Wait()
	close(errCh)

	// 에러가 있다면 첫 번째 에러를 반환합니다
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

// RegisterHandler는 새로운 이벤트 핸들러를 등록합니다.
func (p *SimplePublisher) RegisterHandler(handler EventHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	eventType := handler.HandlerType()
	p.handlers[eventType] = append(p.handlers[eventType], handler)
	return nil
}

// UnregisterHandler는 등록된 이벤트 핸들러를 제거합니다.
func (p *SimplePublisher) UnregisterHandler(handler EventHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	eventType := handler.HandlerType()
	handlers := p.handlers[eventType]

	// 핸들러 목록에서 해당 핸들러를 찾아 제거합니다
	for i, h := range handlers {
		if h == handler {
			p.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// Subscribe는 특정 타입의 이벤트를 구독합니다.
func (p *SimplePublisher) Subscribe(eventType string, handler EventHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers[eventType] = append(p.handlers[eventType], handler)
	return nil
}

// Unsubscribe는 특정 타입의 이벤트 구독을 취소합니다.
func (p *SimplePublisher) Unsubscribe(eventType string, handler EventHandler) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if handlers, exists := p.handlers[eventType]; exists {
		for i, h := range handlers {
			if h == handler {
				p.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (p *SimplePublisher) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.handlers = make(map[string][]EventHandler)
	return nil
}
