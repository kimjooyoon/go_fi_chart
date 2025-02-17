package github

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
)

// Collector GitHub 메트릭을 수집하는 컬렉터입니다.
type Collector struct {
	mu        sync.RWMutex
	metrics   []*domain.Metric
	publisher domain.Publisher
}

// NewCollector 새로운 GitHub 메트릭 컬렉터를 생성합니다.
func NewCollector(publisher domain.Publisher) *Collector {
	return &Collector{
		metrics:   make([]*domain.Metric, 0),
		publisher: publisher,
	}
}

// AddActionStatusMetric GitHub 액션 상태 메트릭을 추가합니다.
func (c *Collector) AddActionStatusMetric(repository, workflow string, status ActionStatus) error {
	metric := NewActionMetric(repository, workflow, status, 0, time.Now())
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = append(c.metrics, metric.ToDomain())
	return nil
}

// AddActionDurationMetric GitHub 액션 실행 시간 메트릭을 추가합니다.
func (c *Collector) AddActionDurationMetric(repository, workflow string, duration time.Duration) error {
	metric := NewActionMetric(repository, workflow, ActionStatusSuccess, duration, time.Now())
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = append(c.metrics, metric.ToDomain())
	return nil
}

// Collect 수집된 메트릭을 반환하고 이벤트를 발행합니다.
func (c *Collector) Collect(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.metrics) == 0 {
		return nil
	}

	evt := domain.NewMonitoringEvent(domain.EventTypeMetricCollected, c.metrics)
	if err := c.publisher.Publish(evt); err != nil {
		return err
	}

	c.metrics = make([]*domain.Metric, 0)
	return nil
}

// Reset 수집된 메트릭을 초기화합니다.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metrics = make([]*domain.Metric, 0)
}
