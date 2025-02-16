package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// MetricType 메트릭의 타입을 나타냅니다.
type MetricType string

const (
	TypeCounter   MetricType = "counter"
	TypeGauge     MetricType = "gauge"
	TypeHistogram MetricType = "histogram"
	TypeSummary   MetricType = "summary"
)

// Metric 모니터링 시스템의 메트릭을 나타냅니다.
type Metric struct {
	Name        string            `json:"name"`
	Type        MetricType        `json:"type"`
	Value       float64           `json:"value"`
	Labels      map[string]string `json:"labels,omitempty"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description"`
}

// Collector 메트릭을 수집하는 인터페이스입니다.
type Collector interface {
	// Collect 메트릭을 수집합니다.
	Collect(ctx context.Context) ([]domain.Metric, error)
}

// SimpleCollector 기본적인 메트릭 수집기 구현체입니다.
type SimpleCollector struct {
	mu        sync.RWMutex
	metrics   []domain.Metric
	publisher domain.Publisher
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector(publisher domain.Publisher) *SimpleCollector {
	return &SimpleCollector{
		metrics:   make([]domain.Metric, 0),
		publisher: publisher,
	}
}

// Collect 수집된 메트릭을 반환합니다.
func (c *SimpleCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	metrics := make([]domain.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	// 메트릭 수집 이벤트 발행
	event := domain.NewMonitoringEvent(
		domain.TypeMetricCollected,
		"collector",
		domain.MetricPayload{Metrics: metrics},
		nil,
	)

	if err := c.publisher.Publish(ctx, event); err != nil {
		return nil, domain.NewError("metrics", domain.ErrCodeInternal, "메트릭 수집 이벤트 발행 실패")
	}

	return metrics, nil
}

// AddMetric 메트릭을 추가합니다.
func (c *SimpleCollector) AddMetric(metric domain.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if metric.Timestamp.IsZero() {
		metric.Timestamp = time.Now()
	}
	c.metrics = append(c.metrics, metric)
}

// Reset 수집된 메트릭을 초기화합니다.
func (c *SimpleCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]domain.Metric, 0)
}
