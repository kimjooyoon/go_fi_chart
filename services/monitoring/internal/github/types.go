package github

import (
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
)

// ActionStatus GitHub 액션의 상태입니다.
type ActionStatus string

const (
	ActionStatusSuccess    ActionStatus = "success"
	ActionStatusFailure    ActionStatus = "failure"
	ActionStatusInProgress ActionStatus = "in_progress"
)

// ActionMetric GitHub Actions 실행 메트릭을 나타냅니다.
type ActionMetric struct {
	WorkflowName string
	Status       ActionStatus
	Duration     time.Duration
	StartedAt    time.Time
	FinishedAt   time.Time
}

// NewActionMetric 새로운 액션 메트릭을 생성합니다.
func NewActionMetric(name string, status ActionStatus, duration time.Duration) ActionMetric {
	now := time.Now()
	return ActionMetric{
		WorkflowName: name,
		Status:       status,
		Duration:     duration,
		StartedAt:    now.Add(-duration),
		FinishedAt:   now,
	}
}

// ToMetric 액션 메트릭을 일반 메트릭으로 변환합니다.
func (m ActionMetric) ToMetric() domain.Metric {
	var value float64
	switch m.Status {
	case ActionStatusSuccess:
		value = 1
	case ActionStatusFailure:
		value = 0
	case ActionStatusInProgress:
		value = 2
	default:
		value = -1
	}

	return domain.NewBaseMetric(
		m.WorkflowName,
		domain.TypeGauge,
		domain.NewValue(value, map[string]string{
			"status": string(m.Status),
		}),
		"GitHub 액션 실행 상태",
	)
}

// ToDurationMetric 액션 실행 시간 메트릭을 생성합니다.
func (m ActionMetric) ToDurationMetric() domain.Metric {
	return domain.NewBaseMetric(
		m.WorkflowName+"_duration",
		domain.TypeGauge,
		domain.NewValue(m.Duration.Seconds(), map[string]string{
			"action": m.WorkflowName,
		}),
		"GitHub 액션 실행 시간 (초)",
	)
}
