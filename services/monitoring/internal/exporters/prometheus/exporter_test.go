package prometheus

import (
	"context"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func Test_NewExporter_should_create_empty_exporter(t *testing.T) {
	// When
	exporter := NewExporter()

	// Then
	assert.NotNil(t, exporter)
	assert.NotNil(t, exporter.registry)
	assert.Empty(t, exporter.metrics)
}

func Test_Exporter_should_export_metrics(t *testing.T) {
	// Given
	exporter := NewExporter()
	now := time.Now()

	testMetrics := []metrics.Metric{
		{
			Name:        "test_counter",
			Type:        metrics.TypeCounter,
			Value:       42.0,
			Labels:      map[string]string{"label": "value"},
			Timestamp:   now,
			Description: "Test counter metric",
		},
		{
			Name:        "test_gauge",
			Type:        metrics.TypeGauge,
			Value:       123.45,
			Labels:      map[string]string{"label": "value"},
			Timestamp:   now,
			Description: "Test gauge metric",
		},
	}

	// When
	err := exporter.Export(context.Background(), testMetrics)

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 2)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 2)

	for _, family := range metricFamilies {
		assert.Contains(t, []string{"test_counter", "test_gauge"}, family.GetName())
		assert.Len(t, family.GetMetric(), 1)
	}
}

func Test_Exporter_should_handle_different_metric_types(t *testing.T) {
	// Given
	exporter := NewExporter()
	now := time.Now()

	testMetrics := []metrics.Metric{
		{
			Name:        "test_counter",
			Type:        metrics.TypeCounter,
			Value:       1.0,
			Timestamp:   now,
			Description: "Test counter",
		},
		{
			Name:        "test_gauge",
			Type:        metrics.TypeGauge,
			Value:       2.0,
			Timestamp:   now,
			Description: "Test gauge",
		},
		{
			Name:        "test_histogram",
			Type:        metrics.TypeHistogram,
			Value:       3.0,
			Timestamp:   now,
			Description: "Test histogram",
		},
		{
			Name:        "test_summary",
			Type:        metrics.TypeSummary,
			Value:       4.0,
			Timestamp:   now,
			Description: "Test summary",
		},
	}

	// When
	err := exporter.Export(context.Background(), testMetrics)

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 4)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 4)
}

func Test_Exporter_should_update_existing_metrics(t *testing.T) {
	// Given
	exporter := NewExporter()
	now := time.Now()

	metric := metrics.Metric{
		Name:        "test_counter",
		Type:        metrics.TypeCounter,
		Value:       1.0,
		Timestamp:   now,
		Description: "Test counter",
	}

	// When
	err := exporter.Export(context.Background(), []metrics.Metric{metric})
	assert.NoError(t, err)

	metric.Value = 2.0
	err = exporter.Export(context.Background(), []metrics.Metric{metric})

	// Then
	assert.NoError(t, err)
	assert.Len(t, exporter.metrics, 1)

	metricFamilies, err := exporter.registry.Gather()
	assert.NoError(t, err)
	assert.Len(t, metricFamilies, 1)

	family := metricFamilies[0]
	assert.Equal(t, "test_counter", family.GetName())
	assert.Equal(t, 3.0, family.GetMetric()[0].GetCounter().GetValue())
}

func Test_Exporter_should_handle_invalid_metric_type(t *testing.T) {
	// Given
	exporter := NewExporter()
	now := time.Now()

	metric := metrics.Metric{
		Name:        "test_invalid",
		Type:        "invalid",
		Value:       1.0,
		Timestamp:   now,
		Description: "Test invalid metric",
	}

	// When
	err := exporter.Export(context.Background(), []metrics.Metric{metric})

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "지원하지 않는 메트릭 타입")
}

func Test_Exporter_should_handle_duplicate_registration(t *testing.T) {
	// Given
	exporter := NewExporter()
	now := time.Now()

	metric := metrics.Metric{
		Name:        "test_counter",
		Type:        metrics.TypeCounter,
		Value:       1.0,
		Timestamp:   now,
		Description: "Test counter",
	}

	// When
	err := exporter.Export(context.Background(), []metrics.Metric{metric})
	assert.NoError(t, err)

	// Try to register a different collector with the same name
	err = exporter.registry.Register(prometheus.NewCounter(prometheus.CounterOpts{
		Name: "test_counter",
		Help: "Duplicate counter",
	}))

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "previously registered descriptor")
}
