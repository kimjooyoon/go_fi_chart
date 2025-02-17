package domain

import "context"

// MetricStorage는 메트릭을 저장하고 조회하는 인터페이스입니다.
type MetricStorage interface {
	// Save는 메트릭을 저장합니다.
	Save(ctx context.Context, metrics []Metric) error
	// Get은 주어진 조건에 맞는 메트릭을 조회합니다.
	Get(ctx context.Context, filter MetricFilter) ([]Metric, error)
}

// MetricFilter는 메트릭 조회 조건을 정의합니다.
type MetricFilter struct {
	MetricType MetricType
	StartTime  int64
	EndTime    int64
	Labels     map[string]string
}

// NewMetricFilter는 새로운 MetricFilter를 생성합니다.
func NewMetricFilter(metricType MetricType, startTime, endTime int64, labels map[string]string) MetricFilter {
	return MetricFilter{
		MetricType: metricType,
		StartTime:  startTime,
		EndTime:    endTime,
		Labels:     labels,
	}
}
