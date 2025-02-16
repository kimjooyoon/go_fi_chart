package asset

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_NewAsset_should_create_asset_with_valid_data(t *testing.T) {
	// Given
	userID := "test-user"
	assetType := Cash
	name := "현금 자산"
	amount := 1000000.0
	currency := "KRW"

	// When
	asset, err := NewAsset(userID, assetType, name, amount, currency)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, asset.ID)
	assert.Equal(t, userID, asset.UserID)
	assert.Equal(t, assetType, asset.Type)
	assert.Equal(t, name, asset.Name)
	money, err := NewMoney(amount, currency)
	assert.NoError(t, err)
	assert.Equal(t, money, asset.Amount)
	assert.NotZero(t, asset.CreatedAt)
	assert.NotZero(t, asset.UpdatedAt)
	assert.Equal(t, asset.CreatedAt, asset.UpdatedAt)
}

func Test_NewTransaction_should_create_transaction_with_valid_data(t *testing.T) {
	// Given
	assetID := "test-asset-1"
	txType := Income
	money := NewTestMoney(500000, "KRW")
	category := "급여"
	description := "3월 급여"

	// When
	tx, err := NewTransaction(assetID, txType, money, category, description)

	// Then
	assert.NoError(t, err)
	assert.NotEmpty(t, tx.ID)
	assert.Equal(t, assetID, tx.AssetID)
	assert.Equal(t, txType, tx.Type)
	assert.Equal(t, money, tx.Amount)
	assert.Equal(t, category, tx.Category)
	assert.Equal(t, description, tx.Description)
	assert.False(t, tx.Date.IsZero())
	assert.False(t, tx.CreatedAt.IsZero())
}

func Test_NewPortfolio_should_create_portfolio_with_valid_data(t *testing.T) {
	// Given
	userID := "test-user"
	assets := []PortfolioAsset{
		{
			AssetID: "asset-1",
			Weight:  0.6,
		},
		{
			AssetID: "asset-2",
			Weight:  0.4,
		},
	}

	// When
	portfolio := NewPortfolio(userID, assets)

	// Then
	assert.NotEmpty(t, portfolio.ID)
	assert.Equal(t, userID, portfolio.UserID)
	assert.Equal(t, assets, portfolio.Assets)
	assert.NotZero(t, portfolio.CreatedAt)
	assert.NotZero(t, portfolio.UpdatedAt)
	assert.Equal(t, portfolio.CreatedAt, portfolio.UpdatedAt)
}

func Test_Asset_GetID_should_return_id(t *testing.T) {
	// Given
	asset := &Asset{ID: "test-id"}

	// When
	id := asset.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Asset_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	asset := &Asset{CreatedAt: now}

	// When
	createdAt := asset.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Asset_GetUpdatedAt_should_return_updated_at(t *testing.T) {
	// Given
	now := time.Now()
	asset := &Asset{UpdatedAt: now}

	// When
	updatedAt := asset.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}

func Test_Transaction_GetID_should_return_id(t *testing.T) {
	// Given
	tx := &Transaction{ID: "test-id"}

	// When
	id := tx.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Transaction_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	tx := &Transaction{CreatedAt: now}

	// When
	createdAt := tx.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Transaction_GetUpdatedAt_should_return_date(t *testing.T) {
	// Given
	now := time.Now()
	tx := &Transaction{Date: now}

	// When
	updatedAt := tx.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}

func Test_Portfolio_GetID_should_return_id(t *testing.T) {
	// Given
	portfolio := &Portfolio{ID: "test-id"}

	// When
	id := portfolio.GetID()

	// Then
	assert.Equal(t, "test-id", id)
}

func Test_Portfolio_GetCreatedAt_should_return_created_at(t *testing.T) {
	// Given
	now := time.Now()
	portfolio := &Portfolio{CreatedAt: now}

	// When
	createdAt := portfolio.GetCreatedAt()

	// Then
	assert.Equal(t, now, createdAt)
}

func Test_Portfolio_GetUpdatedAt_should_return_updated_at(t *testing.T) {
	// Given
	now := time.Now()
	portfolio := &Portfolio{UpdatedAt: now}

	// When
	updatedAt := portfolio.GetUpdatedAt()

	// Then
	assert.Equal(t, now, updatedAt)
}

func TestNewPercentage(t *testing.T) {
	tests := []struct {
		name    string
		value   float64
		wantErr bool
	}{
		{
			name:    "유효한 퍼센트 값",
			value:   50.0,
			wantErr: false,
		},
		{
			name:    "0% 값",
			value:   0.0,
			wantErr: false,
		},
		{
			name:    "100% 값",
			value:   100.0,
			wantErr: false,
		},
		{
			name:    "음수 값",
			value:   -1.0,
			wantErr: true,
		},
		{
			name:    "100% 초과 값",
			value:   101.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPercentage(tt.value)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, got.Value)
			}
		})
	}
}

