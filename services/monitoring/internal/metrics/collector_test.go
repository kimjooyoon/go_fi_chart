package metrics

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleCollector_should_create_empty_collector(t *testing.T) {
	// When
	collector := NewSimpleCollector()

	// Then
	metrics, err := collector.Collect(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, metrics)
}

func Test_SimpleCollector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	collector := NewSimpleCollector()
	metric := Metric{
		Name:   "test_metric",
		Type:   TypeGauge,
		Value:  42.0,
		Labels: map[string]string{"test": "value"},
	}

	// When
	collector.AddMetric(metric)

	// Then
	metrics, err := collector.Collect(context.Background())
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, metric.Name, metrics[0].Name)
	assert.Equal(t, metric.Type, metrics[0].Type)
	assert.Equal(t, metric.Value, metrics[0].Value)
	assert.Equal(t, metric.Labels, metrics[0].Labels)
	assert.NotZero(t, metrics[0].Timestamp)
}

func Test_SimpleCollector_should_reset_metrics(t *testing.T) {
	// Given
	collector := NewSimpleCollector()
	metric := Metric{
		Name:  "test_metric",
		Type:  TypeGauge,
		Value: 42.0,
	}
	collector.AddMetric(metric)

	// When
	collector.Reset()

	// Then
	metrics, err := collector.Collect(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, metrics)
}

func Test_SimpleCollector_should_be_thread_safe(_ *testing.T) {
	// Given
	collector := NewSimpleCollector()
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			collector.AddMetric(Metric{
				Name:  "test_metric",
				Type:  TypeGauge,
				Value: float64(i),
			})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			_, _ = collector.Collect(context.Background())
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/10; i++ {
			collector.Reset()
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
