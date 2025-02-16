package metrics

import (
	"context"
	"time"
)

// Type 메트릭의 타입을 나타냅니다.
type Type string

const (
	TypeCounter   Type = "counter"
	TypeGauge     Type = "gauge"
	TypeHistogram Type = "histogram"
	TypeSummary   Type = "summary"
)

// Value 메트릭의 값을 나타냅니다.
type Value struct {
	Raw       float64           `json:"raw"`
	Labels    map[string]string `json:"labels,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// NewValue 새로운 메트릭 값을 생성합니다.
func NewValue(raw float64, labels map[string]string) Value {
	return Value{
		Raw:       raw,
		Labels:    labels,
		Timestamp: time.Now(),
	}
}

// WithLabel 레이블을 추가한 새로운 Value를 반환합니다.
func (v Value) WithLabel(key, value string) Value {
	newLabels := make(map[string]string, len(v.Labels)+1)
	for k, v := range v.Labels {
		newLabels[k] = v
	}
	newLabels[key] = value
	return Value{
		Raw:       v.Raw,
		Labels:    newLabels,
		Timestamp: v.Timestamp,
	}
}

// WithLabels 여러 레이블을 추가한 새로운 Value를 반환합니다.
func (v Value) WithLabels(labels map[string]string) Value {
	newLabels := make(map[string]string, len(v.Labels)+len(labels))
	for k, v := range v.Labels {
		newLabels[k] = v
	}
	for k, v := range labels {
		newLabels[k] = v
	}
	return Value{
		Raw:       v.Raw,
		Labels:    newLabels,
		Timestamp: v.Timestamp,
	}
}

// Metric 메트릭 인터페이스입니다.
type Metric interface {
	// Name 메트릭의 이름을 반환합니다.
	Name() string
	// Type 메트릭의 타입을 반환합니다.
	Type() Type
	// Value 메트릭의 값을 반환합니다.
	Value() Value
	// Description 메트릭의 설명을 반환합니다.
	Description() string
}

// Collector 메트릭 수집기 인터페이스입니다.
type Collector interface {
	// Collect 메트릭을 수집합니다.
	Collect(ctx context.Context) ([]Metric, error)
}
