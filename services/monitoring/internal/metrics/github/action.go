package github

import (
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
)

// ActionStatus는 GitHub 액션의 상태를 나타냅니다.
type ActionStatus string

const (
	// ActionStatusSuccess는 성공 상태를 나타냅니다.
	ActionStatusSuccess ActionStatus = "success"
	// ActionStatusFailure는 실패 상태를 나타냅니다.
	ActionStatusFailure ActionStatus = "failure"
	// ActionStatusInProgress는 진행 중 상태를 나타냅니다.
	ActionStatusInProgress ActionStatus = "in_progress"
)

// ActionMetric은 GitHub 액션 메트릭을 나타냅니다.
type ActionMetric struct {
	repository string
	workflow   string
	status     ActionStatus
	duration   time.Duration
	timestamp  time.Time
}

// NewActionMetric은 새로운 GitHub 액션 메트릭을 생성합니다.
func NewActionMetric(repository, workflow string, status ActionStatus, duration time.Duration, timestamp time.Time) *ActionMetric {
	return &ActionMetric{
		repository: repository,
		workflow:   workflow,
		status:     status,
		duration:   duration,
		timestamp:  timestamp,
	}
}

// Repository는 메트릭이 속한 레포지토리를 반환합니다.
func (m *ActionMetric) Repository() string {
	return m.repository
}

// Workflow는 메트릭이 속한 워크플로우를 반환합니다.
func (m *ActionMetric) Workflow() string {
	return m.workflow
}

// Status는 액션의 상태를 반환합니다.
func (m *ActionMetric) Status() ActionStatus {
	return m.status
}

// Duration은 액션의 실행 시간을 반환합니다.
func (m *ActionMetric) Duration() time.Duration {
	return m.duration
}

// Timestamp는 메트릭의 타임스탬프를 반환합니다.
func (m *ActionMetric) Timestamp() time.Time {
	return m.timestamp
}

// ToDomain은 GitHub 액션 메트릭을 도메인 메트릭으로 변환합니다.
func (m *ActionMetric) ToDomain() *domain.Metric {
	var value float64
	switch m.status {
	case ActionStatusSuccess:
		value = 1
	case ActionStatusFailure:
		value = 0
	case ActionStatusInProgress:
		value = 2
	default:
		value = -1
	}

	labels := map[string]string{
		"repository":  m.repository,
		"workflow":    m.workflow,
		"status":      string(m.status),
		"metric_type": "action_status",
	}

	if m.duration > 0 {
		labels["metric_type"] = "action_duration"
		value = m.duration.Seconds()
	}

	return domain.NewMetric(
		m.repository+"_"+m.workflow,
		domain.MetricTypeGitHub,
		domain.NewMetricValue(value, labels),
		m.timestamp,
	)
}
