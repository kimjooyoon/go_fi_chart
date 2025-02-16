package domain

import (
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/google/uuid"
)

const (
	EventTypeAssetCreated       = "asset.created"
	EventTypeAssetUpdated       = "asset.updated"
	EventTypeAssetDeleted       = "asset.deleted"
	EventTypeAssetAmountChanged = "asset.amount_changed"
)

// AssetCreatedEvent는 자산이 생성되었을 때 발생하는 이벤트입니다.
type AssetCreatedEvent struct {
	events.BaseEvent
	AssetID   string    `json:"assetId"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
}

// NewAssetCreatedEvent는 새로운 AssetCreatedEvent를 생성합니다.
func NewAssetCreatedEvent(asset *Asset) events.Event {
	return events.NewEvent(
		EventTypeAssetCreated,
		uuid.MustParse(asset.ID),
		"asset",
		1,
		AssetCreatedEvent{
			AssetID:   asset.ID,
			UserID:    asset.UserID,
			Type:      string(asset.Type),
			Name:      asset.Name,
			Amount:    asset.Amount.Amount,
			Currency:  asset.Amount.Currency,
			CreatedAt: asset.CreatedAt,
		},
		nil,
	)
}

// AssetUpdatedEvent는 자산이 업데이트되었을 때 발생하는 이벤트입니다.
type AssetUpdatedEvent struct {
	events.BaseEvent
	AssetID   string    `json:"assetId"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// NewAssetUpdatedEvent는 새로운 AssetUpdatedEvent를 생성합니다.
func NewAssetUpdatedEvent(asset *Asset) events.Event {
	return events.NewEvent(
		EventTypeAssetUpdated,
		uuid.MustParse(asset.ID),
		"asset",
		1,
		AssetUpdatedEvent{
			AssetID:   asset.ID,
			Name:      asset.Name,
			UpdatedAt: asset.UpdatedAt,
		},
		nil,
	)
}

// AssetAmountChangedEvent는 자산의 금액이 변경되었을 때 발생하는 이벤트입니다.
type AssetAmountChangedEvent struct {
	events.BaseEvent
	AssetID      string    `json:"assetId"`
	Amount       float64   `json:"amount"`
	Currency     string    `json:"currency"`
	PrevAmount   float64   `json:"prevAmount"`
	PrevCurrency string    `json:"prevCurrency"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// NewAssetAmountChangedEvent는 새로운 AssetAmountChangedEvent를 생성합니다.
func NewAssetAmountChangedEvent(asset *Asset, prevAmount valueobjects.Money) events.Event {
	return events.NewEvent(
		EventTypeAssetAmountChanged,
		uuid.MustParse(asset.ID),
		"asset",
		1,
		AssetAmountChangedEvent{
			AssetID:      asset.ID,
			Amount:       asset.Amount.Amount,
			Currency:     asset.Amount.Currency,
			PrevAmount:   prevAmount.Amount,
			PrevCurrency: prevAmount.Currency,
			UpdatedAt:    asset.UpdatedAt,
		},
		nil,
	)
}

// AssetDeletedEvent는 자산이 삭제되었을 때 발생하는 이벤트입니다.
type AssetDeletedEvent struct {
	events.BaseEvent
	AssetID   string    `json:"assetId"`
	DeletedAt time.Time `json:"deletedAt"`
}

// NewAssetDeletedEvent는 새로운 AssetDeletedEvent를 생성합니다.
func NewAssetDeletedEvent(asset *Asset) events.Event {
	return events.NewEvent(
		EventTypeAssetDeleted,
		uuid.MustParse(asset.ID),
		"asset",
		1,
		AssetDeletedEvent{
			AssetID:   asset.ID,
			DeletedAt: time.Now(),
		},
		nil,
	)
}
