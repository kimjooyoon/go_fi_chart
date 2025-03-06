package infrastructure

import (
	"context"
	"sync"
	"testing"

	"github.com/aske/go_fi_chart/internal/common/repository"
	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventBus는 테스트용 이벤트 버스입니다.
type MockEventBus struct {
	mock.Mock
	publishedEvents []events.Event
	mu              sync.Mutex
}

func (m *MockEventBus) Publish(ctx context.Context, event events.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	args := m.Called(ctx, event)
	if args.Error(0) == nil {
		m.publishedEvents = append(m.publishedEvents, event)
	}
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventType string, handler events.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handler events.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestNewMemoryAssetRepository(t *testing.T) {
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.assets)
	assert.Empty(t, repo.assets)
	assert.Equal(t, eventBus, repo.eventBus)
}

func TestMemoryAssetRepository_Save(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := domain.NewAsset("user-1", domain.Stock, "테스트 자산", amount)

	// 이벤트 발행 설정
	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)

	// 실행
	err := repo.Save(ctx, asset)

	// 검증
	assert.NoError(t, err)
	assert.Len(t, repo.assets, 1)
	assert.Equal(t, asset, repo.assets[asset.ID])
	eventBus.AssertExpectations(t)

	// 중복 저장 시도
	err = repo.Save(ctx, asset)
	assert.Error(t, err)
}

func TestMemoryAssetRepository_FindByID(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := domain.NewAsset("user-1", domain.Stock, "테스트 자산", amount)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset)

	// 실행 & 검증
	found, err := repo.FindByID(ctx, asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, asset, found)

	// 존재하지 않는 ID로 조회
	notFound, err := repo.FindByID(ctx, "non-existent")
	assert.Error(t, err)
	assert.Nil(t, notFound)
}

func TestMemoryAssetRepository_Update(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := domain.NewAsset("user-1", domain.Stock, "테스트 자산", amount)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset)

	// 업데이트
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	asset.Update("업데이트된 자산", domain.Stock, newAmount)

	// 실행
	err := repo.Update(ctx, asset)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, asset, repo.assets[asset.ID])
	eventBus.AssertExpectations(t)

	// 존재하지 않는 자산 업데이트
	nonExistentAsset := domain.NewAsset("user-2", domain.Stock, "존재하지 않는 자산", amount)
	err = repo.Update(ctx, nonExistentAsset)
	assert.Error(t, err)
}

func TestMemoryAssetRepository_Delete(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := domain.NewAsset("user-1", domain.Stock, "테스트 자산", amount)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset)

	// 실행
	err := repo.Delete(ctx, asset.ID)

	// 검증
	assert.NoError(t, err)
	assert.Empty(t, repo.assets)
	eventBus.AssertExpectations(t)

	// 존재하지 않는 자산 삭제
	err = repo.Delete(ctx, "non-existent")
	assert.Error(t, err)
}

func TestMemoryAssetRepository_FindByUserID(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount1, _ := valueobjects.NewMoney(1000.0, "USD")
	asset1 := domain.NewAsset("user-1", domain.Stock, "자산 1", amount1)

	amount2, _ := valueobjects.NewMoney(2000.0, "USD")
	asset2 := domain.NewAsset("user-1", domain.Bond, "자산 2", amount2)

	amount3, _ := valueobjects.NewMoney(3000.0, "USD")
	asset3 := domain.NewAsset("user-2", domain.Stock, "자산 3", amount3)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset1)
	repo.Save(ctx, asset2)
	repo.Save(ctx, asset3)

	// 실행
	assets, err := repo.FindByUserID(ctx, "user-1")

	// 검증
	assert.NoError(t, err)
	assert.Len(t, assets, 2)
	assert.Contains(t, assets, asset1)
	assert.Contains(t, assets, asset2)
}

func TestMemoryAssetRepository_FindByType(t *testing.T) {
	// Given
	ctx := context.Background()
	eventBus := &MockEventBus{}
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
	repo := NewMemoryAssetRepository(eventBus)

	// 여러 타입의 자산 생성
	stockAmount, _ := valueobjects.NewMoney(1000.0, "USD")
	bondAmount, _ := valueobjects.NewMoney(2000.0, "USD")

	stockAsset := domain.NewAsset("user-123", domain.Stock, "Stock Asset", stockAmount)
	bondAsset := domain.NewAsset("user-123", domain.Bond, "Bond Asset", bondAmount)

	err := repo.Save(ctx, stockAsset)
	assert.NoError(t, err)
	err = repo.Save(ctx, bondAsset)
	assert.NoError(t, err)

	// When
	stockAssets, err := repo.FindByType(ctx, domain.Stock)
	assert.NoError(t, err)
	bondAssets, err := repo.FindByType(ctx, domain.Bond)
	assert.NoError(t, err)

	// Then
	assert.Len(t, stockAssets, 1)
	assert.Len(t, bondAssets, 1)
	assert.Equal(t, stockAsset.ID, stockAssets[0].ID)
	assert.Equal(t, bondAsset.ID, bondAssets[0].ID)
}

