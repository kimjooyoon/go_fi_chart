package github

import (
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
)

// ActionStatus GitHub 액션의 상태입니다.
type ActionStatus string

const (
	ActionStatusSuccess ActionStatus = "success"
	ActionStatusFailure ActionStatus = "failure"
	ActionStatusSkipped ActionStatus = "skipped"
)

// ActionStatusMetric GitHub 액션 상태 메트릭입니다.
type ActionStatusMetric struct {
	name   string
	status ActionStatus
}

// NewActionStatusMetric 새로운 액션 상태 메트릭을 생성합니다.
func NewActionStatusMetric(name string, status ActionStatus) *ActionStatusMetric {
	return &ActionStatusMetric{
		name:   name,
		status: status,
	}
}

// Name 메트릭 이름을 반환합니다.
func (m *ActionStatusMetric) Name() string {
	return m.name
}

// Type 메트릭 타입을 반환합니다.
func (m *ActionStatusMetric) Type() metrics.Type {
	return metrics.TypeGauge
}

// Value 메트릭 값을 반환합니다.
func (m *ActionStatusMetric) Value() metrics.Value {
	var value float64
	switch m.status {
	case ActionStatusSuccess:
		value = 1
	case ActionStatusFailure:
		value = 0
	case ActionStatusSkipped:
		value = -1
	default:
		value = -2
	}

	return metrics.NewValue(value, map[string]string{
		"action": m.name,
	})
}

// Description 메트릭 설명을 반환합니다.
func (m *ActionStatusMetric) Description() string {
	return "GitHub 액션 실행 상태"
}

// ActionDurationMetric GitHub 액션 실행 시간 메트릭입니다.
type ActionDurationMetric struct {
	name     string
	duration time.Duration
}

// NewActionDurationMetric 새로운 액션 실행 시간 메트릭을 생성합니다.
func NewActionDurationMetric(name string, duration time.Duration) *ActionDurationMetric {
	return &ActionDurationMetric{
		name:     name,
		duration: duration,
	}
}

// Name 메트릭 이름을 반환합니다.
func (m *ActionDurationMetric) Name() string {
	return m.name
}

// Type 메트릭 타입을 반환합니다.
func (m *ActionDurationMetric) Type() metrics.Type {
	return metrics.TypeGauge
}

// Value 메트릭 값을 반환합니다.
func (m *ActionDurationMetric) Value() metrics.Value {
	return metrics.NewValue(m.duration.Seconds(), map[string]string{
		"action": m.name,
	})
}

// Description 메트릭 설명을 반환합니다.
func (m *ActionDurationMetric) Description() string {
	return "GitHub 액션 실행 시간 (초)"
}
