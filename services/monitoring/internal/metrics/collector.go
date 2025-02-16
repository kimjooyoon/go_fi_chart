package metrics

import (
	"context"
	"sync"
	"time"
)

// MetricType 메트릭의 종류를 나타냅니다.
type MetricType string

const (
	TypeGauge   MetricType = "GAUGE"
	TypeCounter MetricType = "COUNTER"
)

// Metric 수집된 메트릭을 나타냅니다.
type Metric struct {
	Name      string            `json:"name"`
	Type      MetricType        `json:"type"`
	Value     float64           `json:"value"`
	Labels    map[string]string `json:"labels"`
	Timestamp time.Time         `json:"timestamp"`
}

// Collector 메트릭을 수집하는 인터페이스입니다.
type Collector interface {
	// Collect 메트릭을 수집합니다.
	Collect(ctx context.Context) ([]Metric, error)
}

// SimpleCollector 기본적인 메트릭 수집 구현체입니다.
type SimpleCollector struct {
	mu      sync.RWMutex
	metrics []Metric
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector() *SimpleCollector {
	return &SimpleCollector{
		metrics: make([]Metric, 0),
	}
}

// Collect 현재 저장된 메트릭을 반환합니다.
func (c *SimpleCollector) Collect(_ context.Context) ([]Metric, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]Metric, len(c.metrics))
	copy(result, c.metrics)
	return result, nil
}

// AddMetric 새로운 메트릭을 추가합니다.
func (c *SimpleCollector) AddMetric(metric Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metric.Timestamp = time.Now()
	c.metrics = append(c.metrics, metric)
}

// Reset 저장된 메트릭을 초기화합니다.
func (c *SimpleCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = make([]Metric, 0)
}
