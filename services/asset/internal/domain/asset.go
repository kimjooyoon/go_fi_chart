package domain

import (
	"context"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
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
	events    []events.Event
	mu        sync.RWMutex
}

// NewAsset 새로운 자산을 생성합니다.
func NewAsset(userID string, assetType AssetType, name string, amount valueobjects.Money) *Asset {
	now := time.Now()
	asset := &Asset{
		ID:        uuid.New().String(),
		UserID:    userID,
		Type:      assetType,
		Name:      name,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
		events:    make([]events.Event, 0),
	}

	asset.mu.Lock()
	asset.events = append(asset.events, NewAssetCreatedEvent(asset))
	asset.mu.Unlock()
	return asset
}

// Update 자산 정보를 업데이트합니다.
func (a *Asset) Update(name string, assetType AssetType, amount valueobjects.Money) {
	a.mu.Lock()
	defer a.mu.Unlock()

	prevAmount := a.Amount
	a.Name = name
	a.Type = assetType
	a.Amount = amount
	a.UpdatedAt = time.Now()

	a.events = append(a.events, NewAssetUpdatedEvent(a))
	if !prevAmount.Equals(amount) {
		a.events = append(a.events, NewAssetAmountChangedEvent(a, prevAmount))
	}
}

// UpdateAmount 자산의 금액을 업데이트합니다.
func (a *Asset) UpdateAmount(amount valueobjects.Money) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.Amount.Equals(amount) {
		return
	}

	prevAmount := a.Amount
	a.Amount = amount
	a.UpdatedAt = time.Now()

	a.events = append(a.events, NewAssetAmountChangedEvent(a, prevAmount))
}

// MarkAsDeleted 자산을 삭제 상태로 표시합니다.
func (a *Asset) MarkAsDeleted() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = append(a.events, NewAssetDeletedEvent(a))
}

// Events 발생한 이벤트 목록을 반환합니다.
func (a *Asset) Events() []events.Event {
	a.mu.RLock()
	defer a.mu.RUnlock()
	events := make([]events.Event, len(a.events))
	copy(events, a.events)
	return events
}

// ClearEvents 이벤트 목록을 초기화합니다.
func (a *Asset) ClearEvents() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.events = make([]events.Event, 0)
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

func IsValidAssetType(assetType AssetType) bool {
	switch assetType {
	case Stock, Bond, Cash, RealEstate, Crypto:
		return true
	default:
		return false
	}
}
