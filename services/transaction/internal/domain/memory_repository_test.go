package domain

import (
	"context"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func createTestTransaction() *Transaction {
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
	return transaction
}

func TestNewMemoryTransactionRepository(t *testing.T) {
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.transactions)
	assert.Empty(t, repo.transactions)
	assert.Equal(t, eventBus, repo.eventBus)
}

func TestMemoryTransactionRepository_Save(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	transaction := createTestTransaction()

	// When
	err := repo.Save(context.Background(), transaction)

	// Then
	assert.NoError(t, err)
	saved, err := repo.FindByID(context.Background(), transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction, saved)
}

func TestMemoryTransactionRepository_FindByID(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	transaction := createTestTransaction()
	err := repo.Save(context.Background(), transaction)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByID(context.Background(), transaction.ID)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, transaction, found)
}

func TestMemoryTransactionRepository_FindByUserID(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	userID := uuid.New()
	transaction1 := createTestTransaction()
	transaction2 := createTestTransaction()
	transaction1.UserID = userID
	transaction2.UserID = userID
	err := repo.Save(context.Background(), transaction1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), transaction2)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByUserID(context.Background(), userID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	assert.Contains(t, found, transaction1)
	assert.Contains(t, found, transaction2)
}

func TestMemoryTransactionRepository_FindByPortfolioID(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	portfolioID := uuid.New()
	transaction1 := createTestTransaction()
	transaction2 := createTestTransaction()
	transaction1.PortfolioID = portfolioID
	transaction2.PortfolioID = portfolioID
	err := repo.Save(context.Background(), transaction1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), transaction2)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByPortfolioID(context.Background(), portfolioID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	assert.Contains(t, found, transaction1)
	assert.Contains(t, found, transaction2)
}

func TestMemoryTransactionRepository_FindByAssetID(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	assetID := uuid.New()
	transaction1 := createTestTransaction()
	transaction2 := createTestTransaction()
	transaction1.AssetID = assetID
	transaction2.AssetID = assetID
	err := repo.Save(context.Background(), transaction1)
	assert.NoError(t, err)
	err = repo.Save(context.Background(), transaction2)
	assert.NoError(t, err)

	// When
	found, err := repo.FindByAssetID(context.Background(), assetID)

	// Then
	assert.NoError(t, err)
	assert.Len(t, found, 2)
	assert.Contains(t, found, transaction1)
	assert.Contains(t, found, transaction2)
}

func TestMemoryTransactionRepository_Update(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	transaction := createTestTransaction()
	err := repo.Save(context.Background(), transaction)
	assert.NoError(t, err)

	// When
	newAmount, _ := valueobjects.NewMoney(200.0, "USD")
	transaction.Update(Sell, newAmount, 4.0, transaction.ExecutedPrice, transaction.ExecutedAt)
	err = repo.Update(context.Background(), transaction)

	// Then
	assert.NoError(t, err)
	updated, err := repo.FindByID(context.Background(), transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction, updated)
}

func TestMemoryTransactionRepository_Delete(t *testing.T) {
	// Given
	eventBus := events.NewSimplePublisher()
	repo := NewMemoryTransactionRepository(eventBus)
	transaction := createTestTransaction()
	err := repo.Save(context.Background(), transaction)
	assert.NoError(t, err)

	// When
	err = repo.Delete(context.Background(), transaction.ID)

	// Then
	assert.NoError(t, err)
	_, err = repo.FindByID(context.Background(), transaction.ID)
	assert.Error(t, err)
}