func TestPercentage_Operations(t *testing.T) {
	p50, _ := NewPercentage(50.0)
	p25, _ := NewPercentage(25.0)

	t.Run("Add", func(t *testing.T) {
		sum, err := p50.Add(p25)
		assert.NoError(t, err)
		assert.Equal(t, 75.0, sum.Value)

		// 100% 초과 테스트
		_, err = sum.Add(p50)
		assert.Error(t, err)
	})

	t.Run("Subtract", func(t *testing.T) {
		diff, err := p50.Subtract(p25)
		assert.NoError(t, err)
		assert.Equal(t, 25.0, diff.Value)

		// 음수 결과 테스트
		_, err = p25.Subtract(p50)
		assert.Error(t, err)
	})

	t.Run("Multiply", func(t *testing.T) {
		result, err := p50.Multiply(0.5)
		assert.NoError(t, err)
		assert.Equal(t, 25.0, result.Value)

		// 100% 초과 테스트
		_, err = p50.Multiply(3.0)
		assert.Error(t, err)
	})
}

func TestPercentage_Conversions(t *testing.T) {
	p50, _ := NewPercentage(50.0)

	t.Run("ToDecimal", func(t *testing.T) {
		decimal := p50.ToDecimal()
		assert.Equal(t, 0.5, decimal)
	})

	t.Run("FromDecimal", func(t *testing.T) {
		p, err := FromDecimal(0.5)
		assert.NoError(t, err)
		assert.Equal(t, 50.0, p.Value)

		// 1.0 초과 테스트
		_, err = FromDecimal(1.5)
		assert.Error(t, err)
	})
}

func TestPercentage_Checks(t *testing.T) {
	p0, _ := NewPercentage(0.0)
	p100, _ := NewPercentage(100.0)
	p50, _ := NewPercentage(50.0)

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, p0.IsZero())
		assert.False(t, p50.IsZero())
		assert.False(t, p100.IsZero())
	})

	t.Run("IsComplete", func(t *testing.T) {
		assert.False(t, p0.IsComplete())
		assert.False(t, p50.IsComplete())
		assert.True(t, p100.IsComplete())
	})
}

func TestNewTimeRange(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tests := []struct {
		name    string
		start   time.Time
		end     time.Time
		wantErr bool
	}{
		{
			name:    "유효한 시간 범위",
			start:   now,
			end:     later,
			wantErr: false,
		},
		{
			name:    "시작 시간이 종료 시간보다 늦은 경우",
			start:   later,
			end:     now,
			wantErr: true,
		},
		{
			name:    "시작 시간과 종료 시간이 같은 경우",
			start:   now,
			end:     now,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTimeRange(tt.start, tt.end)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.start, got.Start)
				assert.Equal(t, tt.end, got.End)
			}
		})
	}
}

func TestTimeRange_Duration(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)

	tr, _ := NewTimeRange(now, later)
	duration := tr.Duration()

	assert.Equal(t, 24*time.Hour, duration)
}

func TestTimeRange_Contains(t *testing.T) {
	now := time.Now()
	middle := now.Add(12 * time.Hour)
	later := now.Add(24 * time.Hour)
	outside := now.Add(48 * time.Hour)

	tr, _ := NewTimeRange(now, later)

	tests := []struct {
		name     string
		time     time.Time
		expected bool
	}{
		{
			name:     "시작 시간",
			time:     now,
			expected: true,
		},
		{
			name:     "중간 시간",
			time:     middle,
			expected: true,
		},
		{
			name:     "종료 시간",
			time:     later,
			expected: true,
		},
		{
			name:     "범위 밖 시간",
			time:     outside,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tr.Contains(tt.time))
		})
	}
}

func TestTimeRange_Overlaps(t *testing.T) {
	now := time.Now()
	middle := now.Add(12 * time.Hour)
	later := now.Add(24 * time.Hour)
	outside := now.Add(48 * time.Hour)

	tr, _ := NewTimeRange(now, later)

	tests := []struct {
		name     string
		other    TimeRange
		expected bool
	}{
		{
			name:     "완전히 포함되는 경우",
			other:    TimeRange{Start: now.Add(6 * time.Hour), End: now.Add(18 * time.Hour)},
			expected: true,
		},
		{
			name:     "부분적으로 겹치는 경우",
			other:    TimeRange{Start: middle, End: outside},
			expected: true,
		},
		{
			name:     "겹치지 않는 경우",
			other:    TimeRange{Start: later.Add(time.Hour), End: outside},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tr.Overlaps(tt.other))
		})
	}
}

func TestTimeRange_Split(t *testing.T) {
	now := time.Now()
	later := now.Add(10 * time.Hour)
	tr, _ := NewTimeRange(now, later)

	tests := []struct {
		name         string
		interval     time.Duration
		expectedLen  int
		expectedLast time.Time
	}{
		{
			name:         "2시간 간격으로 분할",
			interval:     2 * time.Hour,
			expectedLen:  5,
			expectedLast: later,
		},
		{
			name:         "음수 간격",
			interval:     -1 * time.Hour,
			expectedLen:  1,
			expectedLast: later,
		},
		{
			name:         "전체 기간보다 큰 간격",
			interval:     24 * time.Hour,
			expectedLen:  1,
			expectedLast: later,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranges := tr.Split(tt.interval)
			assert.Equal(t, tt.expectedLen, len(ranges))
			assert.Equal(t, now, ranges[0].Start)
			assert.Equal(t, tt.expectedLast, ranges[len(ranges)-1].End)
		})
	}
}

