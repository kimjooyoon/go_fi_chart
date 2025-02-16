package collectors

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
)

// BaseCollector 모든 컬렉터의 기본 구현을 제공합니다.
type BaseCollector struct {
	mu        sync.RWMutex
	metrics   []domain.Metric
	publisher domain.Publisher
}

// NewBaseCollector 새로운 BaseCollector를 생성합니다.
func NewBaseCollector(publisher domain.Publisher) *BaseCollector {
	return &BaseCollector{
		metrics:   make([]domain.Metric, 0),
		publisher: publisher,
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *BaseCollector) AddMetric(metric domain.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *BaseCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]domain.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	if err := c.publisher.Publish(ctx, metrics); err != nil {
		return nil, err
	}

	c.metrics = make([]domain.Metric, 0)
	return metrics, nil
}

// Reset 수집된 메트릭을 초기화합니다.
func (c *BaseCollector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]domain.Metric, 0)
}
