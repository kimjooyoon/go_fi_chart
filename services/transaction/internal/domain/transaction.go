package domain

import (
	"context"
	"errors"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
)

// TransactionType은 거래 유형을 나타냅니다
type TransactionType string

const (
	Buy  TransactionType = "BUY"
	Sell TransactionType = "SELL"
)

// Transaction은 자산 거래를 나타내는 도메인 모델입니다
type Transaction struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	PortfolioID   uuid.UUID
	AssetID       uuid.UUID
	Type          TransactionType
	Amount        valueobjects.Money
	Quantity      float64
	ExecutedPrice valueobjects.Money
	ExecutedAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	events        []events.Event
}

// NewTransaction은 새로운 거래를 생성합니다
func NewTransaction(
	userID uuid.UUID,
	portfolioID uuid.UUID,
	assetID uuid.UUID,
	transactionType TransactionType,
	amount valueobjects.Money,
	quantity float64,
	executedPrice valueobjects.Money,
	executedAt time.Time,
) (*Transaction, error) {
	if userID == uuid.Nil {
		return nil, errors.New("user ID is required")
	}
	if portfolioID == uuid.Nil {
		return nil, errors.New("portfolio ID is required")
	}
	if assetID == uuid.Nil {
		return nil, errors.New("asset ID is required")
	}
	if quantity <= 0 {
		return nil, errors.New("quantity must be positive")
	}
	if !amount.IsPositive() {
		return nil, errors.New("amount must be positive")
	}
	if !executedPrice.IsPositive() {
		return nil, errors.New("executed price must be positive")
	}

	now := time.Now()
	transaction := &Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		PortfolioID:   portfolioID,
		AssetID:       assetID,
		Type:          transactionType,
		Amount:        amount,
		Quantity:      quantity,
		ExecutedPrice: executedPrice,
		ExecutedAt:    executedAt,
		CreatedAt:     now,
		UpdatedAt:     now,
		events:        make([]events.Event, 0),
	}

	transaction.events = append(transaction.events, NewTransactionCreatedEvent(transaction))
	return transaction, nil
}

// Validate는 거래의 유효성을 검증합니다
func (t *Transaction) Validate() error {
	if t.Quantity <= 0 {
		return errors.New("quantity must be positive")
	}
	if t.Type != Buy && t.Type != Sell {
		return errors.New("invalid transaction type")
	}
	return nil
}

// CalculateTotalAmount는 거래의 총 금액을 계산합니다
func (t *Transaction) CalculateTotalAmount() valueobjects.Money {
	return t.Amount
}

// Update는 거래 정보를 업데이트합니다
func (t *Transaction) Update(
	transactionType TransactionType,
	amount valueobjects.Money,
	quantity float64,
	executedPrice valueobjects.Money,
	executedAt time.Time,
) {
	prevAmount := t.Amount
	prevQuantity := t.Quantity

	t.Type = transactionType
	t.Amount = amount
	t.Quantity = quantity
	t.ExecutedPrice = executedPrice
	t.ExecutedAt = executedAt
	t.UpdatedAt = time.Now()

	t.events = append(t.events, NewTransactionUpdatedEvent(t, prevAmount, prevQuantity))
}

// MarkAsDeleted는 거래를 삭제 상태로 표시합니다
func (t *Transaction) MarkAsDeleted() {
	t.events = append(t.events, NewTransactionDeletedEvent(t))
}

// Events는 발생한 이벤트 목록을 반환합니다
func (t *Transaction) Events() []events.Event {
	return t.events
}

// ClearEvents는 이벤트 목록을 초기화합니다
func (t *Transaction) ClearEvents() {
	t.events = make([]events.Event, 0)
}

// TransactionRepository는 거래 저장소 인터페이스를 정의합니다
type TransactionRepository interface {
	Save(ctx context.Context, transaction *Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Transaction, error)
	FindByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]*Transaction, error)
	FindByAssetID(ctx context.Context, assetID uuid.UUID) ([]*Transaction, error)
	Update(ctx context.Context, transaction *Transaction) error
	Delete(ctx context.Context, id uuid.UUID) error
}
