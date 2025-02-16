package collectors

import (
	"github.com/aske/go_fi_chart/services/monitoring/metrics/domain"
)

// SimpleCollector 기본적인 메트릭 수집기 구현체입니다.
type SimpleCollector struct {
	*BaseCollector
}

// NewSimpleCollector 새로운 SimpleCollector를 생성합니다.
func NewSimpleCollector(publisher domain.Publisher) *SimpleCollector {
	return &SimpleCollector{
		BaseCollector: NewBaseCollector(publisher),
	}
}
