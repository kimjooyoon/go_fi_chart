package github

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	"github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
)

// Collector GitHub 메트릭을 수집하는 컬렉터입니다.
type Collector struct {
	mu        sync.Mutex
	metrics   []ActionMetric
	publisher domain.Publisher
}

// NewCollector 새로운 GitHub 메트릭 컬렉터를 생성합니다.
func NewCollector(publisher domain.Publisher) *Collector {
	return &Collector{
		metrics:   make([]ActionMetric, 0),
		publisher: publisher,
	}
}

// AddActionStatusMetric GitHub 액션 상태 메트릭을 추가합니다.
func (c *Collector) AddActionStatusMetric(name string, status ActionStatus) error {
	metric := NewActionMetric(name, status, 0)
	c.mu.Lock()
	c.metrics = append(c.metrics, metric)
	c.mu.Unlock()
	return nil
}

// AddActionDurationMetric GitHub 액션 실행 시간 메트릭을 추가합니다.
func (c *Collector) AddActionDurationMetric(name string, duration time.Duration) error {
	metric := NewActionMetric(name, ActionStatusSuccess, duration)
	c.mu.Lock()
	c.metrics = append(c.metrics, metric)
	c.mu.Unlock()
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *Collector) Collect(ctx context.Context) ([]metrics.Metric, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	result := make([]metrics.Metric, 0, len(c.metrics)*2)
	for _, m := range c.metrics {
		metric := m.ToMetric()
		result = append(result, metric)

		if m.Duration > 0 {
			result = append(result, m.ToDurationMetric())
		}

		event := domain.NewMonitoringEvent(domain.TypeMetricCollected, metric)
		if err := c.publisher.Publish(ctx, event); err != nil {
			return nil, err
		}
	}

	c.metrics = c.metrics[:0]
	return result, nil
}
