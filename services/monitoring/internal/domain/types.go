package domain

import (
	"errors"
	"time"
)

// EventType은 모니터링 이벤트의 타입을 나타냅니다.
type EventType string

const (
	// EventTypeMetricCollected 메트릭이 수집되었을 때의 이벤트 타입
	EventTypeMetricCollected EventType = "METRIC_COLLECTED"
	TypeAlertTriggered       EventType = "alert_triggered"
)

// Event는 모니터링 이벤트를 나타냅니다.
type Event struct {
	Type      EventType
	Timestamp time.Time
	Data      interface{}
}

// NewMonitoringEvent는 새로운 모니터링 이벤트를 생성합니다.
func NewMonitoringEvent(eventType EventType, data interface{}) Event {
	return Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// Handler 이벤트 핸들러 인터페이스입니다.
type Handler interface {
	Handle(event Event) error
}

// Publisher는 이벤트를 발행하는 인터페이스입니다.
type Publisher interface {
	Publish(event Event) error
}

// MetricType은 메트릭의 타입을 나타냅니다.
type MetricType string

const (
	// MetricTypeAssetValue 자산 가치 메트릭
	MetricTypeAssetValue MetricType = "ASSET_VALUE"
	// MetricTypeTransactionCount 거래 수 메트릭
	MetricTypeTransactionCount MetricType = "TRANSACTION_COUNT"
	// MetricTypePortfolioValue 포트폴리오 가치 메트릭
	MetricTypePortfolioValue MetricType = "PORTFOLIO_VALUE"
	// MetricTypeUserCount 사용자 수 메트릭
	MetricTypeUserCount MetricType = "USER_COUNT"
	// MetricTypeGitHub GitHub 관련 메트릭
	MetricTypeGitHub MetricType = "GITHUB"
	// MetricTypeGauge 게이지 타입 메트릭
	MetricTypeGauge MetricType = "GAUGE"
)

// MetricValue는 메트릭의 값을 나타냅니다.
type MetricValue struct {
	value  float64
	labels map[string]string
}

// NewMetricValue는 새로운 메트릭 값을 생성합니다.
func NewMetricValue(value float64, labels map[string]string) *MetricValue {
	if labels == nil {
		labels = make(map[string]string)
	}
	return &MetricValue{
		value:  value,
		labels: labels,
	}
}

// Value는 메트릭의 값을 반환합니다.
func (v *MetricValue) Value() float64 {
	return v.value
}

// Labels는 메트릭의 레이블을 반환합니다.
func (v *MetricValue) Labels() map[string]string {
	return v.labels
}

// Equals는 다른 메트릭 값과 동일한지 비교합니다.
func (v *MetricValue) Equals(other *MetricValue) bool {
	if v.value != other.value {
		return false
	}
	if len(v.labels) != len(other.labels) {
		return false
	}
	for k, v := range v.labels {
		if otherV, exists := other.labels[k]; !exists || v != otherV {
			return false
		}
	}
	return true
}

// Add는 두 메트릭 값을 더합니다.
func (v *MetricValue) Add(other *MetricValue) *MetricValue {
	return NewMetricValue(v.value+other.value, v.labels)
}

// Metric은 메트릭을 나타냅니다.
type Metric struct {
	id         string
	metricType MetricType
	value      *MetricValue
	timestamp  time.Time
}

// NewMetric은 새로운 메트릭을 생성합니다.
func NewMetric(id string, metricType MetricType, value *MetricValue, timestamp time.Time) *Metric {
	return &Metric{
		id:         id,
		metricType: metricType,
		value:      value,
		timestamp:  timestamp,
	}
}

// ID는 메트릭의 ID를 반환합니다.
func (m *Metric) ID() string {
	return m.id
}

// Type은 메트릭의 타입을 반환합니다.
func (m *Metric) Type() MetricType {
	return m.metricType
}

// Value는 메트릭의 값을 반환합니다.
func (m *Metric) Value() *MetricValue {
	return m.value
}

// Timestamp는 메트릭의 타임스탬프를 반환합니다.
func (m *Metric) Timestamp() time.Time {
	return m.timestamp
}

// Validate는 메트릭의 유효성을 검사합니다.
func (m *Metric) Validate() error {
	if m.id == "" {
		return ErrInvalidMetricID
	}
	if !isValidMetricType(m.metricType) {
		return ErrInvalidMetricType
	}
	return nil
}

// isValidMetricType은 메트릭 타입의 유효성을 검사합니다.
func isValidMetricType(t MetricType) bool {
	switch t {
	case MetricTypeAssetValue,
		MetricTypeTransactionCount,
		MetricTypePortfolioValue,
		MetricTypeUserCount,
		MetricTypeGitHub,
		MetricTypeGauge:
		return true
	default:
		return false
	}
}

var (
	// ErrInvalidMetricID 잘못된 메트릭 ID 에러
	ErrInvalidMetricID = errors.New("invalid metric id")
	// ErrInvalidMetricType 잘못된 메트릭 타입 에러
	ErrInvalidMetricType = errors.New("invalid metric type")
	// ErrMetricCollectionFailed 메트릭 수집 실패 에러
	ErrMetricCollectionFailed = errors.New("metric collection failed")
	// ErrMetricSaveFailed 메트릭 저장 실패 에러
	ErrMetricSaveFailed = errors.New("metric save failed")
)
