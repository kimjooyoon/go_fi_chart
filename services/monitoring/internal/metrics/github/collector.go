package github

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	pkgmetrics "github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
)

// Collector GitHub 메트릭을 수집하는 컬렉터입니다.
type Collector struct {
	mu        sync.RWMutex
	metrics   []pkgmetrics.Metric
	publisher domain.Publisher
}

// NewCollector 새로운 GitHub 메트릭 컬렉터를 생성합니다.
func NewCollector(publisher domain.Publisher) *Collector {
	return &Collector{
		metrics:   make([]pkgmetrics.Metric, 0),
		publisher: publisher,
	}
}

// AddActionStatusMetric GitHub 액션 상태 메트릭을 추가합니다.
func (c *Collector) AddActionStatusMetric(name string, status ActionStatus) error {
	metric := NewActionMetric(name, status, 0)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = append(c.metrics, metric.ToMetric())
	return nil
}

// AddActionDurationMetric GitHub 액션 실행 시간 메트릭을 추가합니다.
func (c *Collector) AddActionDurationMetric(name string, duration time.Duration) error {
	metric := NewActionMetric(name, ActionStatusSuccess, duration)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = append(c.metrics, metric.ToMetric())
	if metric.Duration > 0 {
		c.metrics = append(c.metrics, metric.ToDurationMetric())
	}
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *Collector) Collect(ctx context.Context) ([]pkgmetrics.Metric, error) {
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
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]pkgmetrics.Metric, 0)
}
