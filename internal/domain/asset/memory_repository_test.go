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
	asset, err := NewAsset("test-user", Cash, "Test Asset", 1000000, "KRW")
	assert.NoError(t, err)

	// When
	err = repo.Save(context.Background(), asset)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, asset, found)
}

func Test_memory_repo_should_update_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset, err := NewAsset("test-user", Cash, "Test Asset", 1000000, "KRW")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	asset.Name = "Updated Asset"
	err = repo.Update(context.Background(), asset)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Asset", found.Name)
}

func Test_memory_repo_should_delete_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset, err := NewAsset("test-user", Cash, "Test Asset", 1000000, "KRW")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	err = repo.Delete(context.Background(), asset.ID)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.Error(t, err)
	assert.Nil(t, found)
}

func Test_memory_repo_should_find_assets_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset1, err := NewAsset("test-user", Cash, "Test Asset 1", 1000000, "KRW")
	assert.NoError(t, err)
	asset2, err := NewAsset("test-user", Stock, "Test Asset 2", 2000000, "KRW")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset2)
	assert.NoError(t, err)

	// When
	assets, err := repo.FindByUserID(context.Background(), "test-user")

	// Then
	assert.NoError(t, err)
	assert.Len(t, assets, 2)
}

func Test_memory_repo_should_find_assets_by_type(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset1, err := NewAsset("test-user", Cash, "Test Asset 1", 1000000, "KRW")
	assert.NoError(t, err)
	asset2, err := NewAsset("test-user", Stock, "Test Asset 2", 2000000, "KRW")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset2)
	assert.NoError(t, err)

	// When
	assets, err := repo.FindByType(context.Background(), Cash)

	// Then
	assert.NoError(t, err)
	assert.Len(t, assets, 1)
	assert.Equal(t, Cash, assets[0].Type)
}

func Test_memory_repo_should_update_asset_amount(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	asset, err := NewAsset("test-user", Cash, "Test Asset", 1000000, "KRW")
	assert.NoError(t, err)
	err = repo.Save(context.Background(), asset)
	assert.NoError(t, err)

	// When
	money := NewTestMoney(2000000, "KRW")
	err = repo.UpdateAmount(context.Background(), asset.ID, money)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindByID(context.Background(), asset.ID)
	assert.NoError(t, err)
	assert.Equal(t, money, found.Amount)
}

func Test_memory_repo_should_save_and_find_transaction_by_id(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	tx := NewTestTransaction()

	// When
	err := repo.SaveTransaction(context.Background(), tx)

	// Then
	assert.NoError(t, err)
	found, err := repo.FindTransactionByID(context.Background(), tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, tx, found)
}

func Test_memory_repo_should_find_transactions_by_date_range(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	tx1 := NewTestTransaction()
	tx2 := NewTestTransaction()
	err := repo.SaveTransaction(context.Background(), tx1)
	assert.NoError(t, err)
	err = repo.SaveTransaction(context.Background(), tx2)
	assert.NoError(t, err)

	// When
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)
	transactions, err := repo.FindTransactionsByDateRange(context.Background(), start, end)

	// Then
	assert.NoError(t, err)
	assert.Len(t, transactions, 2)
}

func Test_memory_repo_should_calculate_total_amount(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	tx1 := NewTestTransaction()
	tx2 := NewTestTransaction()
	err := repo.SaveTransaction(context.Background(), tx1)
	assert.NoError(t, err)
	err = repo.SaveTransaction(context.Background(), tx2)
	assert.NoError(t, err)

	// When
	total, err := repo.CalculateTotalAmount(context.Background(), "KRW")

	// Then
	assert.NoError(t, err)
	expected := NewTestMoney(1000000, "KRW")
	assert.Equal(t, expected, total)
}

func Test_should_save_and_find_portfolio_by_id(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	portfolio := NewTestPortfolio()
	ctx := context.Background()

	// When
	err := repo.Save(ctx, portfolio)

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, portfolio.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, portfolio, found)
}

func Test_memory_repo_should_find_portfolio_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	portfolio := NewTestPortfolio()
	ctx := context.Background()
	err := repo.Save(ctx, portfolio)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByUserID(ctx, portfolio.UserID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, portfolio, found)
}

func Test_should_update_portfolio_assets(t *testing.T) {
	// Given
	repo := NewMemoryPortfolioRepository()
	portfolio := NewTestPortfolio()
	ctx := context.Background()
	err := repo.Save(ctx, portfolio)
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
	err = repo.UpdateAssets(ctx, portfolio.ID, newAssets)

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, portfolio.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, newAssets, found.Assets)
}
