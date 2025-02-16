package github

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/stretchr/testify/assert"
)

type mockPublisher struct {
	mu      sync.RWMutex
	metrics []domain.Metric
}

func (p *mockPublisher) Publish(_ context.Context, event domain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	metrics := event.Payload().([]domain.Metric)
	p.metrics = make([]domain.Metric, len(metrics))
	copy(p.metrics, metrics)
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
	startTime := time.Now()

	// when
	err := collector.AddActionMetric("test-action", ActionStatusSuccess)
	assert.NoError(t, err)

	metrics, err := collector.Collect(context.Background())

	// then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metric := metrics[0]
	assert.Equal(t, "github_action_status", metric.Name)
	assert.Equal(t, domain.MetricTypeGauge, metric.Type)
	assert.Equal(t, float64(0), metric.Value)
	assert.Equal(t, "test-action", metric.Labels["action"])
	assert.Equal(t, "GitHub 액션의 상태를 나타냅니다. 0: 성공, 1: 실패, 2: 진행 중", metric.Description)
	assert.True(t, metric.Timestamp.After(startTime))
	assert.True(t, metric.Timestamp.Before(time.Now()))
}

func Test_ActionCollector_should_add_duration_metric(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)
	duration := 10 * time.Second
	startTime := time.Now()

	// when
	err := collector.AddDurationMetric("test-action", duration)
	assert.NoError(t, err)

	metrics, err := collector.Collect(context.Background())

	// then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)

	metric := metrics[0]
	assert.Equal(t, "github_action_duration_seconds", metric.Name)
	assert.Equal(t, domain.MetricTypeGauge, metric.Type)
	assert.Equal(t, float64(10), metric.Value)
	assert.Equal(t, "test-action", metric.Labels["action"])
	assert.Equal(t, "GitHub 액션의 실행 시간(초)입니다.", metric.Description)
	assert.True(t, metric.Timestamp.After(startTime))
	assert.True(t, metric.Timestamp.Before(time.Now()))
}

func Test_ActionCollector_should_handle_different_statuses(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)
	startTime := time.Now()

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
		assert.Equal(t, domain.MetricTypeGauge, metric.Type)
		assert.Equal(t, "GitHub 액션의 상태를 나타냅니다. 0: 성공, 1: 실패, 2: 진행 중", metric.Description)
		assert.True(t, metric.Timestamp.After(startTime))
		assert.True(t, metric.Timestamp.Before(time.Now()))

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

func Test_ActionCollector_should_be_thread_safe(t *testing.T) {
	// given
	publisher := &mockPublisher{}
	collector := NewActionCollector(publisher)
	iterations := 1000
	done := make(chan bool)
	errChan := make(chan error)

	// when
	go func() {
		for i := 0; i < iterations; i++ {
			if err := collector.AddActionMetric("test-action", ActionStatusSuccess); err != nil {
				errChan <- err
				return
			}
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations; i++ {
			if err := collector.AddDurationMetric("test-action", time.Second); err != nil {
				errChan <- err
				return
			}
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			if _, err := collector.Collect(context.Background()); err != nil {
				errChan <- err
				return
			}
		}
		done <- true
	}()

	// then
	for i := 0; i < 3; i++ {
		select {
		case err := <-errChan:
			t.Errorf("unexpected error: %v", err)
			return
		case <-done:
			continue
		}
	}
}
