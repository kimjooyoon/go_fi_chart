package domain

import (
	"context"
	"time"
)

// MetricRepository는 메트릭 저장소 인터페이스입니다.
type MetricRepository interface {
	// Save는 메트릭을 저장합니다.
	Save(ctx context.Context, metric *Metric) error
	// FindByID는 ID로 메트릭을 조회합니다.
	FindByID(ctx context.Context, id string) (*Metric, error)
	// FindByType은 타입으로 메트릭을 조회합니다.
	FindByType(ctx context.Context, metricType MetricType) ([]*Metric, error)
	// FindByTimeRange는 시간 범위로 메트릭을 조회합니다.
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]*Metric, error)
}
