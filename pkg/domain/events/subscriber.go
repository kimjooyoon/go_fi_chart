package events

import (
	"context"
	"sync"
)

// EventSubscriber는 이벤트 구독을 담당하는 인터페이스입니다.
type EventSubscriber interface {
	// Subscribe는 특정 이벤트 타입을 구독합니다.
	Subscribe(eventType string, handler EventHandler) error

	// Unsubscribe는 특정 이벤트 타입의 구독을 취소합니다.
	Unsubscribe(eventType string, handler EventHandler) error

	// Start는 구독자를 시작합니다.
	Start(ctx context.Context) error

	// Stop은 구독자를 중지합니다.
	Stop() error
}

// SimpleSubscriber는 EventSubscriber의 기본 구현을 제공합니다.
type SimpleSubscriber struct {
	handlers   map[string][]EventHandler
	eventChan  chan Event
	mu         sync.RWMutex
	isRunning  bool
	cancelFunc context.CancelFunc
}

// NewSimpleSubscriber는 새로운 SimpleSubscriber를 생성합니다.
func NewSimpleSubscriber() *SimpleSubscriber {
	return &SimpleSubscriber{
		handlers:  make(map[string][]EventHandler),
		eventChan: make(chan Event, 100),
	}
}

// Subscribe는 새로운 이벤트 핸들러를 등록합니다.
func (s *SimpleSubscriber) Subscribe(eventType string, handler EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.handlers[eventType] = append(s.handlers[eventType], handler)
	return nil
}

// Unsubscribe는 등록된 이벤트 핸들러를 제거합니다.
func (s *SimpleSubscriber) Unsubscribe(eventType string, handler EventHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	handlers := s.handlers[eventType]
	for i, h := range handlers {
		if h == handler {
			s.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// Start는 이벤트 처리를 시작합니다.
func (s *SimpleSubscriber) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.isRunning {
		s.mu.Unlock()
		return nil
	}

	ctx, cancel := context.WithCancel(ctx)
	s.cancelFunc = cancel
	s.isRunning = true
	s.mu.Unlock()

	go s.processEvents(ctx)
	return nil
}

// Stop은 이벤트 처리를 중지합니다.
func (s *SimpleSubscriber) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	s.cancelFunc()
	s.isRunning = false
	return nil
}

// processEvents는 이벤트 채널로부터 이벤트를 받아 처리합니다.
func (s *SimpleSubscriber) processEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-s.eventChan:
			s.handleEvent(ctx, event)
		}
	}
}

// handleEvent는 단일 이벤트를 처리합니다.
func (s *SimpleSubscriber) handleEvent(ctx context.Context, event Event) {
	s.mu.RLock()
	handlers := s.handlers[event.EventType()]
	s.mu.RUnlock()

	var wg sync.WaitGroup
	for _, handler := range handlers {
		wg.Add(1)
		go func(h EventHandler) {
			defer wg.Done()
			_ = h.HandleEvent(ctx, event)
		}(handler)
	}
	wg.Wait()
}

// PublishEvent는 이벤트를 구독자의 채널로 발행합니다.
func (s *SimpleSubscriber) PublishEvent(event Event) error {
	s.mu.RLock()
	isRunning := s.isRunning
	s.mu.RUnlock()

	if !isRunning {
		return nil
	}

	select {
	case s.eventChan <- event:
		return nil
	default:
		return nil
	}
}
