package asset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestFixture(t *testing.T) {
	// When
	fixture := NewTestFixture()

	// Then
	assert.NotNil(t, fixture)
	assert.NotNil(t, fixture.assets)
	assert.NotNil(t, fixture.transactions)
	assert.NotNil(t, fixture.portfolios)
	assert.Empty(t, fixture.assets)
	assert.Empty(t, fixture.transactions)
	assert.Empty(t, fixture.portfolios)
}

func TestCreateFixture(t *testing.T) {
	// When
	fixture := CreateFixture()

	// Then
	assert.NotNil(t, fixture)
	assert.Len(t, fixture.assets, 2)
	assert.Len(t, fixture.transactions, 3)
	assert.Len(t, fixture.portfolios, 1)
}

func TestTestFixture_GetAssetByID(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var assetID string
	for id := range fixture.assets {
		assetID = id
		break
	}

	// When
	asset := fixture.GetAssetByID(assetID)

	// Then
	assert.NotNil(t, asset)
	assert.Equal(t, assetID, asset.ID)
}

func TestTestFixture_GetTransactionByID(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var txID string
	for id := range fixture.transactions {
		txID = id
		break
	}

	// When
	tx := fixture.GetTransactionByID(txID)

	// Then
	assert.NotNil(t, tx)
	assert.Equal(t, txID, tx.ID)
}

func TestTestFixture_GetPortfolioByID(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var portfolioID string
	for id := range fixture.portfolios {
		portfolioID = id
		break
	}

	// When
	portfolio := fixture.GetPortfolioByID(portfolioID)

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, portfolioID, portfolio.ID)
}

func TestTestFixture_GetAssetsByUserID(t *testing.T) {
	// Given
	fixture := CreateFixture()

	// When
	assets := fixture.GetAssetsByUserID("user-1")

	// Then
	assert.NotEmpty(t, assets)
	for _, asset := range assets {
		assert.Equal(t, "user-1", asset.UserID)
	}
}

func TestTestFixture_GetTransactionsByAssetID(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var assetID string
	for id := range fixture.assets {
		assetID = id
		break
	}

	// When
	transactions := fixture.GetTransactionsByAssetID(assetID)

	// Then
	assert.NotEmpty(t, transactions)
	for _, tx := range transactions {
		assert.Equal(t, assetID, tx.AssetID)
	}
}

func TestTestFixture_GetPortfolioByUserID(t *testing.T) {
	// Given
	fixture := CreateFixture()

	// When
	portfolio := fixture.GetPortfolioByUserID("user-1")

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, "user-1", portfolio.UserID)
}

func Test_fixture_should_create_with_valid_data(t *testing.T) {
	// When
	fixture := CreateFixture()

	// Then
	assert.NotNil(t, fixture)
	assert.NotNil(t, fixture.assets)
	assert.NotNil(t, fixture.transactions)
	assert.NotNil(t, fixture.portfolios)
	assert.Len(t, fixture.assets, 2)
	assert.Len(t, fixture.transactions, 3)
	assert.Len(t, fixture.portfolios, 1)
}

func Test_fixture_should_find_asset_by_id(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var assetID string
	for id := range fixture.assets {
		assetID = id
		break
	}

	// When
	asset := fixture.GetAssetByID(assetID)

	// Then
	assert.NotNil(t, asset)
	assert.Equal(t, assetID, asset.ID)
}

func Test_fixture_should_find_transaction_by_id(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var txID string
	for id := range fixture.transactions {
		txID = id
		break
	}

	// When
	tx := fixture.GetTransactionByID(txID)

	// Then
	assert.NotNil(t, tx)
	assert.Equal(t, txID, tx.ID)
}

func Test_fixture_should_find_portfolio_by_id(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var portfolioID string
	for id := range fixture.portfolios {
		portfolioID = id
		break
	}

	// When
	portfolio := fixture.GetPortfolioByID(portfolioID)

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, portfolioID, portfolio.ID)
}

func Test_fixture_should_find_assets_by_user_id(t *testing.T) {
	// Given
	fixture := CreateFixture()

	// When
	assets := fixture.GetAssetsByUserID("user-1")

	// Then
	assert.NotEmpty(t, assets)
	for _, asset := range assets {
		assert.Equal(t, "user-1", asset.UserID)
	}
}

func Test_fixture_should_find_transactions_by_asset_id(t *testing.T) {
	// Given
	fixture := CreateFixture()
	var assetID string
	for id := range fixture.assets {
		assetID = id
		break
	}

	// When
	transactions := fixture.GetTransactionsByAssetID(assetID)

	// Then
	assert.NotEmpty(t, transactions)
	for _, tx := range transactions {
		assert.Equal(t, assetID, tx.AssetID)
	}
}

func Test_fixture_should_find_portfolio_by_user_id(t *testing.T) {
	// Given
	fixture := CreateFixture()

	// When
	portfolio := fixture.GetPortfolioByUserID("user-1")

	// Then
	assert.NotNil(t, portfolio)
	assert.Equal(t, "user-1", portfolio.UserID)
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
