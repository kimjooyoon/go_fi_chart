package github

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
	pkgdomain "github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	mu     sync.RWMutex
	events []pkgdomain.Event
}

func (p *mockPublisher) Publish(_ context.Context, evt pkgdomain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, evt)
	return nil
}

func (p *mockPublisher) Subscribe(_ pkgdomain.Handler) error {
	return nil
}

func (p *mockPublisher) Unsubscribe(_ pkgdomain.Handler) error {
	return nil
}

func Test_NewCollector_should_create_empty_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]pkgdomain.Event, 0)}

	// When
	collector := NewCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.Empty(t, collector.metrics)
}

func Test_Collector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]pkgdomain.Event, 0)}
	collector := NewCollector(publisher)
	metric := domain.NewBaseMetric(
		"test_metric",
		domain.TypeGauge,
		domain.NewValue(42.0, map[string]string{"test": "label"}),
		"Test metric",
	)

	// When
	collector.Add(metric)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metric.Name(), metrics[0].Name())
	assert.Equal(t, metric.Value().Raw, metrics[0].Value().Raw)
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, pkgdomain.TypeMetricCollected, publisher.events[0].Type)
}

func Test_Collector_should_return_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]pkgdomain.Event, 0)}
	collector := NewCollector(publisher)
	metric := domain.NewBaseMetric(
		"test_metric",
		domain.TypeGauge,
		domain.NewValue(42.0, map[string]string{"test": "label"}),
		"Test metric",
	)

	// When
	collector.Add(metric)
	metrics := collector.Metrics()

	// Then
	assert.Len(t, metrics, 1)
	assert.Equal(t, metric.Name(), metrics[0].Name())
	assert.Equal(t, metric.Value().Raw, metrics[0].Value().Raw)
}

func Test_Collector_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]pkgdomain.Event, 0)}
	collector := NewCollector(publisher)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			collector.Add(domain.NewBaseMetric(
				"test_metric",
				domain.TypeGauge,
				domain.NewValue(float64(i), map[string]string{"test": "label"}),
				"Test metric",
			))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_, _ = collector.Collect(context.Background())
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_ = collector.Metrics()
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
