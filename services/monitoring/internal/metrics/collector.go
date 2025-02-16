package metrics

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	pkgmetrics "github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
)

// BaseCollector 모든 컬렉터의 기본 구현을 제공합니다.
type BaseCollector struct {
	mu        sync.RWMutex
	metrics   []pkgmetrics.Metric
	publisher domain.Publisher
}

// NewBaseCollector 새로운 BaseCollector를 생성합니다.
func NewBaseCollector(publisher domain.Publisher) *BaseCollector {
	return &BaseCollector{
		metrics:   make([]pkgmetrics.Metric, 0),
		publisher: publisher,
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *BaseCollector) AddMetric(metric pkgmetrics.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *BaseCollector) Collect(ctx context.Context) ([]pkgmetrics.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]pkgmetrics.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	evt := domain.NewMonitoringEvent(domain.TypeMetricCollected, metrics)
	if err := c.publisher.Publish(ctx, evt); err != nil {
		return nil, err
	}

	c.metrics = make([]pkgmetrics.Metric, 0)
	return metrics, nil
}

// Reset 수집된 메트릭을 초기화합니다.
func (c *BaseCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]pkgmetrics.Metric, 0)
}

// SimpleCollector 기본적인 메트릭 수집기 구현체입니다.
type SimpleCollector struct {
	*BaseCollector
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector(publisher domain.Publisher) *SimpleCollector {
	return &SimpleCollector{
		BaseCollector: NewBaseCollector(publisher),
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *SimpleCollector) AddMetric(metric pkgmetrics.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *SimpleCollector) Collect(ctx context.Context) ([]pkgmetrics.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]pkgmetrics.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	evt := domain.NewMonitoringEvent(domain.TypeMetricCollected, metrics)
	if err := c.publisher.Publish(ctx, evt); err != nil {
		return nil, err
	}

	c.metrics = make([]pkgmetrics.Metric, 0)
	return metrics, nil
}

// Reset 수집된 메트릭을 초기화합니다.
func (c *SimpleCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]pkgmetrics.Metric, 0)
}
