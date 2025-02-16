package asset

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_memory_repo_should_save_and_find_asset_by_id(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset := NewAsset("user-1", Cash, "현금 자산", 1000000.0, "KRW")

	// When
	err := repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	found, err := repo.FindByID(context.Background(), asset.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, asset.ID, found.ID)
	assert.Equal(t, asset.Amount, found.Amount)
}

func Test_memory_repo_should_update_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset := NewAsset("user-1", Cash, "현금 자산", 1000000.0, "KRW")
	err := repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	asset.Amount = Money{Amount: 2000000.0, Currency: "KRW"}
	err = repo.Update(context.Background(), asset)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, asset.Amount, found.Amount)
}

func Test_memory_repo_should_delete_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset := NewAsset("user-1", Cash, "현금 자산", 1000000.0, "KRW")
	err := repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	err = repo.Delete(context.Background(), asset.ID)

	// Then
	assert.NoError(t, err)
	_, err = repo.FindByID(context.Background(), asset.ID)
	assert.Error(t, err)
}

func Test_memory_repo_should_find_assets_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	userID := "user-1"
	asset1 := NewAsset(userID, Cash, "현금 자산", 1000000.0, "KRW")
	asset2 := NewAsset(userID, Stock, "주식 자산", 2000000.0, "KRW")
	err := repo.Save(context.Background(), asset1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset2)
	assert.NoError(t, err)

	// When
	assets, err := repo.FindByUserID(context.Background(), userID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, assets, 2)
}

func Test_memory_repo_should_find_assets_by_type(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset1 := NewAsset("user-1", Cash, "현금 자산 1", 1000000.0, "KRW")
	asset2 := NewAsset("user-2", Cash, "현금 자산 2", 2000000.0, "KRW")
	err := repo.Save(context.Background(), asset1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset2)
	assert.NoError(t, err)

	// When
	assets, err := repo.FindByType(context.Background(), Cash)

	// Then
	assert.NoError(t, err)
	assert.Len(t, assets, 2)
}

func Test_memory_repo_should_update_asset_amount(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset := NewAsset("user-1", Cash, "현금 자산", 1000000.0, "KRW")
	err := repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	newAmount := 2000000.0
	err = repo.UpdateAmount(context.Background(), asset.ID, newAmount)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, Money{Amount: newAmount, Currency: "KRW"}, found.Amount)
}

func Test_memory_repo_should_save_and_find_transaction_by_id(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	money := NewMoney(500000, "KRW")
	tx, err := NewTransaction("asset-1", Income, money, "급여", "2월 급여")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}

	// When
	err = repo.Save(context.Background(), tx)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, tx.ID, found.ID)
	assert.Equal(t, tx.Amount, found.Amount)
}

func Test_memory_repo_should_find_transactions_by_date_range(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	money := NewMoney(500000, "KRW")
	tx1, err := NewTransaction("asset-1", Income, money, "급여", "2월 급여")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}
	tx1.Date = time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	money2 := NewMoney(300000, "KRW")
	tx2, err := NewTransaction("asset-1", Expense, money2, "식비", "2월 식비")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}
	tx2.Date = time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

	money3 := NewMoney(200000, "KRW")
	tx3, err := NewTransaction("asset-1", Income, money3, "부수입", "2월 부수입")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}
	tx3.Date = time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	repo.Save(context.Background(), tx1)
	repo.Save(context.Background(), tx2)
	repo.Save(context.Background(), tx3)

	// When
	start := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 2, 28, 23, 59, 59, 0, time.UTC)
	found, err := repo.FindByDateRange(context.Background(), start, end)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 2)
}

func Test_memory_repo_should_calculate_total_amount(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	money1 := NewMoney(500000, "KRW")
	tx1, err := NewTransaction("asset-1", Income, money1, "급여", "2월 급여")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}

	money2 := NewMoney(300000, "KRW")
	tx2, err := NewTransaction("asset-1", Expense, money2, "식비", "2월 식비")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}

	money3 := NewMoney(200000, "KRW")
	tx3, err := NewTransaction("asset-1", Income, money3, "부수입", "2월 부수입")
	if err != nil {
		t.Fatalf("거래 생성 중 오류 발생: %v", err)
	}

	repo.Save(context.Background(), tx1)
	repo.Save(context.Background(), tx2)
	repo.Save(context.Background(), tx3)

	// When
	total, err := repo.GetTotalAmount(context.Background(), "asset-1")

	// Then
	assert.NoError(t, err)
	expectedTotal := NewMoney(400000, "KRW")
	assert.Equal(t, expectedTotal, total)
}

func Test_should_save_and_find_portfolio_by_id(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	fixture := NewTestFixture()
	ctx := context.Background()

	// When
	err := repo.Save(ctx, fixture.Portfolios[0])

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Portfolios[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fixture.Portfolios[0], found)
}

func Test_memory_repo_should_find_portfolio_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	err := repo.Save(ctx, fixture.Portfolios[0])
	assert.NoError(t, err)

	// When
	found, err := repo.FindByUserID(ctx, fixture.Portfolios[0].UserID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fixture.Portfolios[0], found)
}

func Test_should_update_portfolio_assets(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	err := repo.Save(ctx, fixture.Portfolios[0])
	assert.NoError(t, err)

	newAssets := []PortfolioAsset{
		{
			AssetID: "asset-1",
			Weight:  0.3,
		},
		{
			AssetID: "asset-2",
			Weight:  0.7,
		},
	}

	// When
	err = repo.UpdateAssets(ctx, fixture.Portfolios[0].ID, newAssets)

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Portfolios[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, newAssets, found.Assets)
}
