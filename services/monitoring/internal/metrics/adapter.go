package metrics

import (
	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	metricsdomain "github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
)

// MetricAdapter는 internal/domain.Metric을 metrics/domain.Metric으로 변환합니다.
type MetricAdapter struct {
	metric     *domain.Metric
	metricID   string
	metricDesc string
}

// NewMetricAdapter는 새로운 MetricAdapter를 생성합니다.
func NewMetricAdapter(metric *domain.Metric) *MetricAdapter {
	return &MetricAdapter{
		metric:     metric,
		metricID:   metric.ID(),
		metricDesc: "Monitoring metric",
	}
}

// Name은 메트릭의 이름을 반환합니다.
func (a *MetricAdapter) Name() string {
	return a.metricID
}

// Type은 메트릭의 타입을 반환합니다.
func (a *MetricAdapter) Type() metricsdomain.Type {
	return metricsdomain.Type(a.metric.Type())
}

// Value는 메트릭의 값을 반환합니다.
func (a *MetricAdapter) Value() metricsdomain.Value {
	value := a.metric.Value()
	return metricsdomain.NewValue(value.Value(), value.Labels())
}

// Description은 메트릭의 설명을 반환합니다.
func (a *MetricAdapter) Description() string {
	return a.metricDesc
}
