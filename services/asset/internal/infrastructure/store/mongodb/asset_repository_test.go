package mongodb

import (
	"context"
	"testing"
	"time"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	
	"github.com/aske/go_fi_chart/internal/common/repository"
	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
)

// 테스트용 이벤트 버스 구현
type mockEventBus struct {
	events []events.Event
}

func (m *mockEventBus) Publish(_ context.Context, event events.Event) error {
	m.events = append(m.events, event)
	return nil
}

func (m *mockEventBus) Subscribe(_ context.Context, _ string, _ events.EventHandler) error {
	return nil
}

func (m *mockEventBus) Unsubscribe(_ context.Context, _ string, _ events.EventHandler) error {
	return nil
}

func setupMongoDBTest(t *testing.T) (*AssetRepository, *mongo.Client, *mockEventBus) {
	// MongoDB 테스트 연결 설정
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	require.NoError(t, err, "Failed to connect to MongoDB")
	
	// 테스트 DB 초기화
	db := client.Database("asset_test_db")
	err = db.Drop(ctx)
	require.NoError(t, err, "Failed to drop test database")
	
	// 모의 이벤트 버스 생성
	eventBus := &mockEventBus{events: make([]events.Event, 0)}
	
	// 테스트용 저장소 생성
	repo := NewAssetRepository(db, eventBus)
	err = repo.Init(ctx)
	require.NoError(t, err, "Failed to initialize repository")
	
	return repo, client, eventBus
}

func TestAssetRepository_SaveAndFindByID(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 테스트용 자산 생성
	money, err := valueobjects.NewMoney(100.0, "USD")
	require.NoError(t, err)
	
	asset := domain.NewAsset("user1", domain.Stock, "테스트 주식", money)
	
	// 저장
	err = repo.Save(ctx, asset)
	require.NoError(t, err, "Failed to save asset")
	
	// ID로 조회
	found, err := repo.FindByID(ctx, asset.ID)
	require.NoError(t, err, "Failed to find asset by ID")
	assert.Equal(t, asset.ID, found.ID, "Asset ID mismatch")
	assert.Equal(t, asset.UserID, found.UserID, "UserID mismatch")
	assert.Equal(t, asset.Type, found.Type, "Asset type mismatch")
	assert.Equal(t, asset.Name, found.Name, "Asset name mismatch")
	assert.Equal(t, asset.Amount.Amount, found.Amount.Amount, "Asset amount mismatch")
	assert.Equal(t, asset.Amount.Currency, found.Amount.Currency, "Asset currency mismatch")
}

func TestAssetRepository_FindAll(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 여러 자산 저장
	money1, _ := valueobjects.NewMoney(100.0, "USD")
	money2, _ := valueobjects.NewMoney(200.0, "EUR")
	money3, _ := valueobjects.NewMoney(300.0, "USD")
	
	asset1 := domain.NewAsset("user1", domain.Stock, "주식 1", money1)
	asset2 := domain.NewAsset("user2", domain.Bond, "채권 1", money2)
	asset3 := domain.NewAsset("user1", domain.Cash, "현금 1", money3)
	
	err := repo.Save(ctx, asset1)
	require.NoError(t, err)
	err = repo.Save(ctx, asset2)
	require.NoError(t, err)
	err = repo.Save(ctx, asset3)
	require.NoError(t, err)
	
	// 모든 자산 조회
	assets, err := repo.FindAll(ctx)
	require.NoError(t, err)
	assert.Len(t, assets, 3, "Should find 3 assets")
	
	// 페이지네이션 적용
	assets, err = repo.FindAll(ctx, repository.WithLimit(2), repository.WithOffset(0))
	require.NoError(t, err)
	assert.Len(t, assets, 2, "Should find 2 assets with pagination")
	
	assets, err = repo.FindAll(ctx, repository.WithLimit(2), repository.WithOffset(2))
	require.NoError(t, err)
	assert.Len(t, assets, 1, "Should find 1 asset with pagination offset")
}

func TestAssetRepository_Update(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 자산 생성 및 저장
	money, _ := valueobjects.NewMoney(100.0, "USD")
	asset := domain.NewAsset("user1", domain.Stock, "테스트 주식", money)
	
	err := repo.Save(ctx, asset)
	require.NoError(t, err)
	
	// 자산 업데이트
	newMoney, _ := valueobjects.NewMoney(150.0, "USD")
	err = asset.Update("업데이트된 주식", domain.Stock, newMoney)
	require.NoError(t, err)
	
	err = repo.Update(ctx, asset)
	require.NoError(t, err)
	
	// 업데이트된 자산 조회
	updated, err := repo.FindByID(ctx, asset.ID)
	require.NoError(t, err)
	assert.Equal(t, "업데이트된 주식", updated.Name, "Name should be updated")
	assert.Equal(t, 150.0, updated.Amount.Amount, "Amount should be updated")
}

