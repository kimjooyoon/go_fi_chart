package domain

import (
	"context"
	"testing"
	"time"

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
	repo := NewMemoryTransactionRepository()

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.transactions)
	assert.Empty(t, repo.transactions)
}

func TestMemoryTransactionRepository_Save(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	err := repo.Save(ctx, transaction)
	assert.NoError(t, err)
	assert.Len(t, repo.transactions, 1)

	err = repo.Save(ctx, transaction)
	assert.Error(t, err)
}

func TestMemoryTransactionRepository_FindByID(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	found, err := repo.FindByID(ctx, transaction.ID)
	assert.NoError(t, err)
	assert.Equal(t, transaction, found)

	notFound, err := repo.FindByID(ctx, uuid.New())
	assert.Error(t, err)
	assert.Nil(t, notFound)
}

func TestMemoryTransactionRepository_FindByUserID(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	transactions, err := repo.FindByUserID(ctx, transaction.UserID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, transaction, transactions[0])

	emptyTransactions, err := repo.FindByUserID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Empty(t, emptyTransactions)
}

func TestMemoryTransactionRepository_FindByPortfolioID(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	transactions, err := repo.FindByPortfolioID(ctx, transaction.PortfolioID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, transaction, transactions[0])

	emptyTransactions, err := repo.FindByPortfolioID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Empty(t, emptyTransactions)
}

func TestMemoryTransactionRepository_FindByAssetID(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	transactions, err := repo.FindByAssetID(ctx, transaction.AssetID)
	assert.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, transaction, transactions[0])

	emptyTransactions, err := repo.FindByAssetID(ctx, uuid.New())
	assert.NoError(t, err)
	assert.Empty(t, emptyTransactions)
}

func TestMemoryTransactionRepository_Update(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	newAmount, _ := valueobjects.NewMoney(200.0, "USD")
	transaction.Update(Sell, newAmount, 4.0, transaction.ExecutedPrice, transaction.ExecutedAt)

	err := repo.Update(ctx, transaction)
	assert.NoError(t, err)

	found, _ := repo.FindByID(ctx, transaction.ID)
	assert.Equal(t, Sell, found.Type)
	assert.Equal(t, newAmount, found.Amount)
	assert.Equal(t, 4.0, found.Quantity)

	notFoundTransaction := createTestTransaction()
	err = repo.Update(ctx, notFoundTransaction)
	assert.Error(t, err)
}

func TestMemoryTransactionRepository_Delete(t *testing.T) {
	repo := NewMemoryTransactionRepository()
	ctx := context.Background()
	transaction := createTestTransaction()

	repo.Save(ctx, transaction)

	err := repo.Delete(ctx, transaction.ID)
	assert.NoError(t, err)
	assert.Empty(t, repo.transactions)

	err = repo.Delete(ctx, transaction.ID)
	assert.Error(t, err)
}
