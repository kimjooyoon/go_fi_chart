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
	originalUpdatedAt := asset.UpdatedAt

	// When
	time.Sleep(time.Millisecond) // UpdatedAt 변경 확인을 위한 지연
	asset.MarkAsDeleted()

	// Then
	assert.True(t, asset.IsDeleted)
	assert.True(t, asset.UpdatedAt.After(originalUpdatedAt))
	assert.NotNil(t, asset.DeletedAt)

	// 이벤트 검증
	events := asset.Events()
	assert.Len(t, events, 2) // Created + Deleted
	assert.Equal(t, EventTypeAssetDeleted, events[1].EventType())

	// 삭제된 자산 업데이트 시도
	newAmount, _ := valueobjects.NewMoney(2000.0, "USD")
	err := asset.Update("Updated Name", Stock, newAmount)
	assert.Error(t, err)
	assert.Equal(t, ErrAssetDeleted, err)
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

func TestIsValidAssetType(t *testing.T) {
	// Given
	validTypes := []AssetType{Stock, Bond, Cash, RealEstate, Crypto}
	invalidType := AssetType("INVALID")

	// When & Then
	for _, validType := range validTypes {
		assert.True(t, IsValidAssetType(validType), "유효한 자산 타입 %s가 유효하지 않은 것으로 판단됨", validType)
	}
	assert.False(t, IsValidAssetType(invalidType), "유효하지 않은 자산 타입 %s가 유효한 것으로 판단됨", invalidType)
}

func TestAsset_Associate(t *testing.T) {
	t.Run("동일한 통화로 변환", func(t *testing.T) {
		// Given
		asset := createTestAsset()

		// When
		associatedAsset, err := asset.Associate("USD", 1.0)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, asset, associatedAsset)
	})

	t.Run("다른 통화로 변환", func(t *testing.T) {
		// Given
		asset := createTestAsset() // 1000 USD
		exchangeRate := 1300.0     // 1 USD = 1300 KRW

		// When
		associatedAsset, err := asset.Associate("KRW", exchangeRate)

		// Then
		assert.NoError(t, err)
		assert.NotEqual(t, asset.ID, associatedAsset.ID)
		assert.Equal(t, asset.UserID, associatedAsset.UserID)
		assert.Equal(t, asset.Type, associatedAsset.Type)
		assert.Equal(t, asset.Name, associatedAsset.Name)
		assert.Equal(t, "KRW", associatedAsset.Amount.Currency)
		assert.Equal(t, 1300000.0, associatedAsset.Amount.Amount)
	})

	t.Run("금액 곱셈 에러", func(t *testing.T) {
		// 실제 테스트는 skip - 에러 상황 시뮬레이션이 어려움
		t.Skip("Money.Multiply 메서드가 에러를 반환하는 상황을 시뮬레이션하기 어려우므로 스킵합니다.")
	})
}
