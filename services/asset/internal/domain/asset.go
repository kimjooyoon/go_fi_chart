package domain

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
)

// AssetType 자산의 유형을 나타냅니다.
type AssetType string

const (
	Cash       AssetType = "CASH"
	Stock      AssetType = "STOCK"
	Bond       AssetType = "BOND"
	RealEstate AssetType = "REAL_ESTATE"
	Crypto     AssetType = "CRYPTO"
)

// Asset 자산을 나타냅니다.
type Asset struct {
	ID        string
	UserID    string
	Type      AssetType
	Name      string
	Amount    valueobjects.Money
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAsset 새로운 자산을 생성합니다.
func NewAsset(userID string, assetType AssetType, name string, amount valueobjects.Money) *Asset {
	now := time.Now()
	return &Asset{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      assetType,
		Name:      name,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Update 자산 정보를 업데이트합니다.
func (a *Asset) Update(name string, assetType AssetType, amount valueobjects.Money) {
	a.Name = name
	a.Type = assetType
	a.Amount = amount
	a.UpdatedAt = time.Now()
}

// Associate는 자산을 다른 통화로 변환합니다.
func (a *Asset) Associate(targetCurrency string, exchangeRate float64) (*Asset, error) {
	if targetCurrency == a.Amount.Currency {
		return a, nil
	}

	newAmount, err := a.Amount.Multiply(exchangeRate)
	if err != nil {
		return nil, err
	}

	// 금액을 소수점 2자리로 반올림
	roundedAmount := newAmount.Round(2, valueobjects.RoundHalfUp)

	associatedAsset := &Asset{
		ID:        uuid.New().String(),
		UserID:    a.UserID,
		Type:      a.Type,
		Name:      a.Name,
		Amount:    roundedAmount,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return associatedAsset, nil
}

// AssetRepository 자산 저장소 인터페이스입니다.
type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id string) (*Asset, error)
	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) ([]*Asset, error)
	FindByType(ctx context.Context, assetType AssetType) ([]*Asset, error)
}
