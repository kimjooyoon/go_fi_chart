package github

import (
	"context"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
	"github.com/stretchr/testify/assert"
)

func Test_NewActionCollector_should_create_empty_collector(t *testing.T) {
	// When
	collector := NewActionCollector()

	// Then
	metrics, err := collector.Collect(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, metrics)
}

func Test_ActionCollector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	collector := NewActionCollector()
	now := time.Now()
	actionMetric := ActionMetric{
		WorkflowName: "ci",
		Status:       StatusSuccess,
		Duration:     5 * time.Second,
		StartedAt:    now.Add(-5 * time.Second),
		FinishedAt:   now,
	}

	// When
	collector.AddActionMetric(actionMetric)

	// Then
	metrics, err := collector.Collect(context.Background())
	assert.NoError(t, err)
	assert.Len(t, metrics, 2)

	// Duration 메트릭 검증
	durationMetric := findMetricByName(metrics, "github_action_duration_seconds")
	assert.NotNil(t, durationMetric)
	assert.Equal(t, 5.0, durationMetric.Value)
	assert.Equal(t, "ci", durationMetric.Labels["workflow"])
	assert.Equal(t, "success", durationMetric.Labels["status"])
	assertTimeNear(t, now, durationMetric.Timestamp)

	// Status 메트릭 검증
	statusMetric := findMetricByName(metrics, "github_action_status")
	assert.NotNil(t, statusMetric)
	assert.Equal(t, 1.0, statusMetric.Value)
	assert.Equal(t, "ci", statusMetric.Labels["workflow"])
	assertTimeNear(t, now, statusMetric.Timestamp)
}

func Test_ActionCollector_should_handle_different_statuses(t *testing.T) {
	// Given
	collector := NewActionCollector()
	now := time.Now()
	testCases := []struct {
		status ActionStatus
		value  float64
	}{
		{StatusSuccess, 1.0},
		{StatusFailure, 0.0},
		{StatusRunning, 0.5},
	}

	for _, tc := range testCases {
		// When
		collector.AddActionMetric(ActionMetric{
			WorkflowName: "ci",
			Status:       tc.status,
			Duration:     time.Second,
			StartedAt:    now.Add(-time.Second),
			FinishedAt:   now,
		})

		// Then
		metrics, err := collector.Collect(context.Background())
		assert.NoError(t, err)

		statusMetric := findMetricByName(metrics, "github_action_status")
		assert.NotNil(t, statusMetric)
		assert.Equal(t, tc.value, statusMetric.Value)
		assert.Equal(t, "ci", statusMetric.Labels["workflow"])

		collector.baseCollector.Reset()
	}
}

func findMetricByName(metrics []metrics.Metric, name string) *metrics.Metric {
	for _, m := range metrics {
		if m.Name == name {
			return &m
		}
	}
	return nil
}

// assertTimeNear 두 시간이 1초 이내의 차이를 가지는지 확인합니다.
func assertTimeNear(t *testing.T, expected, actual time.Time) {
	diff := expected.Sub(actual)
	if diff < 0 {
		diff = -diff
	}
	assert.True(t, diff < time.Second,
		"시간 차이가 너무 큽니다. expected: %v, actual: %v, diff: %v",
		expected, actual, diff)
}
