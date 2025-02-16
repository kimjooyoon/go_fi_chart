package github

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// ActionStatus GitHub 액션의 상태를 나타냅니다.
type ActionStatus string

const (
	ActionStatusSuccess    ActionStatus = "success"
	ActionStatusFailure    ActionStatus = "failure"
	ActionStatusInProgress ActionStatus = "in_progress"
)

// statusToValue ActionStatus를 숫자로 변환합니다.
func (s ActionStatus) toValue() float64 {
	switch s {
	case ActionStatusSuccess:
		return 0
	case ActionStatusFailure:
		return 1
	case ActionStatusInProgress:
		return 2
	default:
		return -1
	}
}

// ActionMetric GitHub Actions 실행 메트릭을 나타냅니다.
type ActionMetric struct {
	WorkflowName string
	Status       ActionStatus
	Duration     time.Duration
	StartedAt    time.Time
	FinishedAt   time.Time
}

// ActionCollector GitHub 액션 메트릭을 수집하는 컬렉터입니다.
type ActionCollector struct {
	mu        sync.RWMutex
	metrics   []domain.Metric
	publisher domain.Publisher
}

// NewActionCollector 새로운 ActionCollector를 생성합니다.
func NewActionCollector(publisher domain.Publisher) *ActionCollector {
	return &ActionCollector{
		metrics:   make([]domain.Metric, 0),
		publisher: publisher,
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *ActionCollector) AddMetric(metric domain.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *ActionCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]domain.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	event := domain.NewMonitoringEvent(domain.TypeMetricCollected, "github-collector", metrics, nil)
	if err := c.publisher.Publish(ctx, event); err != nil {
		return nil, err
	}

	c.metrics = make([]domain.Metric, 0)
	return metrics, nil
}

// AddActionMetric GitHub 액션 메트릭을 추가합니다.
func (c *ActionCollector) AddActionMetric(name string, status ActionStatus) error {
	return c.AddMetric(domain.Metric{
		Name:        "github_action_status",
		Type:        domain.MetricTypeGauge,
		Value:       status.toValue(),
		Labels:      map[string]string{"action": name},
		Description: "GitHub 액션의 상태를 나타냅니다. 0: 성공, 1: 실패, 2: 진행 중",
		Timestamp:   time.Now(),
	})
}

// AddDurationMetric GitHub 액션 실행 시간 메트릭을 추가합니다.
func (c *ActionCollector) AddDurationMetric(name string, duration time.Duration) error {
	return c.AddMetric(domain.Metric{
		Name:        "github_action_duration_seconds",
		Type:        domain.MetricTypeGauge,
		Value:       duration.Seconds(),
		Labels:      map[string]string{"action": name},
		Description: "GitHub 액션의 실행 시간(초)입니다.",
		Timestamp:   time.Now(),
	})
}
