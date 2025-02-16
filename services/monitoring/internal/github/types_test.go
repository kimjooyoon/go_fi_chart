package github

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewActionMetric_should_create_metric_with_valid_data(t *testing.T) {
	// Given
	name := "test-workflow"
	status := ActionStatusSuccess
	duration := 10 * time.Second

	// When
	metric := NewActionMetric(name, status, duration)

	// Then
	assert.Equal(t, name, metric.WorkflowName)
	assert.Equal(t, status, metric.Status)
	assert.Equal(t, duration, metric.Duration)
	assert.WithinDuration(t, time.Now(), metric.FinishedAt, time.Second)
	assert.WithinDuration(t, metric.FinishedAt.Add(-duration), metric.StartedAt, time.Second)
}

func Test_ActionMetric_ToMetric_should_convert_to_domain_metric(t *testing.T) {
	// Given
	metric := NewActionMetric("test-workflow", ActionStatusSuccess, 10*time.Second)

	// When
	domainMetric := metric.ToMetric()

	// Then
	assert.Equal(t, metric.WorkflowName, domainMetric.Name())
	assert.Equal(t, float64(1), domainMetric.Value().Raw)
	assert.Equal(t, string(metric.Status), domainMetric.Value().Labels["status"])
}

func Test_ActionMetric_ToDurationMetric_should_convert_to_duration_metric(t *testing.T) {
	// Given
	duration := 10 * time.Second
	metric := NewActionMetric("test-workflow", ActionStatusSuccess, duration)

	// When
	durationMetric := metric.ToDurationMetric()

	// Then
	assert.Equal(t, metric.WorkflowName+"_duration", durationMetric.Name())
	assert.Equal(t, duration.Seconds(), durationMetric.Value().Raw)
	assert.Equal(t, metric.WorkflowName, durationMetric.Value().Labels["action"])
}

func Test_ActionMetric_ToMetric_should_handle_different_statuses(t *testing.T) {
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
		// Given
		metric := NewActionMetric("test-workflow", tc.status, 0)

		// When
		domainMetric := metric.ToMetric()

		// Then
		assert.Equal(t, tc.expected, domainMetric.Value().Raw)
		assert.Equal(t, string(tc.status), domainMetric.Value().Labels["status"])
	}
}
