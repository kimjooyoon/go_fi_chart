package domain

import (
	"context"
	"time"
)

// Type 메트릭의 타입을 나타냅니다.
type Type string

const (
	// TypeCounter 카운터 타입 메트릭입니다.
	TypeCounter Type = "counter"
	// TypeGauge 게이지 타입 메트릭입니다.
	TypeGauge Type = "gauge"
	// TypeHistogram 히스토그램 타입 메트릭입니다.
	TypeHistogram Type = "histogram"
	// TypeSummary 요약 타입 메트릭입니다.
	TypeSummary Type = "summary"
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
	// AddMetric 메트릭을 추가합니다.
	AddMetric(metric Metric) error
	// Reset 수집된 메트릭을 초기화합니다.
	Reset()
}

// Publisher 메트릭 이벤트 발행자 인터페이스입니다.
type Publisher interface {
	// Publish 메트릭 이벤트를 발행합니다.
	Publish(ctx context.Context, metrics []Metric) error
}

// Event 메트릭 이벤트입니다.
type Event struct {
	Type      string
	Metrics   []Metric
	Timestamp time.Time
}

const (
	// TypeMetricCollected 메트릭 수집 이벤트 타입입니다.
	TypeMetricCollected = "metric.collected"
)
