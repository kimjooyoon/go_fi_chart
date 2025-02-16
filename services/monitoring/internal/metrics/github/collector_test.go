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

func Test_NewActionCollector_should_create_empty_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}

	// When
	collector := NewActionCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.Empty(t, collector.metrics)
}

func Test_ActionCollector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewActionCollector(publisher)
	startTime := time.Now()
	metric := ActionMetric{
		WorkflowName: "test-action",
		Status:       ActionStatusSuccess,
		StartedAt:    startTime,
	}

	// When
	err := collector.AddMetric(metric)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metric.WorkflowName, metrics[0].WorkflowName)
	assert.Equal(t, metric.Status, metrics[0].Status)
	assert.Equal(t, metric.StartedAt, metrics[0].StartedAt)
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.TypeMetricCollected, publisher.events[0].Type)
}

func Test_ActionCollector_should_add_duration_metric(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewActionCollector(publisher)
	duration := 10 * time.Second

	// When
	err := collector.AddActionDurationMetric("test-action", duration)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metric := metrics[0]
	assert.Equal(t, "test-action", metric.WorkflowName)
	assert.InDelta(t, duration, metric.Duration, float64(10*time.Millisecond))
	assert.True(t, metric.FinishedAt.After(metric.StartedAt))
	assert.InDelta(t, duration.Seconds(), metric.FinishedAt.Sub(metric.StartedAt).Seconds(), 0.1)
}

func Test_ActionCollector_should_handle_different_statuses(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewActionCollector(publisher)
	startTime := time.Now()

	testCases := []struct {
		status ActionStatus
	}{
		{ActionStatusSuccess},
		{ActionStatusFailure},
		{ActionStatusInProgress},
		{"unknown"},
	}

	for _, tc := range testCases {
		// When
		err := collector.AddActionStatusMetric("test-action", tc.status)
		assert.NoError(t, err)
	}

	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, len(testCases))

	for i, metric := range metrics {
		assert.Equal(t, "test-action", metric.WorkflowName)
		assert.Equal(t, testCases[i].status, metric.Status)
		assert.True(t, metric.StartedAt.After(startTime))
		assert.True(t, metric.StartedAt.Before(time.Now()))
	}
}

func Test_ActionCollector_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewActionCollector(publisher)
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
