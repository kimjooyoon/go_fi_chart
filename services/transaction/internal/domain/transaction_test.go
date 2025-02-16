package domain

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, _ := valueobjects.NewMoney(100.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
	executedAt := time.Now()

	transaction, err := NewTransaction(
		userID,
		portfolioID,
		assetID,
		Buy,
		amount,
		2.0,
		executedPrice,
		executedAt,
	)

	assert.NoError(t, err)
	assert.NotNil(t, transaction)
	assert.Equal(t, userID, transaction.UserID)
	assert.Equal(t, portfolioID, transaction.PortfolioID)
	assert.Equal(t, assetID, transaction.AssetID)
	assert.Equal(t, Buy, transaction.Type)
	assert.Equal(t, amount, transaction.Amount)
	assert.Equal(t, 2.0, transaction.Quantity)
	assert.Equal(t, executedPrice, transaction.ExecutedPrice)
	assert.Equal(t, executedAt, transaction.ExecutedAt)
	assert.NotZero(t, transaction.CreatedAt)
	assert.NotZero(t, transaction.UpdatedAt)

	invalidAmount, _ := valueobjects.NewMoney(-100.0, "USD")
	transaction, err = NewTransaction(
		userID,
		portfolioID,
		assetID,
		Buy,
		invalidAmount,
		2.0,
		executedPrice,
		executedAt,
	)

	assert.Error(t, err)
	assert.Nil(t, transaction)
}

func TestTransaction_Validate(t *testing.T) {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, _ := valueobjects.NewMoney(100.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
	executedAt := time.Now()

	transaction, _ := NewTransaction(
		userID,
		portfolioID,
		assetID,
		Buy,
		amount,
		2.0,
		executedPrice,
		executedAt,
	)

	err := transaction.Validate()
	assert.NoError(t, err)

	transaction.Quantity = -1
	err = transaction.Validate()
	assert.Error(t, err)

	transaction.Quantity = 2.0
	transaction.Type = "INVALID"
	err = transaction.Validate()
	assert.Error(t, err)
}

func TestTransaction_CalculateTotalAmount(t *testing.T) {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, _ := valueobjects.NewMoney(100.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
	executedAt := time.Now()

	transaction, _ := NewTransaction(
		userID,
		portfolioID,
		assetID,
		Buy,
		amount,
		2.0,
		executedPrice,
		executedAt,
	)

	total := transaction.CalculateTotalAmount()
	expected, _ := valueobjects.NewMoney(100.0, "USD")
	assert.Equal(t, expected, total)
}

func TestTransaction_Update(t *testing.T) {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, _ := valueobjects.NewMoney(100.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
	executedAt := time.Now()

	transaction, _ := NewTransaction(
		userID,
		portfolioID,
		assetID,
		Buy,
		amount,
		2.0,
		executedPrice,
		executedAt,
	)

	newAmount, _ := valueobjects.NewMoney(200.0, "USD")
	newExecutedPrice, _ := valueobjects.NewMoney(100.0, "USD")
	newExecutedAt := time.Now().Add(time.Hour)

	transaction.Update(Sell, newAmount, 4.0, newExecutedPrice, newExecutedAt)

	assert.Equal(t, Sell, transaction.Type)
	assert.Equal(t, newAmount, transaction.Amount)
	assert.Equal(t, 4.0, transaction.Quantity)
	assert.Equal(t, newExecutedPrice, transaction.ExecutedPrice)
	assert.Equal(t, newExecutedAt, transaction.ExecutedAt)
}
