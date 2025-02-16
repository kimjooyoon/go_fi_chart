package github

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
)

// ActionStatus GitHub 액션의 상태를 나타냅니다.
type ActionStatus string

const (
	ActionStatusSuccess    ActionStatus = "success"
	ActionStatusFailure    ActionStatus = "failure"
	ActionStatusInProgress ActionStatus = "in_progress"
)

// MetricType GitHub 메트릭의 타입을 나타냅니다.
type MetricType string

const (
	TypeGauge MetricType = "gauge"
)

// Metric GitHub 메트릭을 나타냅니다.
type Metric struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description"`
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
	metrics   []ActionMetric
	publisher domain.Publisher
}

// NewActionCollector 새로운 ActionCollector를 생성합니다.
func NewActionCollector(publisher domain.Publisher) *ActionCollector {
	return &ActionCollector{
		metrics:   make([]ActionMetric, 0),
		publisher: publisher,
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *ActionCollector) AddMetric(metric ActionMetric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *ActionCollector) Collect(ctx context.Context) ([]ActionMetric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]ActionMetric, len(c.metrics))
	copy(metrics, c.metrics)

	evt := domain.NewMonitoringEvent(domain.TypeMetricCollected, metrics)
	if err := c.publisher.Publish(ctx, evt); err != nil {
		return nil, err
	}

	c.metrics = make([]ActionMetric, 0)
	return metrics, nil
}

// AddActionStatusMetric GitHub 액션 상태 메트릭을 추가합니다.
func (c *ActionCollector) AddActionStatusMetric(name string, status ActionStatus) error {
	return c.AddMetric(ActionMetric{
		WorkflowName: name,
		Status:       status,
		StartedAt:    time.Now(),
	})
}

// AddActionDurationMetric GitHub 액션 실행 시간 메트릭을 추가합니다.
func (c *ActionCollector) AddActionDurationMetric(name string, duration time.Duration) error {
	return c.AddMetric(ActionMetric{
		WorkflowName: name,
		Duration:     duration,
		StartedAt:    time.Now().Add(-duration),
		FinishedAt:   time.Now(),
	})
}
