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
	tx := NewTransaction("asset-1", Income, 500000.0, "급여", "3월 급여")

	// When
	err := repo.Save(context.Background(), tx)
	assert.NoError(t, err)

	found, err := repo.FindByID(context.Background(), tx.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, tx.ID, found.ID)
	assert.Equal(t, tx.Amount, found.Amount)
}

func Test_memory_repo_should_find_transactions_by_date_range(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	now := time.Now()
	tx1 := NewTransaction("asset-1", Income, 500000.0, "급여", "3월 급여")
	tx2 := NewTransaction("asset-1", Expense, 100000.0, "식비", "3월 식비")
	tx1.Date = now.Add(-24 * time.Hour)
	tx2.Date = now

	err := repo.Save(context.Background(), tx1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), tx2)
	assert.NoError(t, err)

	// When
	start := now.Add(-48 * time.Hour)
	end := now.Add(24 * time.Hour)
	transactions, err := repo.FindByDateRange(context.Background(), start, end)

	// Then
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
}

func Test_memory_repo_should_calculate_total_amount(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	assetID := "asset-1"
	tx1 := NewTransaction(assetID, Income, 500000.0, "급여", "3월 급여")
	tx2 := NewTransaction(assetID, Expense, 100000.0, "식비", "3월 식비")

	err := repo.Save(context.Background(), tx1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), tx2)
	assert.NoError(t, err)

	// When
	total, err := repo.GetTotalAmount(context.Background(), assetID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 400000.0, total)
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
