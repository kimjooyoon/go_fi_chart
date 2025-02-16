package simple

import (
	"github.com/aske/go_fi_chart/services/monitoring/internal/metrics"
)

// Metric 기본 메트릭 구현체입니다.
type Metric struct {
	name        string
	metricType  metrics.Type
	value       metrics.Value
	description string
}

// NewMetric 새로운 메트릭을 생성합니다.
func NewMetric(name string, metricType metrics.Type, value float64, description string) *Metric {
	return &Metric{
		name:        name,
		metricType:  metricType,
		value:       metrics.NewValue(value, nil),
		description: description,
	}
}

// Name 메트릭의 이름을 반환합니다.
func (m *Metric) Name() string {
	return m.name
}

// Type 메트릭의 타입을 반환합니다.
func (m *Metric) Type() metrics.Type {
	return m.metricType
}

// Value 메트릭의 값을 반환합니다.
func (m *Metric) Value() metrics.Value {
	return m.value
}

// Description 메트릭의 설명을 반환합니다.
func (m *Metric) Description() string {
	return m.description
}

// WithLabels 레이블이 추가된 새로운 메트릭을 반환합니다.
func (m *Metric) WithLabels(labels map[string]string) *Metric {
	return &Metric{
		name:        m.name,
		metricType:  m.metricType,
		value:       m.value.WithLabels(labels),
		description: m.description,
	}
}
