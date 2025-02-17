package internal

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	"github.com/aske/go_fi_chart/services/asset/internal/infrastructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testEventHandler struct {
	events []events.Event
	mu     sync.Mutex
}

func (h *testEventHandler) HandleEvent(_ context.Context, event events.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, event)
	return nil
}

func (h *testEventHandler) GetEvents() []events.Event {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.events
}

func (h *testEventHandler) HandlerType() string {
	return "test.handler"
}

type testEventBus struct {
	handlers map[string][]events.EventHandler
	mu       sync.RWMutex
}

func (b *testEventBus) Close() error {
	return nil
}

func (b *testEventBus) Subscribe(eventType string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.handlers == nil {
		b.handlers = make(map[string][]events.EventHandler)
	}
	b.handlers[eventType] = append(b.handlers[eventType], handler)
	return nil
}

func (b *testEventBus) Unsubscribe(eventType string, handler events.EventHandler) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if handlers, ok := b.handlers[eventType]; ok {
		for i, h := range handlers {
			if h == handler {
				b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
	return nil
}

func (b *testEventBus) Publish(ctx context.Context, event events.Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if handlers, ok := b.handlers[event.EventType()]; ok {
		for _, handler := range handlers {
			if err := handler.HandleEvent(ctx, event); err != nil {
				return err
			}
		}
	}
	return nil
}

func newTestEventBus() *testEventBus {
	return &testEventBus{
		handlers: make(map[string][]events.EventHandler),
	}
}

func TestAssetService_Integration(t *testing.T) {
	// Given
	eventBus := newTestEventBus()
	handler := &testEventHandler{}

	// 모든 이벤트 타입에 대해 핸들러 등록
	eventTypes := []string{
		domain.EventTypeAssetCreated,
		domain.EventTypeAssetUpdated,
		domain.EventTypeAssetAmountChanged,
		domain.EventTypeAssetDeleted,
	}
	for _, eventType := range eventTypes {
		err := eventBus.Subscribe(eventType, handler)
		require.NoError(t, err)
	}

	repo := infrastructure.NewMemoryAssetRepository(eventBus)
	amount, err := valueobjects.NewMoney(1000.0, "USD")
	require.NoError(t, err)

	// When
	asset := domain.NewAsset("user-1", domain.Stock, "Test Asset", amount)
	err = repo.Save(context.Background(), asset)
	require.NoError(t, err)

	// 이벤트 발생 확인
	assert.Eventually(t, func() bool {
		events := handler.GetEvents()
		return len(events) > 0 && events[0].EventType() == domain.EventTypeAssetCreated
	}, time.Second, 10*time.Millisecond)

	// 자산 업데이트
	newAmount, err := valueobjects.NewMoney(2000.0, "USD")
	require.NoError(t, err)
	asset.Update("Updated Asset", domain.Stock, newAmount)
	err = repo.Update(context.Background(), asset)
	require.NoError(t, err)

	// 업데이트 이벤트 확인
	assert.Eventually(t, func() bool {
		events := handler.GetEvents()
		return len(events) >= 3 &&
			events[1].EventType() == domain.EventTypeAssetUpdated &&
			events[2].EventType() == domain.EventTypeAssetAmountChanged
	}, time.Second, 10*time.Millisecond)

	// 자산 삭제
	err = repo.Delete(context.Background(), asset.ID)
	require.NoError(t, err)

	// 삭제 이벤트 확인
	assert.Eventually(t, func() bool {
		events := handler.GetEvents()
		return len(events) >= 4 && events[3].EventType() == domain.EventTypeAssetDeleted
	}, time.Second, 10*time.Millisecond)
}

func TestAssetService_ConcurrentEventHandling(t *testing.T) {
	// Given
	eventBus := newTestEventBus()
	handler := &testEventHandler{}

	// 모든 이벤트 타입에 대해 핸들러 등록
	err := eventBus.Subscribe(domain.EventTypeAssetCreated, handler)
	require.NoError(t, err)

	repo := infrastructure.NewMemoryAssetRepository(eventBus)

	// When
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			amount, err := valueobjects.NewMoney(float64(100*(idx+1)), "USD")
			require.NoError(t, err)

			asset := domain.NewAsset(
				fmt.Sprintf("user-%d", idx),
				domain.Stock,
				fmt.Sprintf("Test Asset %d", idx),
				amount,
			)

			err = repo.Save(context.Background(), asset)
			require.NoError(t, err)
		}(i)
	}
	wg.Wait()

	// Then
	assert.Eventually(t, func() bool {
		events := handler.GetEvents()
		return len(events) == 10
	}, 2*time.Second, 10*time.Millisecond)

	events := handler.GetEvents()
	assert.Len(t, events, 10)
	for _, event := range events {
		assert.Equal(t, domain.EventTypeAssetCreated, event.EventType())
	}
}

func TestAssetService_EventReplay(t *testing.T) {
	// Given
	eventBus := newTestEventBus()
	handler := &testEventHandler{}

	// 모든 이벤트 타입에 대해 핸들러 등록
	eventTypes := []string{
		domain.EventTypeAssetCreated,
		domain.EventTypeAssetUpdated,
		domain.EventTypeAssetAmountChanged,
		domain.EventTypeAssetDeleted,
	}
	for _, eventType := range eventTypes {
		err := eventBus.Subscribe(eventType, handler)
		require.NoError(t, err)
	}

	repo := infrastructure.NewMemoryAssetRepository(eventBus)
	amount, err := valueobjects.NewMoney(1000.0, "USD")
	require.NoError(t, err)

	// When
	asset := domain.NewAsset("user-1", domain.Stock, "Test Asset", amount)
	err = repo.Save(context.Background(), asset)
	require.NoError(t, err)

	// 이벤트 발생 확인
	assert.Eventually(t, func() bool {
		events := handler.GetEvents()
		return len(events) > 0 && events[0].EventType() == domain.EventTypeAssetCreated
	}, time.Second, 10*time.Millisecond)

	// 이벤트 리플레이
	events := handler.GetEvents()
	require.NotEmpty(t, events)

	// 새로운 핸들러로 이벤트 리플레이
	replayHandler := &testEventHandler{}
	for _, event := range events {
		err := replayHandler.HandleEvent(context.Background(), event)
		require.NoError(t, err)
	}

	// Then
	replayedEvents := replayHandler.GetEvents()
	assert.Equal(t, len(events), len(replayedEvents))
	for i, event := range events {
		assert.Equal(t, event.EventType(), replayedEvents[i].EventType())
		assert.Equal(t, event.AggregateID(), replayedEvents[i].AggregateID())
	}
}
