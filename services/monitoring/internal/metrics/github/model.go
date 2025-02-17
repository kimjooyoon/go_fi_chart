package github

import (
	"errors"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
)

// MetricType은 GitHub 메트릭의 타입을 나타냅니다.
type MetricType string

const (
	// MetricTypeStars 스타 수 메트릭
	MetricTypeStars MetricType = "STARS"
	// MetricTypeForks 포크 수 메트릭
	MetricTypeForks MetricType = "FORKS"
	// MetricTypeIssues 이슈 수 메트릭
	MetricTypeIssues MetricType = "ISSUES"
	// MetricTypePullRequests PR 수 메트릭
	MetricTypePullRequests MetricType = "PULL_REQUESTS"
	// MetricTypeContributors 기여자 수 메트릭
	MetricTypeContributors MetricType = "CONTRIBUTORS"
)

// Metric은 GitHub 메트릭을 나타냅니다.
type Metric struct {
	repository string
	metricType MetricType
	value      float64
	timestamp  time.Time
}

// NewMetric은 새로운 GitHub 메트릭을 생성합니다.
func NewMetric(repository string, metricType MetricType, value float64, timestamp time.Time) *Metric {
	return &Metric{
		repository: repository,
		metricType: metricType,
		value:      value,
		timestamp:  timestamp,
	}
}

// Repository는 메트릭이 속한 레포지토리를 반환합니다.
func (m *Metric) Repository() string {
	return m.repository
}

// MetricType은 메트릭의 타입을 반환합니다.
func (m *Metric) MetricType() MetricType {
	return m.metricType
}

// Value는 메트릭의 값을 반환합니다.
func (m *Metric) Value() float64 {
	return m.value
}

// Timestamp는 메트릭의 타임스탬프를 반환합니다.
func (m *Metric) Timestamp() time.Time {
	return m.timestamp
}

// Validate는 메트릭의 유효성을 검사합니다.
func (m *Metric) Validate() error {
	if m.repository == "" {
		return ErrInvalidRepository
	}
	if !isValidMetricType(m.metricType) {
		return ErrInvalidMetricType
	}
	return nil
}

// ToDomain은 GitHub 메트릭을 도메인 메트릭으로 변환합니다.
func (m *Metric) ToDomain() *domain.Metric {
	labels := map[string]string{
		"repository": m.repository,
		"type":       string(m.metricType),
	}
	return domain.NewMetric(
		m.repository+"."+string(m.metricType),
		domain.MetricTypeGitHub,
		domain.NewMetricValue(m.value, labels),
		m.timestamp,
	)
}

// isValidMetricType는 메트릭 타입의 유효성을 검사합니다.
func isValidMetricType(t MetricType) bool {
	switch t {
	case MetricTypeStars, MetricTypeForks, MetricTypeIssues, MetricTypePullRequests, MetricTypeContributors:
		return true
	default:
		return false
	}
}

// String은 메트릭 타입의 문자열 표현을 반환합니다.
func (t MetricType) String() string {
	return string(t)
}

// IsValid는 메트릭 타입의 유효성을 검사합니다.
func (t MetricType) IsValid() bool {
	return isValidMetricType(t)
}

var (
	// ErrInvalidRepository 잘못된 레포지토리 에러
	ErrInvalidRepository = errors.New("invalid repository")
	// ErrInvalidMetricType 잘못된 메트릭 타입 에러
	ErrInvalidMetricType = errors.New("invalid metric type")
)
