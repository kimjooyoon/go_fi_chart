package alerts

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/pkg/domain"
	"github.com/stretchr/testify/assert"
)

type mockNotifier struct {
	mu     sync.Mutex
	alerts []Alert
}

func (n *mockNotifier) Notify(_ context.Context, alert Alert) error {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.alerts = append(n.alerts, alert)
	return nil
}

type mockPublisher struct {
	mu     sync.Mutex
	events []domain.Event
}

func (p *mockPublisher) Publish(_ context.Context, event domain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, event)
	return nil
}

func (p *mockPublisher) Subscribe(_ domain.Handler) error {
	return nil
}

func (p *mockPublisher) Unsubscribe(_ domain.Handler) error {
	return nil
}

func Test_NewSimpleNotifier_should_create_empty_notifier(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}

	// When
	notifier := NewSimpleNotifier(publisher)

	// Then
	assert.NotNil(t, notifier)
	assert.Empty(t, notifier.handlers)
}

func Test_SimpleNotifier_should_notify_handlers_and_publish_event(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	notifier := NewSimpleNotifier(publisher)
	handler := &mockNotifier{alerts: make([]Alert, 0)}
	notifier.AddHandler(handler)

	alert := Alert{
		ID:        "test-alert",
		Level:     LevelWarning,
		Source:    "test",
		Message:   "Test alert",
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"key": "value",
		},
	}

	// When
	err := notifier.Notify(context.Background(), alert)

	// Then
	assert.NoError(t, err)
	assert.Len(t, handler.alerts, 1)
	assert.Equal(t, alert.ID, handler.alerts[0].ID)
	assert.Equal(t, alert.Level, handler.alerts[0].Level)
	assert.Equal(t, alert.Message, handler.alerts[0].Message)

	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.TypeAlertTriggered, publisher.events[0].Type)
}

func Test_SimpleNotifier_should_remove_handler(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	notifier := NewSimpleNotifier(publisher)
	handler := &mockNotifier{alerts: make([]Alert, 0)}
	notifier.AddHandler(handler)

	alert := Alert{
		ID:        "test-alert",
		Level:     LevelWarning,
		Source:    "test",
		Message:   "Test alert",
		Timestamp: time.Now(),
	}

	// When
	notifier.RemoveHandler(handler)
	err := notifier.Notify(context.Background(), alert)

	// Then
	assert.NoError(t, err)
	assert.Empty(t, handler.alerts)
	assert.Len(t, publisher.events, 1)
}

func Test_SimpleNotifier_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	notifier := NewSimpleNotifier(publisher)
	handler := &mockNotifier{alerts: make([]Alert, 0)}
	notifier.AddHandler(handler)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			_ = notifier.Notify(context.Background(), Alert{
				ID:        "test-alert",
				Level:     LevelWarning,
				Source:    "test",
				Message:   "Test alert",
				Timestamp: time.Now(),
			})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			notifier.AddHandler(&mockNotifier{alerts: make([]Alert, 0)})
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			notifier.RemoveHandler(handler)
			notifier.AddHandler(handler)
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}
