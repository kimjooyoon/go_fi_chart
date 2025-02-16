package github

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
	pkgdomain "github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
)

// Collector GitHub 메트릭 수집기입니다.
type Collector struct {
	metrics   []domain.Metric
	publisher pkgdomain.Publisher
	mu        sync.RWMutex
}

// NewCollector 새로운 GitHub 메트릭 수집기를 생성합니다.
func NewCollector(publisher pkgdomain.Publisher) *Collector {
	return &Collector{
		publisher: publisher,
	}
}

// Add 메트릭을 추가합니다.
func (c *Collector) Add(metric domain.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
}

// Collect 메트릭을 수집하고 이벤트를 발행합니다.
func (c *Collector) Collect(ctx context.Context) ([]domain.Metric, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]domain.Metric, len(c.metrics))
	copy(result, c.metrics)

	for _, metric := range c.metrics {
		event := pkgdomain.NewMonitoringEvent(pkgdomain.TypeMetricCollected, metric)
		if err := c.publisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}

	c.metrics = make([]domain.Metric, 0)
	return result, nil
}

// Metrics 수집된 메트릭 목록을 반환합니다.
func (c *Collector) Metrics() []domain.Metric {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]domain.Metric, len(c.metrics))
	copy(result, c.metrics)
	return result
}
