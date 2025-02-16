package valueobjects

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
