package domain

import (
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
)

const (
	EventTypeTransactionCreated = "transaction.created"
	EventTypeTransactionUpdated = "transaction.updated"
	EventTypeTransactionDeleted = "transaction.deleted"
)

// TransactionCreatedEvent는 거래가 생성되었을 때 발생하는 이벤트입니다.
type TransactionCreatedEvent struct {
	events.BaseEvent
	TransactionID string    `json:"transactionId"`
	UserID        string    `json:"userId"`
	PortfolioID   string    `json:"portfolioId"`
	AssetID       string    `json:"assetId"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Quantity      float64   `json:"quantity"`
	ExecutedPrice float64   `json:"executedPrice"`
	ExecutedAt    time.Time `json:"executedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

// NewTransactionCreatedEvent는 새로운 TransactionCreatedEvent를 생성합니다.
func NewTransactionCreatedEvent(transaction *Transaction) events.Event {
	return events.NewEvent(
		EventTypeTransactionCreated,
		transaction.ID,
		"transaction",
		1,
		TransactionCreatedEvent{
			TransactionID: transaction.ID.String(),
			UserID:        transaction.UserID.String(),
			PortfolioID:   transaction.PortfolioID.String(),
			AssetID:       transaction.AssetID.String(),
			Type:          string(transaction.Type),
			Amount:        transaction.Amount.Amount,
			Currency:      transaction.Amount.Currency,
			Quantity:      transaction.Quantity,
			ExecutedPrice: transaction.ExecutedPrice.Amount,
			ExecutedAt:    transaction.ExecutedAt,
			CreatedAt:     transaction.CreatedAt,
		},
		nil,
	)
}

// TransactionUpdatedEvent는 거래가 업데이트되었을 때 발생하는 이벤트입니다.
type TransactionUpdatedEvent struct {
	events.BaseEvent
	TransactionID string    `json:"transactionId"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Quantity      float64   `json:"quantity"`
	ExecutedPrice float64   `json:"executedPrice"`
	ExecutedAt    time.Time `json:"executedAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	PrevAmount    float64   `json:"prevAmount"`
	PrevCurrency  string    `json:"prevCurrency"`
	PrevQuantity  float64   `json:"prevQuantity"`
}

// NewTransactionUpdatedEvent는 새로운 TransactionUpdatedEvent를 생성합니다.
func NewTransactionUpdatedEvent(transaction *Transaction, prevAmount valueobjects.Money, prevQuantity float64) events.Event {
	return events.NewEvent(
		EventTypeTransactionUpdated,
		transaction.ID,
		"transaction",
		1,
		TransactionUpdatedEvent{
			TransactionID: transaction.ID.String(),
			Type:          string(transaction.Type),
			Amount:        transaction.Amount.Amount,
			Currency:      transaction.Amount.Currency,
			Quantity:      transaction.Quantity,
			ExecutedPrice: transaction.ExecutedPrice.Amount,
			ExecutedAt:    transaction.ExecutedAt,
			UpdatedAt:     transaction.UpdatedAt,
			PrevAmount:    prevAmount.Amount,
			PrevCurrency:  prevAmount.Currency,
			PrevQuantity:  prevQuantity,
		},
		nil,
	)
}

// TransactionDeletedEvent는 거래가 삭제되었을 때 발생하는 이벤트입니다.
type TransactionDeletedEvent struct {
	events.BaseEvent
	TransactionID string    `json:"transactionId"`
	DeletedAt     time.Time `json:"deletedAt"`
}

// NewTransactionDeletedEvent는 새로운 TransactionDeletedEvent를 생성합니다.
func NewTransactionDeletedEvent(transaction *Transaction) events.Event {
	return events.NewEvent(
		EventTypeTransactionDeleted,
		transaction.ID,
		"transaction",
		1,
		TransactionDeletedEvent{
			TransactionID: transaction.ID.String(),
			DeletedAt:     time.Now(),
		},
		nil,
	)
}
