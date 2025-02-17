package metrics

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/services/monitoring/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPublisher struct {
	mu     sync.Mutex
	events []domain.Event
}

func (p *mockPublisher) Publish(evt domain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = append(p.events, evt)
	return nil
}

func (p *mockPublisher) Subscribe(_ domain.Handler) error {
	return nil
}

func (p *mockPublisher) Unsubscribe(_ domain.Handler) error {
	return nil
}

func Test_NewSimpleCollector_should_create_empty_collector(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}

	// When
	collector := NewSimpleCollector(publisher)

	// Then
	assert.NotNil(t, collector)
	assert.NotNil(t, collector.BaseCollector)
}

func Test_SimpleCollector_should_add_and_collect_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	metric := domain.NewMetric(
		"test_metric",
		domain.MetricTypeGauge,
		domain.NewMetricValue(42.0, map[string]string{"test": "label"}),
		time.Now(),
	)
	adapter := NewMetricAdapter(metric)

	// When
	err := collector.AddMetric(adapter)
	assert.NoError(t, err)
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Len(t, metrics, 1)
	assert.Equal(t, adapter.Name(), metrics[0].Name())
	assert.Equal(t, adapter.Value().Raw, metrics[0].Value().Raw)
	assert.Len(t, publisher.events, 1)
	assert.Equal(t, domain.EventTypeMetricCollected, publisher.events[0].Type)
}

func Test_SimpleCollector_should_reset_metrics(t *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	metric := domain.NewMetric(
		"test_metric",
		domain.MetricTypeGauge,
		domain.NewMetricValue(42.0, map[string]string{"test": "label"}),
		time.Now(),
	)
	adapter := NewMetricAdapter(metric)

	// When
	collector.AddMetric(adapter)
	collector.Reset()
	metrics, err := collector.Collect(context.Background())

	// Then
	assert.NoError(t, err)
	assert.Empty(t, metrics)
}

func Test_SimpleCollector_should_be_thread_safe(_ *testing.T) {
	// Given
	publisher := &mockPublisher{events: make([]domain.Event, 0)}
	collector := NewSimpleCollector(publisher)
	iterations := 1000
	done := make(chan bool)

	// When
	go func() {
		for i := 0; i < iterations; i++ {
			metric := domain.NewMetric(
				"test_metric",
				domain.MetricTypeGauge,
				domain.NewMetricValue(float64(i), map[string]string{"test": "label"}),
				time.Now(),
			)
			collector.AddMetric(NewMetricAdapter(metric))
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			_, _ = collector.Collect(context.Background())
		}
		done <- true
	}()

	go func() {
		for i := 0; i < iterations/2; i++ {
			collector.Reset()
		}
		done <- true
	}()

	// Then
	for i := 0; i < 3; i++ {
		<-done
	}
}

// MockMetricRepository는 테스트용 메트릭 저장소입니다.
type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) Save(ctx context.Context, metric *domain.Metric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func (m *MockMetricRepository) FindByID(ctx context.Context, id string) (*domain.Metric, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Metric), args.Error(1)
}

func (m *MockMetricRepository) FindByType(ctx context.Context, metricType domain.MetricType) ([]*domain.Metric, error) {
	args := m.Called(ctx, metricType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Metric), args.Error(1)
}

func (m *MockMetricRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.Metric, error) {
	args := m.Called(ctx, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Metric), args.Error(1)
}

// MockMetricCollector는 테스트용 메트릭 수집기입니다.
type MockMetricCollector struct {
	mock.Mock
}

func (m *MockMetricCollector) Collect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockMetricCollector) Start(ctx context.Context) {
	m.Called(ctx)
}

func (m *MockMetricCollector) Stop() {
	m.Called()
}

func TestNewCollector(t *testing.T) {
	repo := new(MockMetricRepository)
	collectors := []MetricCollector{new(MockMetricCollector)}
	manager := NewCollector(repo, collectors)

	assert.NotNil(t, manager)
	assert.Equal(t, repo, manager.repository)
	assert.Equal(t, collectors, manager.collectors)
}

func TestCollector_Start(t *testing.T) {
	repo := new(MockMetricRepository)
	mockCollector := new(MockMetricCollector)
	manager := NewCollector(repo, []MetricCollector{mockCollector})

	mockCollector.On("Start", mock.Anything).Return()
	mockCollector.On("Stop").Return()
	mockCollector.On("Collect", mock.Anything).Return(nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := manager.Start(ctx)
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)
	manager.Stop()

	mockCollector.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCollector_Stop(t *testing.T) {
	repo := new(MockMetricRepository)
	mockCollector := new(MockMetricCollector)
	manager := NewCollector(repo, []MetricCollector{mockCollector})

	mockCollector.On("Start", mock.Anything).Return()
	mockCollector.On("Stop").Return()
	mockCollector.On("Collect", mock.Anything).Return(nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := manager.Start(ctx)
		assert.NoError(t, err)
	}()

	time.Sleep(50 * time.Millisecond)
	manager.Stop()
	wg.Wait()

	mockCollector.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCollector_Collect(t *testing.T) {
	repo := new(MockMetricRepository)
	mockCollector := new(MockMetricCollector)
	manager := NewCollector(repo, []MetricCollector{mockCollector})

	mockCollector.On("Collect", mock.Anything).Return(nil)
	repo.On("Save", mock.Anything, mock.Anything).Return(nil)

	ctx := context.Background()
	err := manager.Collect(ctx)
	assert.NoError(t, err)

	mockCollector.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCollector_CollectError(t *testing.T) {
	repo := new(MockMetricRepository)
	mockCollector := new(MockMetricCollector)
	collector := NewCollector(repo, []MetricCollector{mockCollector})

	// 수집 에러 설정
	expectedErr := domain.ErrMetricCollectionFailed
	mockCollector.On("Collect", mock.Anything).Return(expectedErr)

	// 수집 실행
	err := collector.Collect(context.Background())
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// 모든 기대값이 충족되었는지 확인
	mockCollector.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestCollector_SaveError(t *testing.T) {
	repo := new(MockMetricRepository)
	mockCollector := new(MockMetricCollector)
	collector := NewCollector(repo, []MetricCollector{mockCollector})

	// 메트릭 수집 성공, 저장 실패 설정
	mockCollector.On("Collect", mock.Anything).Return(nil)
	expectedErr := domain.ErrMetricSaveFailed
	repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Metric")).Return(expectedErr)

	// 수집 실행
	err := collector.Collect(context.Background())
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	// 모든 기대값이 충족되었는지 확인
	mockCollector.AssertExpectations(t)
	repo.AssertExpectations(t)
}
