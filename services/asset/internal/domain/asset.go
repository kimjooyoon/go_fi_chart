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

// AssetRepository 자산 저장소 인터페이스입니다.
type AssetRepository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id string) (*Asset, error)
	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) ([]*Asset, error)
	FindByType(ctx context.Context, assetType AssetType) ([]*Asset, error)
}
