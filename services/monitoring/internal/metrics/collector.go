package metrics

import (
	"context"

	"github.com/aske/go_fi_chart/services/monitoring/metrics/collectors"
	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
	pkgdomain "github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
)

// DomainPublisherAdapter pkg/domain.Publisher를 metrics/domain.Publisher로 변환합니다.
type DomainPublisherAdapter struct {
	publisher pkgdomain.Publisher
}

// NewDomainPublisherAdapter 새로운 DomainPublisherAdapter를 생성합니다.
func NewDomainPublisherAdapter(publisher pkgdomain.Publisher) *DomainPublisherAdapter {
	return &DomainPublisherAdapter{
		publisher: publisher,
	}
}

// Publish 메트릭 이벤트를 발행합니다.
func (a *DomainPublisherAdapter) Publish(ctx context.Context, metrics []domain.Metric) error {
	event := pkgdomain.NewMonitoringEvent(pkgdomain.TypeMetricCollected, metrics)
	return a.publisher.Publish(ctx, event)
}

// SimpleCollector 기본적인 메트릭 수집기 구현체입니다.
type SimpleCollector struct {
	*collectors.BaseCollector
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector(publisher pkgdomain.Publisher) *SimpleCollector {
	adapter := NewDomainPublisherAdapter(publisher)
	return &SimpleCollector{
		BaseCollector: collectors.NewBaseCollector(adapter),
	}
}
