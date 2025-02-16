package github

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	"github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
)

// Collector GitHub 메트릭 수집기입니다.
type Collector struct {
	metrics   []metrics.Metric
	publisher domain.Publisher
	mu        sync.RWMutex
}

// NewCollector 새로운 GitHub 메트릭 수집기를 생성합니다.
func NewCollector(publisher domain.Publisher) *Collector {
	return &Collector{
		publisher: publisher,
	}
}

// Add 메트릭을 추가합니다.
func (c *Collector) Add(metric metrics.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
}

// Collect 메트릭을 수집하고 이벤트를 발행합니다.
func (c *Collector) Collect(ctx context.Context) ([]metrics.Metric, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]metrics.Metric, len(c.metrics))
	copy(result, c.metrics)

	for _, metric := range c.metrics {
		event := domain.NewMonitoringEvent(domain.TypeMetricCollected, metric)
		if err := c.publisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}

	c.metrics = make([]metrics.Metric, 0)
	return result, nil
}

// Metrics 수집된 메트릭 목록을 반환합니다.
func (c *Collector) Metrics() []metrics.Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]metrics.Metric, len(c.metrics))
	copy(result, c.metrics)
	return result
}
