package github

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	mu     sync.Mutex
	events []domain.Event
}

func (p *mockPublisher) Publish(evt domain.Event) error {
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
	err := collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusSuccess)
	assert.NoError(t, err)

	err = collector.Collect(context.Background())
	assert.NoError(t, err)

	// Then
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventTypeMetricCollected, publisher.events[0].Type)
}

func Test_Collector_should_add_duration_metric(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)
	duration := 5 * time.Second

	// When
	err := collector.AddActionDurationMetric("test-repo", "test-workflow", duration)
	assert.NoError(t, err)

	err = collector.Collect(context.Background())
	assert.NoError(t, err)

	// Then
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventTypeMetricCollected, publisher.events[0].Type)
}

func Test_Collector_should_reset_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)

	// When
	err := collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusSuccess)
	assert.NoError(t, err)

	collector.Reset()

	err = collector.Collect(context.Background())
	assert.NoError(t, err)

	// Then
	assert.Empty(t, publisher.events)
}

func Test_Collector_should_handle_multiple_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)

	// When
	err := collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusSuccess)
	assert.NoError(t, err)

	err = collector.AddActionDurationMetric("test-repo", "test-workflow", 5*time.Second)
	assert.NoError(t, err)

	err = collector.Collect(context.Background())
	assert.NoError(t, err)

	// Then
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventTypeMetricCollected, publisher.events[0].Type)
}

func Test_Collector_should_handle_different_statuses(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewCollector(publisher)

	// When
	err := collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusSuccess)
	assert.NoError(t, err)

	err = collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusFailure)
	assert.NoError(t, err)

	err = collector.Collect(context.Background())
	assert.NoError(t, err)

	// Then
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventTypeMetricCollected, publisher.events[0].Type)
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
			_ = collector.AddActionStatusMetric("test-repo", "test-workflow", ActionStatusSuccess)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = collector.AddActionDurationMetric("test-repo", "test-workflow", time.Duration(i)*time.Second)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = collector.Collect(context.Background())
		}
		done <- true
	}()

	// Then
	<-done
	<-done
	<-done
}
