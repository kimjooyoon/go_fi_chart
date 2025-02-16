package events

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
	"github.com/stretchr/testify/assert"
)

type mockHandler struct {
	mu     sync.Mutex
	events []domain.Event
}

func (h *mockHandler) Handle(_ context.Context, event domain.Event) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.events = append(h.events, event)
	return nil
}

func Test_NewSimplePublisher_should_create_empty_publisher(t *testing.T) {
	// When
	publisher := NewSimplePublisher()

	// Then
	assert.NotNil(t, publisher)
	assert.Empty(t, publisher.handlers)
}

func Test_SimplePublisher_should_publish_event_to_subscribers(t *testing.T) {
	// Given
	publisher := NewSimplePublisher()
	handler := &mockHandler{events: make([]domain.Event, 0)}
	_ = publisher.Subscribe(handler)

	event := NewMonitoringEvent(
		TypeMetricCollected,
		"test",
		MetricPayload{
			Metrics: []metrics.Metric{
				{
					Name:  "test_metric",
					Type:  metrics.TypeGauge,
					Value: 42.0,
				},
			},
		},
		nil,
	)

	// When
	err := publisher.Publish(context.Background(), event)

	// Then
	assert.NoError(t, err)
	assert.Len(t, handler.events, 1)
	assert.Equal(t, TypeMetricCollected, handler.events[0].EventType())
	assert.Equal(t, "test", handler.events[0].Source())
}

func Test_SimplePublisher_should_unsubscribe_handler(t *testing.T) {
	// Given
	publisher := NewSimplePublisher()
	handler := &mockHandler{events: make([]domain.Event, 0)}
	_ = publisher.Subscribe(handler)

	event := NewMonitoringEvent(
		TypeMetricCollected,
		"test",
		nil,
		nil,
	)

	// When
	_ = publisher.Unsubscribe(handler)
	err := publisher.Publish(context.Background(), event)

	// Then
	assert.NoError(t, err)
	assert.Empty(t, handler.events)
}

func Test_SimplePublisher_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := NewSimplePublisher()
	handler := &mockHandler{events: make([]domain.Event, 0)}
	_ = publisher.Subscribe(handler)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			_ = publisher.Publish(context.Background(), NewMonitoringEvent(
				TypeMetricCollected,
				"test",
				nil,
				nil,
			))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = publisher.Subscribe(&mockHandler{events: make([]domain.Event, 0)})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = publisher.Unsubscribe(handler)
			_ = publisher.Subscribe(handler)
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
