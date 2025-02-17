package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
)

// MetricCollector는 메트릭을 수집하고 저장하는 인터페이스입니다.
type MetricCollector interface {
	Collect(ctx context.Context) error
	Start(ctx context.Context)
	Stop()
}

// BaseMetricCollector는 기본 메트릭 수집기 구현체입니다.
type BaseMetricCollector struct {
	interval time.Duration
	storage  domain.MetricStorage
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewBaseMetricCollector는 새로운 BaseMetricCollector를 생성합니다.
func NewBaseMetricCollector(interval time.Duration, storage domain.MetricStorage) *BaseMetricCollector {
	return &BaseMetricCollector{
		interval: interval,
		storage:  storage,
		stopCh:   make(chan struct{}),
	}
}

// Start는 메트릭 수집을 시작합니다.
func (c *BaseMetricCollector) Start(ctx context.Context) {
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		ticker := time.NewTicker(c.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				return
			case <-ticker.C:
				if err := c.Collect(ctx); err != nil {
					// TODO: 에러 로깅
					continue
				}
			}
		}
	}()
}

// Stop은 메트릭 수집을 중지합니다.
func (c *BaseMetricCollector) Stop() {
	close(c.stopCh)
	c.wg.Wait()
}

// Collect는 메트릭을 수집합니다.
func (c *BaseMetricCollector) Collect(_ context.Context) error {
	// 구체적인 수집 로직은 하위 클래스에서 구현
	return nil
}