func TestTimeRange_Extend(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)
	tr, _ := NewTimeRange(now, later)

	extended, err := tr.Extend(24 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, now, extended.Start)
	assert.Equal(t, later.Add(24*time.Hour), extended.End)
}

func TestTimeRange_Shift(t *testing.T) {
	now := time.Now()
	later := now.Add(24 * time.Hour)
	tr, _ := NewTimeRange(now, later)

	shifted := tr.Shift(24 * time.Hour)
	assert.Equal(t, now.Add(24*time.Hour), shifted.Start)
	assert.Equal(t, later.Add(24*time.Hour), shifted.End)
}

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  bool
	}{
		{
			name:     "유효한 금액과 통화",
			amount:   1000.0,
			currency: "KRW",
			wantErr:  false,
		},
		{
			name:     "음수 금액",
			amount:   -1000.0,
			currency: "KRW",
			wantErr:  true,
		},
		{
			name:     "빈 통화",
			amount:   1000.0,
			currency: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMoney(tt.amount, tt.currency)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.amount, got.Amount)
				assert.Equal(t, tt.currency, got.Currency)
			}
		})
	}
}

func TestMoney_Operations(t *testing.T) {
	krw1000 := NewTestMoney(1000.0, "KRW")
	krw500 := NewTestMoney(500.0, "KRW")
	usd100 := NewTestMoney(100.0, "USD")

	t.Run("Add", func(t *testing.T) {
		sum, err := krw1000.Add(krw500)
		assert.NoError(t, err)
		assert.Equal(t, 1500.0, sum.Amount)
		assert.Equal(t, "KRW", sum.Currency)

		// 다른 통화 더하기 시도
		_, err = krw1000.Add(usd100)
		assert.Error(t, err)
	})

	t.Run("Subtract", func(t *testing.T) {
		diff, err := krw1000.Subtract(krw500)
		assert.NoError(t, err)
		assert.Equal(t, 500.0, diff.Amount)
		assert.Equal(t, "KRW", diff.Currency)

		// 다른 통화 빼기 시도
		_, err = krw1000.Subtract(usd100)
		assert.Error(t, err)

		// 큰 금액 빼기 시도
		_, err = krw500.Subtract(krw1000)
		assert.Error(t, err)
	})

	t.Run("Multiply", func(t *testing.T) {
		result, err := krw1000.Multiply(2.0)
		assert.NoError(t, err)
		assert.Equal(t, 2000.0, result.Amount)
		assert.Equal(t, "KRW", result.Currency)

		// 음수 곱하기 시도
		_, err = krw1000.Multiply(-2.0)
		assert.Error(t, err)
	})

	t.Run("Divide", func(t *testing.T) {
		result, err := krw1000.Divide(2.0)
		assert.NoError(t, err)
		assert.Equal(t, 500.0, result.Amount)
		assert.Equal(t, "KRW", result.Currency)

		// 0으로 나누기 시도
		_, err = krw1000.Divide(0)
		assert.Error(t, err)

		// 음수로 나누기 시도
		_, err = krw1000.Divide(-2.0)
		assert.Error(t, err)
	})
}

func TestMoney_Comparisons(t *testing.T) {
	krw1000 := NewTestMoney(1000.0, "KRW")
	krw500 := NewTestMoney(500.0, "KRW")
	krw1000_2 := NewTestMoney(1000.0, "KRW")
	usd100 := NewTestMoney(100.0, "USD")

	t.Run("Equals", func(t *testing.T) {
		assert.True(t, krw1000.Equals(krw1000_2))
		assert.False(t, krw1000.Equals(krw500))
		assert.False(t, krw1000.Equals(usd100))
	})

	t.Run("GreaterThan", func(t *testing.T) {
		result, err := krw1000.GreaterThan(krw500)
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = krw500.GreaterThan(krw1000)
		assert.NoError(t, err)
		assert.False(t, result)

		// 다른 통화 비교 시도
		_, err = krw1000.GreaterThan(usd100)
		assert.Error(t, err)
	})

	t.Run("LessThan", func(t *testing.T) {
		result, err := krw500.LessThan(krw1000)
		assert.NoError(t, err)
		assert.True(t, result)

		result, err = krw1000.LessThan(krw500)
		assert.NoError(t, err)
		assert.False(t, result)

		// 다른 통화 비교 시도
		_, err = krw1000.LessThan(usd100)
		assert.Error(t, err)
	})
}

func TestMoney_Checks(t *testing.T) {
	krw1000 := NewTestMoney(1000.0, "KRW")
	krw0 := NewTestMoney(0.0, "KRW")

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, krw0.IsZero())
		assert.False(t, krw1000.IsZero())
	})

	t.Run("IsPositive", func(t *testing.T) {
		assert.True(t, krw1000.IsPositive())
		assert.False(t, krw0.IsPositive())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "1000.00 KRW", krw1000.String())
		assert.Equal(t, "0.00 KRW", krw0.String())
	})
}
