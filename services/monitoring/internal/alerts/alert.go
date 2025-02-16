package alerts

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/aske/go_fi_chart/services/monitoring/internal/events"
)

// AlertLevel 알림의 심각도를 나타냅니다.
type AlertLevel string

const (
	LevelInfo     AlertLevel = "INFO"
	LevelWarning  AlertLevel = "WARNING"
	LevelError    AlertLevel = "ERROR"
	LevelCritical AlertLevel = "CRITICAL"
)

// Alert 모니터링 시스템의 알림을 나타냅니다.
type Alert struct {
	ID        string            `json:"id"`
	Level     AlertLevel        `json:"level"`
	Source    string            `json:"source"`
	Message   string            `json:"message"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// Notifier 알림을 처리하는 인터페이스입니다.
type Notifier interface {
	Notify(ctx context.Context, alert Alert) error
}

// SimpleNotifier 기본적인 알림 처리자 구현체입니다.
type SimpleNotifier struct {
	mu        sync.RWMutex
	publisher events.Publisher
	handlers  []Notifier
}

// NewSimpleNotifier 새로운 SimpleNotifier를 생성합니다.
func NewSimpleNotifier(publisher events.Publisher) *SimpleNotifier {
	return &SimpleNotifier{
		publisher: publisher,
		handlers:  make([]Notifier, 0),
	}
}

// AddHandler 알림 핸들러를 추가합니다.
func (n *SimpleNotifier) AddHandler(handler Notifier) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.handlers = append(n.handlers, handler)
}

// RemoveHandler 알림 핸들러를 제거합니다.
func (n *SimpleNotifier) RemoveHandler(handler Notifier) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for i, h := range n.handlers {
		if h == handler {
			n.handlers = append(n.handlers[:i], n.handlers[i+1:]...)
			break
		}
	}
}

// Notify 알림을 처리하고 이벤트를 발행합니다.
func (n *SimpleNotifier) Notify(ctx context.Context, alert Alert) error {
	n.mu.RLock()
	defer n.mu.RUnlock()

	// 모든 핸들러에게 알림 전달
	for _, handler := range n.handlers {
		if err := handler.Notify(ctx, alert); err != nil {
			// 에러가 발생해도 다른 핸들러는 계속 실행
			continue
		}
	}

	// 알림 이벤트 발행
	event := events.NewMonitoringEvent(
		events.TypeAlertTriggered,
		alert.Source,
		alert,
		alert.Metadata,
	)

	if err := n.publisher.Publish(ctx, event); err != nil {
		return domain.NewError("alerts", domain.ErrCodeInternal, "알림 이벤트 발행 실패")
	}

	return nil
}
