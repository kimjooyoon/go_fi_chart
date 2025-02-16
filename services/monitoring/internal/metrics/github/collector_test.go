package github

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	mu     sync.RWMutex
	events []domain.Event
}

func (p *mockPublisher) Publish(_ context.Context, evt domain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, evt)
	return nil
}

func (p *mockPublisher) Subscribe(_ domain.Handler) error {
	return nil
}

func (p *mockPublisher) Unsubscribe(_ domain.Handler) error {
	return nil
}

func Test_NewCollector_should_create_empty_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}

	// When
	collector := NewCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.Empty(t, collector.metrics)
}

func Test_Collector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)

	// When
	err := collector.AddActionStatusMetric("test-action", ActionStatusSuccess)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "test-action", metrics[0].Name())
	assert.Equal(t, float64(1), metrics[0].Value().Raw)
	assert.Equal(t, "success", metrics[0].Value().Labels["status"])
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.TypeMetricCollected, publisher.events[0].Type)
}

func Test_Collector_should_add_duration_metric(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)
	duration := 10 * time.Second

	// When
	err := collector.AddActionDurationMetric("test-action", duration)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)

	statusMetric := metrics[0]
	durationMetric := metrics[1]

	assert.Equal(t, "test-action", statusMetric.Name())
	assert.Equal(t, float64(1), statusMetric.Value().Raw)
	assert.Equal(t, "success", statusMetric.Value().Labels["status"])

	assert.Equal(t, "test-action_duration", durationMetric.Name())
	assert.InDelta(t, duration.Seconds(), durationMetric.Value().Raw, 0.1)
	assert.Equal(t, "test-action", durationMetric.Value().Labels["action"])
	assert.Len(t, publisher.events, 1)
}

func Test_Collector_should_handle_different_statuses(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)

	testCases := []struct {
		status   ActionStatus
		expected float64
	}{
		{ActionStatusSuccess, 1},
		{ActionStatusFailure, 0},
		{ActionStatusInProgress, 2},
		{"unknown", -1},
	}

	for _, tc := range testCases {
		// When
		err := collector.AddActionStatusMetric("test-action", tc.status)
		assert.NoError(t, err)
		metrics, err := collector.Collect(context.Background())

		// Then
		assert.NoError(t, err)
		assert.Len(t, metrics, 1)
		assert.Equal(t, "test-action", metrics[0].Name())
		assert.Equal(t, tc.expected, metrics[0].Value().Raw)
		assert.Equal(t, string(tc.status), metrics[0].Value().Labels["status"])
	}
}

func Test_Collector_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			_ = collector.AddActionStatusMetric("test-action", ActionStatusSuccess)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = collector.AddActionDurationMetric("test-action", time.Duration(i)*time.Second)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_, _ = collector.Collect(context.Background())
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