func TestMemoryAssetRepository_ConcurrentOperations(t *testing.T) {
	// Given
	ctx := context.Background()
	eventBus := &MockEventBus{}
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
	repo := NewMemoryAssetRepository(eventBus)

	// 동시에 여러 자산을 저장
	var wg sync.WaitGroup
	assetCount := 100
	wg.Add(assetCount)

	// When
	for i := 0; i < assetCount; i++ {
		go func(idx int) {
			defer wg.Done()
			amount, _ := valueobjects.NewMoney(float64(1000+idx), "USD")
			asset := domain.NewAsset("user-123", domain.Stock, "Test Asset", amount)
			err := repo.Save(ctx, asset)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Then
	assets, err := repo.FindByUserID(ctx, "user-123")
	assert.NoError(t, err)
	assert.Len(t, assets, assetCount)

	// 이벤트 발행 검증
	eventBus.mu.Lock()
	assert.Len(t, eventBus.publishedEvents, assetCount)
	eventBus.mu.Unlock()
}

func TestMemoryAssetRepository_ConcurrentUpdates(t *testing.T) {
	// Given
	ctx := context.Background()
	eventBus := &MockEventBus{}
	eventBus.On("Publish", mock.Anything, mock.Anything).Return(nil)
	repo := NewMemoryAssetRepository(eventBus)

	// 초기 자산 생성
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := domain.NewAsset("user-123", domain.Stock, "Test Asset", amount)
	err := repo.Save(ctx, asset)
	assert.NoError(t, err)

	// 동시에 여러 번 업데이트
	var wg sync.WaitGroup
	updateCount := 100
	wg.Add(updateCount)

	// When
	for i := 0; i < updateCount; i++ {
		go func(idx int) {
			defer wg.Done()
			newAmount, _ := valueobjects.NewMoney(float64(1000+idx), "USD")
			asset.UpdateAmount(newAmount)
			err := repo.Update(ctx, asset)
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Then
	updatedAsset, err := repo.FindByID(ctx, asset.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedAsset)

	// 이벤트 발행 검증
	eventBus.mu.Lock()
	assert.True(t, len(eventBus.publishedEvents) > updateCount) // 초기 생성 이벤트 + 업데이트 이벤트
	eventBus.mu.Unlock()
}

func TestMemoryAssetRepository_DeleteWithConcurrentAccess(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryAssetRepository(eventBus)

	amount, _ := valueobjects.NewMoney(100.0, "USD")
	asset := domain.NewAsset(uuid.New().String(), domain.Stock, "Test Asset", amount)

	err := repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.Delete(context.Background(), asset.ID)
			if err != nil {
				assert.ErrorContains(t, err, domain.ErrAssetNotFound.Error())
			}
		}()
	}
	wg.Wait()

	// Then
	_, err = repo.FindByID(context.Background(), asset.ID)
	assert.ErrorContains(t, err, domain.ErrAssetNotFound.Error())
}

func TestMemoryAssetRepository_FindAll(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	// 세 개의 자산 준비
	amount1, _ := valueobjects.NewMoney(1000.0, "USD")
	asset1 := domain.NewAsset("user-1", domain.Stock, "자산 1", amount1)

	amount2, _ := valueobjects.NewMoney(2000.0, "USD")
	asset2 := domain.NewAsset("user-1", domain.Bond, "자산 2", amount2)

	amount3, _ := valueobjects.NewMoney(3000.0, "USD")
	asset3 := domain.NewAsset("user-2", domain.Stock, "자산 3", amount3)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset1)
	repo.Save(ctx, asset2)
	repo.Save(ctx, asset3)

	// 모든 자산 조회
	assets, err := repo.FindAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, assets, 3)
	assert.Contains(t, assets, asset1)
	assert.Contains(t, assets, asset2)
	assert.Contains(t, assets, asset3)

	// 페이지네이션 테스트
	limitedAssets, err := repo.FindAll(ctx, repository.WithLimit(2))
	assert.NoError(t, err)
	assert.Len(t, limitedAssets, 2)

	// 오프셋 테스트
	offsetAssets, err := repo.FindAll(ctx, repository.WithOffset(2), repository.WithLimit(2))
	assert.NoError(t, err)
	assert.Len(t, offsetAssets, 1)
}

func TestMemoryAssetRepository_Count(t *testing.T) {
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	// 자산 없는 상태에서 개수 확인
	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	// 세 개의 자산 준비
	amount1, _ := valueobjects.NewMoney(1000.0, "USD")
	asset1 := domain.NewAsset("user-1", domain.Stock, "자산 1", amount1)

	amount2, _ := valueobjects.NewMoney(2000.0, "USD")
	asset2 := domain.NewAsset("user-1", domain.Bond, "자산 2", amount2)

	amount3, _ := valueobjects.NewMoney(3000.0, "USD")
	asset3 := domain.NewAsset("user-2", domain.Stock, "자산 3", amount3)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset1)
	repo.Save(ctx, asset2)
	repo.Save(ctx, asset3)

	// 자산 개수 확인
	count, err = repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)

	// 자산 하나 삭제 후 개수 확인
	repo.Delete(ctx, asset1.ID)
	count, err = repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}
