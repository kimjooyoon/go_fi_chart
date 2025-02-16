package collectors

import (
	"context"
	"testing"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
	"github.com/stretchr/testify/assert"
)

func Test_NewSimpleCollector_should_create_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{}

	// When
	collector := NewSimpleCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.NotNil(t, collector.BaseCollector)
}

func Test_SimpleCollector_should_use_base_collector_functionality(t *testing.T) {
	// Given
	publisher := &mockPublisher{}
	collector := NewSimpleCollector(publisher)
	metric := domain.NewBaseMetric(
		"test_metric",
		domain.TypeGauge,
		domain.NewValue(42.0, map[string]string{"test": "label"}),
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
	assert.Len(t, publisher.metrics, 1)
}
