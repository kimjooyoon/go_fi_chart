package memory

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/internal/domain/event"
)

// EventStore 인메모리 이벤트 저장소 구현체
type EventStore struct {
	events map[string][]event.Event
	mu     sync.RWMutex
}

// NewEventStore 새로운 인메모리 이벤트 저장소를 생성합니다.
func NewEventStore() *EventStore {
	return &EventStore{
		events: make(map[string][]event.Event),
	}
}

// Save 이벤트를 저장합니다.
func (s *EventStore) Save(_ context.Context, events ...event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, evt := range events {
		aggregateID := evt.AggregateID()
		s.events[aggregateID] = append(s.events[aggregateID], evt)
	}

	return nil
}

// Load 특정 애그리게잇의 이벤트들을 로드합니다.
func (s *EventStore) Load(_ context.Context, aggregateID string) ([]event.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, exists := s.events[aggregateID]
	if !exists {
		return []event.Event{}, nil
	}

	result := make([]event.Event, len(events))
	copy(result, events)
	return result, nil
}

// Clear 모든 이벤트를 삭제합니다.
func (s *EventStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = make(map[string][]event.Event)
}
