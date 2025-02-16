package github

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/internal/domain"
)

// ActionCollector GitHub 액션 메트릭을 수집하는 컬렉터입니다.
type ActionCollector struct {
	mu        sync.RWMutex
	metrics   []domain.Metric
	publisher domain.Publisher
}

// NewActionCollector 새로운 ActionCollector를 생성합니다.
func NewActionCollector(publisher domain.Publisher) *ActionCollector {
	return &ActionCollector{
		metrics:   make([]domain.Metric, 0),
		publisher: publisher,
	}
}

// AddMetric 메트릭을 추가합니다.
func (c *ActionCollector) AddMetric(metric domain.Metric) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics = append(c.metrics, metric)
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *ActionCollector) Collect(ctx context.Context) ([]domain.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	metrics := make([]domain.Metric, len(c.metrics))
	copy(metrics, c.metrics)

	event := domain.NewMonitoringEvent(domain.TypeMetricCollected, "github-collector", metrics, nil)
	if err := c.publisher.Publish(ctx, event); err != nil {
		return nil, err
	}

	c.metrics = make([]domain.Metric, 0)
	return metrics, nil
}
