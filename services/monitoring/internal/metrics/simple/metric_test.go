package simple

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_NewMetric_should_create_metric_with_valid_data(t *testing.T) {
	// Given
	name := "test_metric"
	metricType := metrics.TypeGauge
	value := 42.0
	description := "Test metric"

	// When
	metric := NewMetric(name, metricType, value, description)
	metricWithLabels := metric.WithLabels(map[string]string{"test": "label"})

	// Then
	assert.Equal(t, name, metric.Name())
	assert.Equal(t, metricType, metric.Type())
	assert.Equal(t, value, metric.Value().Raw)
	assert.Empty(t, metric.Value().Labels)
	assert.Equal(t, description, metric.Description())
	assert.WithinDuration(t, time.Now(), metric.Value().Timestamp, time.Second)

	assert.Equal(t, name, metricWithLabels.Name())
	assert.Equal(t, metricType, metricWithLabels.Type())
	assert.Equal(t, value, metricWithLabels.Value().Raw)
	assert.Equal(t, map[string]string{"test": "label"}, metricWithLabels.Value().Labels)
	assert.Equal(t, description, metricWithLabels.Description())
}

func Test_NewMetric_should_handle_different_metric_types(t *testing.T) {
	testCases := []struct {
		metricType metrics.Type
	}{
		{metrics.TypeCounter},
		{metrics.TypeGauge},
		{metrics.TypeHistogram},
		{metrics.TypeSummary},
	}

	for _, tc := range testCases {
		// When
		metric := NewMetric("test_metric", tc.metricType, 42.0, "Test metric")

		// Then
		assert.Equal(t, tc.metricType, metric.Type())
	}
}
