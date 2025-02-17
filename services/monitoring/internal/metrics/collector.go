package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	"github.com/aske/go_fi_chart/services/monitoring/metrics/collectors"
	metricsdomain "github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
)

// DomainPublisherAdapter는 도메인 Publisher를 메트릭 Publisher로 변환합니다.
type DomainPublisherAdapter struct {
	publisher domain.Publisher
}

// NewDomainPublisherAdapter는 새로운 DomainPublisherAdapter를 생성합니다.
func NewDomainPublisherAdapter(publisher domain.Publisher) *DomainPublisherAdapter {
	return &DomainPublisherAdapter{
		publisher: publisher,
	}
}

// Publish는 메트릭 이벤트를 발행합니다.
func (a *DomainPublisherAdapter) Publish(_ context.Context, metrics []metricsdomain.Metric) error {
	domainMetrics := make([]domain.Metric, 0, len(metrics))
	for _, m := range metrics {
		value := m.Value()
		domainMetrics = append(domainMetrics, *domain.NewMetric(
			m.Name(),
			domain.MetricType(m.Type()),
			domain.NewMetricValue(value.Raw, value.Labels),
			value.Timestamp,
		))
	}
	event := domain.NewMonitoringEvent(domain.EventTypeMetricCollected, domainMetrics)
	return a.publisher.Publish(event)
}

// SimpleCollector 기본적인 메트릭 수집기 구현체입니다.
type SimpleCollector struct {
	*collectors.BaseCollector
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector(publisher domain.Publisher) *SimpleCollector {
	adapter := NewDomainPublisherAdapter(publisher)
	return &SimpleCollector{
		BaseCollector: collectors.NewBaseCollector(adapter),
	}
}

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

// MetricCollectorManager는 여러 메트릭 수집기를 관리하는 구조체입니다.
type MetricCollectorManager struct {
	repository domain.MetricRepository
	collectors []MetricCollector
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewCollector는 새로운 MetricCollectorManager를 생성합니다.
func NewCollector(repository domain.MetricRepository, collectors []MetricCollector) *MetricCollectorManager {
	return &MetricCollectorManager{
		repository: repository,
		collectors: collectors,
		stopCh:     make(chan struct{}),
	}
}

// Start는 모든 메트릭 수집기를 시작합니다.
func (c *MetricCollectorManager) Start(ctx context.Context) error {
	for _, collector := range c.collectors {
		collector.Start(ctx)
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				return
			default:
				if err := c.Collect(ctx); err != nil {
					// TODO: 에러 로깅
					continue
				}
				time.Sleep(time.Second) // 수집 간격 조절
			}
		}
	}()
	return nil
}

// Stop은 모든 메트릭 수집기를 중지합니다.
func (c *MetricCollectorManager) Stop() {
	for _, collector := range c.collectors {
		collector.Stop()
	}
	close(c.stopCh)
	c.wg.Wait()
}

// Collect는 모든 메트릭 수집기에서 메트릭을 수집합니다.
func (c *MetricCollectorManager) Collect(ctx context.Context) error {
	for _, collector := range c.collectors {
		if err := collector.Collect(ctx); err != nil {
			return err
		}
	}

	// 수집된 메트릭을 저장소에 저장
	metric := domain.NewMetric(
		"test_metric",
		domain.MetricTypeGauge,
		domain.NewMetricValue(42.0, map[string]string{"test": "label"}),
		time.Now(),
	)
	if err := c.repository.Save(ctx, metric); err != nil {
		return err
	}

	return nil
}
