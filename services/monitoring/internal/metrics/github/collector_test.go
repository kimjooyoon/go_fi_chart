package github

import (
	"context"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	metrics []domain.Metric
}

func (p *mockPublisher) Publish(_ context.Context, event domain.Event) error {
	p.metrics = event.Payload().([]domain.Metric)
	return nil
}

func (p *mockPublisher) Subscribe(_ domain.Handler) error {
	return nil
}

func (p *mockPublisher) Unsubscribe(_ domain.Handler) error {
	return nil
}

func Test_NewActionCollector_should_create_empty_collector(t *testing.T) {
	// given
	publisher := &mockPublisher{}

	// when
	collector := NewActionCollector(publisher)

	// then
	assert.NotNil(t, collector)
	assert.Empty(t, collector.metrics)
}

func Test_ActionCollector_should_add_and_collect_metrics(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)

	// when
	err := collector.AddActionMetric("test-action", ActionStatusSuccess)
	assert.NoError(t, err)

	metrics, err := collector.Collect(context.Background())

	// then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "github_action_status", metrics[0].Name)
	assert.Equal(t, float64(0), metrics[0].Value)
	assert.Equal(t, "test-action", metrics[0].Labels["action"])
}

func Test_ActionCollector_should_add_duration_metric(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)
	duration := 10 * time.Second

	// when
	err := collector.AddDurationMetric("test-action", duration)
	assert.NoError(t, err)

	metrics, err := collector.Collect(context.Background())

	// then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, "github_action_duration_seconds", metrics[0].Name)
	assert.Equal(t, float64(10), metrics[0].Value)
	assert.Equal(t, "test-action", metrics[0].Labels["action"])
}

func Test_ActionCollector_should_handle_different_statuses(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)

	// when
	err := collector.AddActionMetric("success-action", ActionStatusSuccess)
	assert.NoError(t, err)
	err = collector.AddActionMetric("failure-action", ActionStatusFailure)
	assert.NoError(t, err)
	err = collector.AddActionMetric("progress-action", ActionStatusInProgress)
	assert.NoError(t, err)

	metrics, err := collector.Collect(context.Background())

	// then
	assert.NoError(t, err)
	assert.Len(t, metrics, 3)

	for _, metric := range metrics {
		assert.Equal(t, "github_action_status", metric.Name)
		switch metric.Labels["action"] {
		case "success-action":
			assert.Equal(t, float64(0), metric.Value)
		case "failure-action":
			assert.Equal(t, float64(1), metric.Value)
		case "progress-action":
			assert.Equal(t, float64(2), metric.Value)
		default:
			t.Errorf("unexpected action: %s", metric.Labels["action"])
		}
	}
}