func TestAssetRepository_Delete(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 자산 생성 및 저장
	money, _ := valueobjects.NewMoney(100.0, "USD")
	asset := domain.NewAsset("user1", domain.Stock, "테스트 주식", money)
	
	err := repo.Save(ctx, asset)
	require.NoError(t, err)
	
	// 자산 삭제
	err = repo.Delete(ctx, asset.ID)
	require.NoError(t, err)
	
	// 삭제된 자산 조회 시도
	_, err = repo.FindByID(ctx, asset.ID)
	assert.Error(t, err, "Should not find deleted asset")
	assert.True(t, errors.Is(err, repository.ErrEntityNotFound), "Error should be ErrEntityNotFound")
}

func TestAssetRepository_FindByUserID(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 여러 사용자의 자산 저장
	money1, _ := valueobjects.NewMoney(100.0, "USD")
	money2, _ := valueobjects.NewMoney(200.0, "EUR")
	money3, _ := valueobjects.NewMoney(300.0, "USD")
	
	asset1 := domain.NewAsset("user1", domain.Stock, "주식 1", money1)
	asset2 := domain.NewAsset("user2", domain.Bond, "채권 1", money2)
	asset3 := domain.NewAsset("user1", domain.Cash, "현금 1", money3)
	
	err := repo.Save(ctx, asset1)
	require.NoError(t, err)
	err = repo.Save(ctx, asset2)
	require.NoError(t, err)
	err = repo.Save(ctx, asset3)
	require.NoError(t, err)
	
	// user1의 자산만 조회
	assets, err := repo.FindByUserID(ctx, "user1")
	require.NoError(t, err)
	assert.Len(t, assets, 2, "Should find 2 assets for user1")
	
	// user2의 자산만 조회
	assets, err = repo.FindByUserID(ctx, "user2")
	require.NoError(t, err)
	assert.Len(t, assets, 1, "Should find 1 asset for user2")
}

func TestAssetRepository_FindByType(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 여러 유형의 자산 저장
	money1, _ := valueobjects.NewMoney(100.0, "USD")
	money2, _ := valueobjects.NewMoney(200.0, "EUR")
	money3, _ := valueobjects.NewMoney(300.0, "USD")
	
	asset1 := domain.NewAsset("user1", domain.Stock, "주식 1", money1)
	asset2 := domain.NewAsset("user2", domain.Bond, "채권 1", money2)
	asset3 := domain.NewAsset("user1", domain.Stock, "주식 2", money3)
	
	err := repo.Save(ctx, asset1)
	require.NoError(t, err)
	err = repo.Save(ctx, asset2)
	require.NoError(t, err)
	err = repo.Save(ctx, asset3)
	require.NoError(t, err)
	
	// Stock 유형 자산만 조회
	assets, err := repo.FindByType(ctx, domain.Stock)
	require.NoError(t, err)
	assert.Len(t, assets, 2, "Should find 2 stock assets")
	
	// Bond 유형 자산만 조회
	assets, err = repo.FindByType(ctx, domain.Bond)
	require.NoError(t, err)
	assert.Len(t, assets, 1, "Should find 1 bond asset")
}

func TestAssetRepository_Count(t *testing.T) {
	repo, client, _ := setupMongoDBTest(t)
	defer client.Disconnect(context.Background())
	
	ctx := context.Background()
	
	// 여러 자산 저장
	money1, _ := valueobjects.NewMoney(100.0, "USD")
	money2, _ := valueobjects.NewMoney(200.0, "EUR")
	money3, _ := valueobjects.NewMoney(300.0, "USD")
	
	asset1 := domain.NewAsset("user1", domain.Stock, "주식 1", money1)
	asset2 := domain.NewAsset("user2", domain.Bond, "채권 1", money2)
	asset3 := domain.NewAsset("user1", domain.Cash, "현금 1", money3)
	
	err := repo.Save(ctx, asset1)
	require.NoError(t, err)
	err = repo.Save(ctx, asset2)
	require.NoError(t, err)
	err = repo.Save(ctx, asset3)
	require.NoError(t, err)
	
	// 전체 자산 개수 조회
	count, err := repo.Count(ctx)
	require.NoError(t, err)
	assert.Equal(t, int64(3), count, "Should count 3 assets")
	
	// user1의 자산 개수 조회
	count, err = repo.CountByUserID(ctx, "user1")
	require.NoError(t, err)
	assert.Equal(t, int64(2), count, "Should count 2 assets for user1")
	
	// Stock 유형 자산 개수 조회
	count, err = repo.CountByType(ctx, domain.Stock)
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Should count 1 stock asset")
}