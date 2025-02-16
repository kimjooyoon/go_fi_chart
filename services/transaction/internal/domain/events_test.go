package domain

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewTransactionCreatedEvent(t *testing.T) {
	// 테스트 데이터 준비
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

	// 이벤트 생성
	event := NewTransactionCreatedEvent(transaction)

	// 검증
	assert.Equal(t, EventTypeTransactionCreated, event.EventType())
	assert.Equal(t, "transaction", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(TransactionCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, transaction.ID.String(), payload.TransactionID)
	assert.Equal(t, transaction.UserID.String(), payload.UserID)
	assert.Equal(t, transaction.PortfolioID.String(), payload.PortfolioID)
	assert.Equal(t, transaction.AssetID.String(), payload.AssetID)
	assert.Equal(t, string(transaction.Type), payload.Type)
	assert.Equal(t, transaction.Amount.Amount, payload.Amount)
	assert.Equal(t, transaction.Amount.Currency, payload.Currency)
	assert.Equal(t, transaction.Quantity, payload.Quantity)
	assert.Equal(t, transaction.ExecutedPrice.Amount, payload.ExecutedPrice)
	assert.Equal(t, transaction.ExecutedAt, payload.ExecutedAt)
	assert.Equal(t, transaction.CreatedAt, payload.CreatedAt)
}

func TestNewTransactionUpdatedEvent(t *testing.T) {
	// 테스트 데이터 준비
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

	// 거래 업데이트
	newAmount, _ := valueobjects.NewMoney(200.0, "USD")
	prevAmount := transaction.Amount
	prevQuantity := transaction.Quantity
	transaction.Update(Sell, newAmount, 4.0, executedPrice, executedAt)

	// 이벤트 생성
	event := NewTransactionUpdatedEvent(transaction, prevAmount, prevQuantity)

	// 검증
	assert.Equal(t, EventTypeTransactionUpdated, event.EventType())
	assert.Equal(t, "transaction", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(TransactionUpdatedEvent)
	assert.True(t, ok)
	assert.Equal(t, transaction.ID.String(), payload.TransactionID)
	assert.Equal(t, string(transaction.Type), payload.Type)
	assert.Equal(t, transaction.Amount.Amount, payload.Amount)
	assert.Equal(t, transaction.Amount.Currency, payload.Currency)
	assert.Equal(t, transaction.Quantity, payload.Quantity)
	assert.Equal(t, transaction.ExecutedPrice.Amount, payload.ExecutedPrice)
	assert.Equal(t, transaction.ExecutedAt, payload.ExecutedAt)
	assert.Equal(t, prevAmount.Amount, payload.PrevAmount)
	assert.Equal(t, prevAmount.Currency, payload.PrevCurrency)
	assert.Equal(t, prevQuantity, payload.PrevQuantity)
}

func TestNewTransactionDeletedEvent(t *testing.T) {
	// 테스트 데이터 준비
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

	// 이벤트 생성
	event := NewTransactionDeletedEvent(transaction)

	// 검증
	assert.Equal(t, EventTypeTransactionDeleted, event.EventType())
	assert.Equal(t, "transaction", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(TransactionDeletedEvent)
	assert.True(t, ok)
	assert.Equal(t, transaction.ID.String(), payload.TransactionID)
	assert.NotZero(t, payload.DeletedAt)
}
