package domain

import (
	"testing"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestNewAssetCreatedEvent(t *testing.T) {
	// 테스트 데이터 준비
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := NewAsset("user-1", Stock, "테스트 자산", amount)

	// 이벤트 생성
	event := NewAssetCreatedEvent(asset)

	// 검증
	assert.Equal(t, EventTypeAssetCreated, event.EventType())
	assert.Equal(t, "asset", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(AssetCreatedEvent)
	assert.True(t, ok)
	assert.Equal(t, asset.ID, payload.AssetID)
	assert.Equal(t, asset.UserID, payload.UserID)
	assert.Equal(t, string(asset.Type), payload.Type)
	assert.Equal(t, asset.Name, payload.Name)
	assert.Equal(t, asset.Amount.Amount, payload.Amount)
	assert.Equal(t, asset.Amount.Currency, payload.Currency)
	assert.Equal(t, asset.CreatedAt, payload.CreatedAt)
}

func TestNewAssetUpdatedEvent(t *testing.T) {
	// 테스트 데이터 준비
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := NewAsset("user-1", Stock, "테스트 자산", amount)

	// 자산 업데이트
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	asset.Update("업데이트된 자산", Stock, newAmount)

	// 이벤트 생성
	event := NewAssetUpdatedEvent(asset)

	// 검증
	assert.Equal(t, EventTypeAssetUpdated, event.EventType())
	assert.Equal(t, "asset", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(AssetUpdatedEvent)
	assert.True(t, ok)
	assert.Equal(t, asset.ID, payload.AssetID)
	assert.Equal(t, asset.Name, payload.Name)
	assert.Equal(t, asset.UpdatedAt, payload.UpdatedAt)
}

func TestNewAssetAmountChangedEvent(t *testing.T) {
	// 테스트 데이터 준비
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := NewAsset("user-1", Stock, "테스트 자산", amount)

	// 금액 변경
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	prevAmount := asset.Amount
	asset.UpdateAmount(newAmount)

	// 이벤트 생성
	event := NewAssetAmountChangedEvent(asset, prevAmount)

	// 검증
	assert.Equal(t, EventTypeAssetAmountChanged, event.EventType())
	assert.Equal(t, "asset", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(AssetAmountChangedEvent)
	assert.True(t, ok)
	assert.Equal(t, asset.ID, payload.AssetID)
	assert.Equal(t, asset.Amount.Amount, payload.Amount)
	assert.Equal(t, asset.Amount.Currency, payload.Currency)
	assert.Equal(t, prevAmount.Amount, payload.PrevAmount)
	assert.Equal(t, prevAmount.Currency, payload.PrevCurrency)
	assert.Equal(t, asset.UpdatedAt, payload.UpdatedAt)
}

func TestNewAssetDeletedEvent(t *testing.T) {
	// 테스트 데이터 준비
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	asset := NewAsset("user-1", Stock, "테스트 자산", amount)

	// 삭제 표시
	asset.MarkAsDeleted()

	// 이벤트 생성
	event := NewAssetDeletedEvent(asset)

	// 검증
	assert.Equal(t, EventTypeAssetDeleted, event.EventType())
	assert.Equal(t, "asset", event.AggregateType())
	assert.Equal(t, uint(1), event.Version())

	payload, ok := event.Payload().(AssetDeletedEvent)
	assert.True(t, ok)
	assert.Equal(t, asset.ID, payload.AssetID)
	assert.NotZero(t, payload.DeletedAt)
}
