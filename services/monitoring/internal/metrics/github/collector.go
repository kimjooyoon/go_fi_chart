package github

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
)

// ActionStatus GitHub Actions의 실행 상태를 나타냅니다.
type ActionStatus string

const (
	StatusSuccess ActionStatus = "success"
	StatusFailure ActionStatus = "failure"
	StatusRunning ActionStatus = "running"
)

// ActionMetric GitHub Actions 실행 메트릭을 나타냅니다.
type ActionMetric struct {
	WorkflowName string
	Status       ActionStatus
	Duration     time.Duration
	StartedAt    time.Time
	FinishedAt   time.Time
}

// ActionCollector GitHub Actions 메트릭을 수집하는 컬렉터입니다.
type ActionCollector struct {
	baseCollector *metrics.SimpleCollector
}

// NewActionCollector 새로운 GitHub Actions 메트릭 컬렉터를 생성합니다.
func NewActionCollector() *ActionCollector {
	return &ActionCollector{
		baseCollector: metrics.NewSimpleCollector(),
	}
}

// Collect GitHub Actions 메트릭을 수집합니다.
func (c *ActionCollector) Collect(ctx context.Context) ([]metrics.Metric, error) {
	return c.baseCollector.Collect(ctx)
}

// AddActionMetric GitHub Actions 메트릭을 추가합니다.
func (c *ActionCollector) AddActionMetric(actionMetric ActionMetric) {
	c.baseCollector.AddMetric(metrics.Metric{
		Name:  "github_action_duration_seconds",
		Type:  metrics.TypeGauge,
		Value: actionMetric.Duration.Seconds(),
		Labels: map[string]string{
			"workflow": actionMetric.WorkflowName,
			"status":   string(actionMetric.Status),
		},
		Timestamp: actionMetric.FinishedAt,
	})

	c.baseCollector.AddMetric(metrics.Metric{
		Name: "github_action_status",
		Type: metrics.TypeGauge,
		Value: map[ActionStatus]float64{
			StatusSuccess: 1.0,
			StatusFailure: 0.0,
			StatusRunning: 0.5,
		}[actionMetric.Status],
		Labels: map[string]string{
			"workflow": actionMetric.WorkflowName,
		},
		Timestamp: actionMetric.FinishedAt,
	})
}
