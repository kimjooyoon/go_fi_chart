package domain

// BaseMetric 기본 메트릭 구현체입니다.
type BaseMetric struct {
	name        string
	metricType  Type
	value       Value
	description string
}

// NewBaseMetric 새로운 BaseMetric을 생성합니다.
func NewBaseMetric(name string, metricType Type, value Value, description string) *BaseMetric {
	return &BaseMetric{
		name:        name,
		metricType:  metricType,
		value:       value,
		description: description,
	}
}

// Name 메트릭의 이름을 반환합니다.
func (m *BaseMetric) Name() string {
	return m.name
}

// Type 메트릭의 타입을 반환합니다.
func (m *BaseMetric) Type() Type {
	return m.metricType
}

// Value 메트릭의 값을 반환합니다.
func (m *BaseMetric) Value() Value {
	return m.value
}

// Description 메트릭의 설명을 반환합니다.
func (m *BaseMetric) Description() string {
	return m.description
}
