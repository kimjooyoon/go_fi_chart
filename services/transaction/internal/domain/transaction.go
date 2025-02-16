package domain

import (
	"context"
	"errors"
	"time"

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

	return &Transaction{
		ID:            uuid.New(),
		UserID:        userID,
		PortfolioID:   portfolioID,
		AssetID:       assetID,
		Type:          transactionType,
		Amount:        amount,
		Quantity:      quantity,
		ExecutedPrice: executedPrice,
		ExecutedAt:    executedAt,
		CreatedAt:     time.Now(),
	}, nil
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
