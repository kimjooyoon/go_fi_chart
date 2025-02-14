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
	fixture := NewTestFixture()
	ctx := context.Background()

	// When
	err := repo.Save(ctx, fixture.Assets[0])

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Assets[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fixture.Assets[0], found)
}

func Test_memory_repo_should_update_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	err := repo.Save(ctx, fixture.Assets[0])
	assert.NoError(t, err)

	// When
	fixture.Assets[0].Amount = 2000000
	err = repo.Update(ctx, fixture.Assets[0])

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Assets[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, float64(2000000), found.Amount)
}

func Test_memory_repo_should_delete_asset(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	err := repo.Save(ctx, fixture.Assets[0])
	assert.NoError(t, err)

	// When
	err = repo.Delete(ctx, fixture.Assets[0].ID)

	// Then
	assert.NoError(t, err)

	// When
	_, err = repo.FindByID(ctx, fixture.Assets[0].ID)

	// Then
	assert.Error(t, err)
}

func Test_memory_repo_should_find_assets_by_user_id(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	for _, asset := range fixture.Assets {
		err := repo.Save(ctx, asset)
		assert.NoError(t, err)
	}

	// When
	found, err := repo.FindByUserID(ctx, fixture.Assets[0].UserID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 3)
}

func Test_memory_repo_should_find_assets_by_type(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	for _, asset := range fixture.Assets {
		err := repo.Save(ctx, asset)
		assert.NoError(t, err)
	}

	// When
	found, err := repo.FindByType(ctx, Cash)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 1)
	assert.Equal(t, Cash, found[0].Type)
}

func Test_memory_repo_should_update_asset_amount(t *testing.T) {
	// Given
	repo := NewMemoryAssetRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	err := repo.Save(ctx, fixture.Assets[0])
	assert.NoError(t, err)

	// When
	err = repo.UpdateAmount(ctx, fixture.Assets[0].ID, 2000000)

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Assets[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, float64(2000000), found.Amount)
}

func Test_memory_repo_should_save_and_find_transaction_by_id(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	fixture := NewTestFixture()
	ctx := context.Background()

	// When
	err := repo.Save(ctx, fixture.Transactions[0])

	// Then
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(ctx, fixture.Transactions[0].ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, fixture.Transactions[0], found)
}

func Test_memory_repo_should_find_transactions_by_date_range(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	for _, tx := range fixture.Transactions {
		err := repo.Save(ctx, tx)
		assert.NoError(t, err)
	}

	start := time.Now().Add(-24 * time.Hour)
	end := time.Now().Add(24 * time.Hour)

	// When
	found, err := repo.FindByDateRange(ctx, start, end)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 3)
}

func Test_memory_repo_should_calculate_total_amount(t *testing.T) {
	// Given
	repo := NewMemoryTransactionRepository()
	fixture := NewTestFixture()
	ctx := context.Background()
	for _, tx := range fixture.Transactions {
		err := repo.Save(ctx, tx)
		assert.NoError(t, err)
	}

	// When
	total, err := repo.GetTotalAmount(ctx, fixture.Transactions[0].AssetID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, float64(400000), total) // 500000(Income) - 100000(Expense)
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
