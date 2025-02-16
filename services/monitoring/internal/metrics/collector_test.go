package metrics

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	pkgmetrics "github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	mu     sync.Mutex
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

func Test_NewSimpleCollector_should_create_empty_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}

	// When
	collector := NewSimpleCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.Empty(t, collector.metrics)
}

func Test_SimpleCollector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	metric := pkgmetrics.NewBaseMetric(
		"test_metric",
		pkgmetrics.TypeGauge,
		pkgmetrics.NewValue(42.0, map[string]string{"test": "label"}),
		"Test metric",
	)

	// When
	err := collector.AddMetric(metric)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metric.Name(), metrics[0].Name())
	assert.Equal(t, metric.Value().Raw, metrics[0].Value().Raw)
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.TypeMetricCollected, publisher.events[0].Type)
}

func Test_SimpleCollector_should_reset_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	metric := pkgmetrics.NewBaseMetric(
		"test_metric",
		pkgmetrics.TypeGauge,
		pkgmetrics.NewValue(42.0, map[string]string{"test": "label"}),
		"Test metric",
	)

	// When
	collector.AddMetric(metric)
	collector.Reset()
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Empty(t, metrics)
}

func Test_SimpleCollector_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			collector.AddMetric(pkgmetrics.NewBaseMetric(
				"test_metric",
				pkgmetrics.TypeGauge,
				pkgmetrics.NewValue(float64(i), map[string]string{"test": "label"}),
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
			collector.Reset()
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
