package prometheus

import (
	"context"
	"fmt"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter Prometheus 익스포터 구현체입니다.
type Exporter struct {
	mu       sync.RWMutex
	registry *prometheus.Registry
	metrics  map[string]prometheus.Collector
}

// NewExporter 새로운 Exporter를 생성합니다.
func NewExporter() *Exporter {
	return &Exporter{
		registry: prometheus.NewRegistry(),
		metrics:  make(map[string]prometheus.Collector),
	}
}

// Export 메트릭을 Prometheus 형식으로 변환하여 등록합니다.
func (e *Exporter) Export(_ context.Context, metrics []metrics.Metric) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, m := range metrics {
		collector, err := e.getOrCreateCollector(m)
		if err != nil {
			return fmt.Errorf("메트릭 컬렉터 생성 실패: %w", err)
		}

		switch m.Type() {
		case "counter":
			if counter, ok := collector.(prometheus.Counter); ok {
				counter.Add(m.Value().Raw)
			}
		case "gauge":
			if gauge, ok := collector.(prometheus.Gauge); ok {
				gauge.Set(m.Value().Raw)
			}
		case "histogram":
			if histogram, ok := collector.(prometheus.Histogram); ok {
				histogram.Observe(m.Value().Raw)
			}
		case "summary":
			if summary, ok := collector.(prometheus.Summary); ok {
				summary.Observe(m.Value().Raw)
			}
		}
	}

	return nil
}

// GetRegistry Prometheus 레지스트리를 반환합니다.
func (e *Exporter) GetRegistry() *prometheus.Registry {
	return e.registry
}

// getOrCreateCollector 메트릭에 대한 Prometheus 컬렉터를 반환하거나 생성합니다.
func (e *Exporter) getOrCreateCollector(m metrics.Metric) (prometheus.Collector, error) {
	if collector, exists := e.metrics[m.Name()]; exists {
		return collector, nil
	}

	var collector prometheus.Collector

	switch m.Type() {
	case "counter":
		collector = prometheus.NewCounter(prometheus.CounterOpts{
			Name: m.Name(),
			Help: m.Description(),
		})
	case "gauge":
		collector = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: m.Name(),
			Help: m.Description(),
		})
	case "histogram":
		collector = prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: m.Name(),
			Help: m.Description(),
		})
	case "summary":
		collector = prometheus.NewSummary(prometheus.SummaryOpts{
			Name: m.Name(),
			Help: m.Description(),
		})
	default:
		return nil, fmt.Errorf("지원하지 않는 메트릭 타입: %s", m.Type())
	}

	if err := e.registry.Register(collector); err != nil {
		return nil, fmt.Errorf("메트릭 등록 실패: %w", err)
	}

	e.metrics[m.Name()] = collector
	return collector, nil
}
