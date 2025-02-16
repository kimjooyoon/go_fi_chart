package infrastructure

import (
	"context"
	"testing"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventBus는 테스트용 이벤트 버스입니다.
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(ctx context.Context, event events.Event) error {
	args := m.Called(ctx, event)
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
	// 준비
	eventBus := new(MockEventBus)
	repo := NewMemoryAssetRepository(eventBus)
	ctx := context.Background()

	amount1, _ := valueobjects.NewMoney(1000.0, "USD")
	asset1 := domain.NewAsset("user-1", domain.Stock, "자산 1", amount1)

	amount2, _ := valueobjects.NewMoney(2000.0, "USD")
	asset2 := domain.NewAsset("user-2", domain.Stock, "자산 2", amount2)

	amount3, _ := valueobjects.NewMoney(3000.0, "USD")
	asset3 := domain.NewAsset("user-1", domain.Bond, "자산 3", amount3)

	eventBus.On("Publish", ctx, mock.AnythingOfType("*events.BaseEvent")).Return(nil)
	repo.Save(ctx, asset1)
	repo.Save(ctx, asset2)
	repo.Save(ctx, asset3)

	// 실행
	assets, err := repo.FindByType(ctx, domain.Stock)

	// 검증
	assert.NoError(t, err)
	assert.Len(t, assets, 2)
	assert.Contains(t, assets, asset1)
	assert.Contains(t, assets, asset2)
}
