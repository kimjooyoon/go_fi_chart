package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewBaseMetric_should_create_metric_with_valid_data(t *testing.T) {
	// Given
	name := "test_metric"
	metricType := TypeGauge
	value := NewValue(42.0, map[string]string{"test": "label"})
	description := "Test metric"

	// When
	metric := NewBaseMetric(name, metricType, value, description)

	// Then
	assert.Equal(t, name, metric.Name())
	assert.Equal(t, metricType, metric.Type())
	assert.Equal(t, value, metric.Value())
	assert.Equal(t, description, metric.Description())
}

func Test_NewValue_should_create_value_with_valid_data(t *testing.T) {
	// Given
	raw := 42.0
	labels := map[string]string{"test": "label"}

	// When
	value := NewValue(raw, labels)

	// Then
	assert.Equal(t, raw, value.Raw)
	assert.Equal(t, labels, value.Labels)
	assert.WithinDuration(t, time.Now(), value.Timestamp, time.Second)
}

func Test_Value_WithLabel_should_add_new_label(t *testing.T) {
	// Given
	value := NewValue(42.0, map[string]string{"test": "label"})

	// When
	newValue := value.WithLabel("new", "label")

	// Then
	assert.Equal(t, value.Raw, newValue.Raw)
	assert.Equal(t, value.Timestamp, newValue.Timestamp)
	assert.Len(t, newValue.Labels, 2)
	assert.Equal(t, "label", newValue.Labels["test"])
	assert.Equal(t, "label", newValue.Labels["new"])
}

func Test_Value_WithLabels_should_add_multiple_labels(t *testing.T) {
	// Given
	value := NewValue(42.0, map[string]string{"test": "label"})
	newLabels := map[string]string{
		"new1": "label1",
		"new2": "label2",
	}

	// When
	newValue := value.WithLabels(newLabels)

	// Then
	assert.Equal(t, value.Raw, newValue.Raw)
	assert.Equal(t, value.Timestamp, newValue.Timestamp)
	assert.Len(t, newValue.Labels, 3)
	assert.Equal(t, "label", newValue.Labels["test"])
	assert.Equal(t, "label1", newValue.Labels["new1"])
	assert.Equal(t, "label2", newValue.Labels["new2"])
}
