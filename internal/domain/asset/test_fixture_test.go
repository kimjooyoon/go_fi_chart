package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_fixture_should_create_with_valid_data(t *testing.T) {
	// When
	fixture := NewTestFixture()

	// Then
	assert.NotNil(t, fixture)
	assert.Len(t, fixture.Assets, 3)
	assert.Len(t, fixture.Transactions, 3)
	assert.Len(t, fixture.Portfolios, 1)
}

func Test_fixture_should_find_asset_by_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	asset := fixture.GetAssetByID("asset-1")

	// Then
	assert.NotNil(t, asset)
	assert.Equal(t, "asset-1", asset.ID)
	assert.Equal(t, Cash, asset.Type)
	assert.Equal(t, "현금 자산", asset.Name)
}

func Test_fixture_should_find_transaction_by_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	tx := fixture.GetTransactionByID("tx-1")

	// Then
	assert.NotNil(t, tx)
	assert.Equal(t, "tx-1", tx.ID)
	assert.Equal(t, Income, tx.Type)
	assert.Equal(t, "급여", tx.Category)
}

func Test_fixture_should_find_portfolio_by_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	portfolio := fixture.GetPortfolioByID("portfolio-1")

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, "portfolio-1", portfolio.ID)
	assert.Len(t, portfolio.Assets, 3)
}

func Test_fixture_should_find_assets_by_user_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	assets := fixture.GetAssetsByUserID("test-user-1")

	// Then
	assert.Len(t, assets, 3)
	for _, asset := range assets {
		assert.Equal(t, "test-user-1", asset.UserID)
	}
}

func Test_fixture_should_find_transactions_by_asset_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	transactions := fixture.GetTransactionsByAssetID("asset-1")

	// Then
	assert.Len(t, transactions, 2)
	for _, tx := range transactions {
		assert.Equal(t, "asset-1", tx.AssetID)
	}
}

func Test_fixture_should_find_portfolio_by_user_id(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When
	portfolio := fixture.GetPortfolioByUserID("test-user-1")

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, "test-user-1", portfolio.UserID)

	// 포트폴리오 비중 합이 1.0(100%)인지 검증
	var totalWeight float64
	for _, asset := range portfolio.Assets {
		totalWeight += asset.Weight
	}
	assert.Equal(t, 1.0, totalWeight)
}

func Test_fixture_should_return_nil_when_not_found(t *testing.T) {
	// Given
	fixture := NewTestFixture()

	// When & Then
	assert.Nil(t, fixture.GetAssetByID("non-existent"))
	assert.Nil(t, fixture.GetTransactionByID("non-existent"))
	assert.Nil(t, fixture.GetPortfolioByID("non-existent"))
	assert.Nil(t, fixture.GetPortfolioByUserID("non-existent"))

	assert.Empty(t, fixture.GetAssetsByUserID("non-existent"))
	assert.Empty(t, fixture.GetTransactionsByAssetID("non-existent"))
}
