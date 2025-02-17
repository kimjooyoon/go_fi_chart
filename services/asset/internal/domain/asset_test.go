package domain

import (
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/stretchr/testify/assert"
)

func TestNewAsset(t *testing.T) {
	// Given
	userID := "user-123"
	assetType := Stock
	name := "Tesla Stock"
	amount, _ := valueobjects.NewMoney(1000.0, "USD")

	// When
	asset := NewAsset(userID, assetType, name, amount)

	// Then
	assert.NotNil(t, asset)
	assert.NotEmpty(t, asset.ID)
	assert.Equal(t, userID, asset.UserID)
	assert.Equal(t, assetType, asset.Type)
	assert.Equal(t, name, asset.Name)
	assert.Equal(t, amount, asset.Amount)
	assert.NotZero(t, asset.CreatedAt)
	assert.NotZero(t, asset.UpdatedAt)
	assert.Len(t, asset.Events(), 1)

	// 이벤트 검증
	events := asset.Events()
	assert.Equal(t, EventTypeAssetCreated, events[0].EventType())
}

func TestAsset_Update(t *testing.T) {
	// Given
	asset := createTestAsset()
	newName := "Updated Tesla Stock"
	newType := Bond
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	originalUpdatedAt := asset.UpdatedAt

	// When
	time.Sleep(time.Millisecond) // UpdatedAt 변경 확인을 위한 지연
	asset.Update(newName, newType, newAmount)

	// Then
	assert.Equal(t, newName, asset.Name)
	assert.Equal(t, newType, asset.Type)
	assert.Equal(t, newAmount, asset.Amount)
	assert.True(t, asset.UpdatedAt.After(originalUpdatedAt))

	// 이벤트 검증
	events := asset.Events()
	assert.Len(t, events, 3) // Created + Updated + AmountChanged
	assert.Equal(t, EventTypeAssetUpdated, events[1].EventType())
	assert.Equal(t, EventTypeAssetAmountChanged, events[2].EventType())
}

func TestAsset_UpdateAmount(t *testing.T) {
	// Given
	asset := createTestAsset()
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	originalUpdatedAt := asset.UpdatedAt

	// When
	time.Sleep(time.Millisecond) // UpdatedAt 변경 확인을 위한 지연
	asset.UpdateAmount(newAmount)

	// Then
	assert.Equal(t, newAmount, asset.Amount)
	assert.True(t, asset.UpdatedAt.After(originalUpdatedAt))

	// 이벤트 검증
	events := asset.Events()
	assert.Len(t, events, 2) // Created + AmountChanged
	assert.Equal(t, EventTypeAssetAmountChanged, events[1].EventType())
}

func TestAsset_MarkAsDeleted(t *testing.T) {
	// Given
	asset := createTestAsset()

	// When
	asset.MarkAsDeleted()

	// Then
	events := asset.Events()
	assert.Len(t, events, 2) // Created + Deleted
	assert.Equal(t, EventTypeAssetDeleted, events[1].EventType())
}

func TestAsset_ClearEvents(t *testing.T) {
	// Given
	asset := createTestAsset()
	assert.NotEmpty(t, asset.Events())

	// When
	asset.ClearEvents()

	// Then
	assert.Empty(t, asset.Events())
}

func createTestAsset() *Asset {
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	return NewAsset("user-123", Stock, "Tesla Stock", amount)
}
