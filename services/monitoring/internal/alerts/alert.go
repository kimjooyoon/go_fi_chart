package alerts

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
)

// AlertLevel 알림의 심각도를 나타냅니다.
type AlertLevel = domain.AlertLevel

const (
	LevelInfo     = domain.LevelInfo
	LevelWarning  = domain.LevelWarning
	LevelError    = domain.LevelError
	LevelCritical = domain.LevelCritical
)

// Alert 모니터링 시스템의 알림을 나타냅니다.
type Alert = domain.Alert

// Notifier 알림을 처리하는 인터페이스입니다.
type Notifier interface {
	Notify(ctx context.Context, alert Alert) error
}

// SimpleNotifier 기본적인 알림 처리자 구현체입니다.
type SimpleNotifier struct {
	mu        sync.RWMutex
	publisher domain.Publisher
	handlers  []Notifier
}

// NewSimpleNotifier 새로운 SimpleNotifier를 생성합니다.
func NewSimpleNotifier(publisher domain.Publisher) *SimpleNotifier {
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
	evt := domain.NewMonitoringEvent(domain.TypeAlertTriggered, alert)
	if err := n.publisher.Publish(ctx, evt); err != nil {
		return err
	}

	return nil
}
