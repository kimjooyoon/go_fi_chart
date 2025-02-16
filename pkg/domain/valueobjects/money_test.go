package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMoney(t *testing.T) {
	tests := []struct {
		name     string
		amount   float64
		currency string
		wantErr  bool
	}{
		{
			name:     "유효한 금액과 통화",
			amount:   100.0,
			currency: "KRW",
			wantErr:  false,
		},
		{
			name:     "0 금액",
			amount:   0.0,
			currency: "USD",
			wantErr:  false,
		},
		{
			name:     "음수 금액",
			amount:   -100.0,
			currency: "EUR",
			wantErr:  true,
		},
		{
			name:     "빈 통화",
			amount:   100.0,
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
	m100KRW, _ := NewMoney(100.0, "KRW")
	m50KRW, _ := NewMoney(50.0, "KRW")
	m100USD, _ := NewMoney(100.0, "USD")

	t.Run("Add", func(t *testing.T) {
		sum, err := m100KRW.Add(m50KRW)
		assert.NoError(t, err)
		assert.Equal(t, 150.0, sum.Amount)
		assert.Equal(t, "KRW", sum.Currency)

		// 다른 통화 더하기 테스트
		_, err = m100KRW.Add(m100USD)
		assert.Error(t, err)
	})

	t.Run("Subtract", func(t *testing.T) {
		diff, err := m100KRW.Subtract(m50KRW)
		assert.NoError(t, err)
		assert.Equal(t, 50.0, diff.Amount)
		assert.Equal(t, "KRW", diff.Currency)

		// 다른 통화 빼기 테스트
		_, err = m100KRW.Subtract(m100USD)
		assert.Error(t, err)
	})

	t.Run("Multiply", func(t *testing.T) {
		result, err := m100KRW.Multiply(2.0)
		assert.NoError(t, err)
		assert.Equal(t, 200.0, result.Amount)
		assert.Equal(t, "KRW", result.Currency)

		// 음수 배수 테스트
		_, err = m100KRW.Multiply(-2.0)
		assert.Error(t, err)
	})

	t.Run("Divide", func(t *testing.T) {
		result, err := m100KRW.Divide(2.0)
		assert.NoError(t, err)
		assert.Equal(t, 50.0, result.Amount)
		assert.Equal(t, "KRW", result.Currency)

		// 0으로 나누기 테스트
		_, err = m100KRW.Divide(0)
		assert.Error(t, err)

		// 음수로 나누기 테스트
		_, err = m100KRW.Divide(-2.0)
		assert.Error(t, err)
	})
}

func TestMoney_Comparisons(t *testing.T) {
	m100KRW, _ := NewMoney(100.0, "KRW")
	m50KRW, _ := NewMoney(50.0, "KRW")
	m100USD, _ := NewMoney(100.0, "USD")
	m0KRW, _ := NewMoney(0.0, "KRW")

	t.Run("IsZero", func(t *testing.T) {
		assert.True(t, m0KRW.IsZero())
		assert.False(t, m100KRW.IsZero())
	})

	t.Run("IsPositive", func(t *testing.T) {
		assert.True(t, m100KRW.IsPositive())
		assert.False(t, m0KRW.IsPositive())
	})

	t.Run("Equals", func(t *testing.T) {
		assert.True(t, m100KRW.Equals(m100KRW))
		assert.False(t, m100KRW.Equals(m50KRW))
		assert.False(t, m100KRW.Equals(m100USD))
	})

	t.Run("GreaterThan", func(t *testing.T) {
		greater, err := m100KRW.GreaterThan(m50KRW)
		assert.NoError(t, err)
		assert.True(t, greater)

		// 다른 통화 비교 테스트
		_, err = m100KRW.GreaterThan(m100USD)
		assert.Error(t, err)
	})

	t.Run("LessThan", func(t *testing.T) {
		less, err := m50KRW.LessThan(m100KRW)
		assert.NoError(t, err)
		assert.True(t, less)

		// 다른 통화 비교 테스트
		_, err = m100KRW.LessThan(m100USD)
		assert.Error(t, err)
	})
}

func TestMoney_String(t *testing.T) {
	m100KRW, _ := NewMoney(100.0, "KRW")
	assert.Equal(t, "100.00 KRW", m100KRW.String())

	m99_99USD, _ := NewMoney(99.99, "USD")
	assert.Equal(t, "99.99 USD", m99_99USD.String())
}
