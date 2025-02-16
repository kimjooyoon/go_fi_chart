package asset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewAsset_should_create_asset_with_valid_data(t *testing.T) {
	// Given
	userID := "test-user"
	assetType := Cash
	name := "현금 자산"
	amount := 1000000.0
	currency := "KRW"

	// When
	asset := NewAsset(userID, assetType, name, amount, currency)

	// Then
	assert.NotEmpty(t, asset.ID)
	assert.Equal(t, userID, asset.UserID)
	assert.Equal(t, assetType, asset.Type)
	assert.Equal(t, name, asset.Name)
	assert.Equal(t, Money{Amount: amount, Currency: currency}, asset.Amount)
	assert.NotZero(t, asset.CreatedAt)
	assert.NotZero(t, asset.UpdatedAt)
	assert.Equal(t, asset.CreatedAt, asset.UpdatedAt)
}

func Test_NewTransaction_should_create_transaction_with_valid_data(t *testing.T) {
	// Given
	assetID := "asset-1"
	transactionType := Income
	amount := 500000.0
	category := "급여"
	description := "3월 급여"

	// When
	tx := NewTransaction(assetID, transactionType, amount, category, description)

	// Then
	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, assetID, tx.AssetID)
	assert.Equal(t, transactionType, tx.Type)
	assert.Equal(t, amount, tx.Amount)
	assert.Equal(t, category, tx.Category)
	assert.Equal(t, description, tx.Description)
	assert.NotZero(t, tx.Date)
	assert.NotZero(t, tx.CreatedAt)
	assert.Equal(t, tx.Date, tx.CreatedAt)
}

func Test_NewPortfolio_should_create_portfolio_with_valid_data(t *testing.T) {
	// Given
	userID := "test-user"
	assets := []PortfolioAsset{
		{
			AssetID: "asset-1",
			Weight:  0.6,
		},
		{
			AssetID: "asset-2",
			Weight:  0.4,
		},
	}

	// When
	portfolio := NewPortfolio(userID, assets)

	// Then
	assert.NotEmpty(t, portfolio.ID)
	assert.Equal(t, userID, portfolio.UserID)
	assert.Equal(t, assets, portfolio.Assets)
	assert.NotZero(t, portfolio.CreatedAt)
	assert.NotZero(t, portfolio.UpdatedAt)
	assert.Equal(t, portfolio.CreatedAt, portfolio.UpdatedAt)
}

func Test_Asset_GetID_should_return_id(t *testing.T) {
	// Given
	asset := &Asset{ID: "test-id"}

	// When
	id := asset.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Asset_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	asset := &Asset{CreatedAt: now}

	// When
	createdAt := asset.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Asset_GetUpdatedAt_should_return_updated_at(t *testing.T) {
	// Given
	now := time.Now()
	asset := &Asset{UpdatedAt: now}

	// When
	updatedAt := asset.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}

func Test_Transaction_GetID_should_return_id(t *testing.T) {
	// Given
	tx := &Transaction{ID: "test-id"}

	// When
	id := tx.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Transaction_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	tx := &Transaction{CreatedAt: now}

	// When
	createdAt := tx.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Transaction_GetUpdatedAt_should_return_date(t *testing.T) {
	// Given
	now := time.Now()
	tx := &Transaction{Date: now}

	// When
	updatedAt := tx.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}

func Test_Portfolio_GetID_should_return_id(t *testing.T) {
	// Given
	portfolio := &Portfolio{ID: "test-id"}

	// When
	id := portfolio.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Portfolio_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	portfolio := &Portfolio{CreatedAt: now}

	// When
	createdAt := portfolio.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Portfolio_GetUpdatedAt_should_return_updated_at(t *testing.T) {
	// Given
	now := time.Now()
	portfolio := &Portfolio{UpdatedAt: now}

	// When
	updatedAt := portfolio.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}
