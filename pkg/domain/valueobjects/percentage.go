package valueobjects

import (
	"fmt"
)

// Percentage 퍼센트 값을 나타냅니다.
type Percentage struct {
	Value float64
}

// NewPercentage Percentage 값 객체를 생성합니다.
func NewPercentage(value float64) (Percentage, error) {
	if value < 0 || value > 100 {
		return Percentage{}, fmt.Errorf("퍼센트 값은 0에서 100 사이여야 합니다: %f", value)
	}
	return Percentage{Value: value}, nil
}

// Add 두 Percentage 값을 더합니다.
func (p Percentage) Add(other Percentage) (Percentage, error) {
	sum := p.Value + other.Value
	return NewPercentage(sum)
}

// Subtract 두 Percentage 값을 뺍니다.
func (p Percentage) Subtract(other Percentage) (Percentage, error) {
	diff := p.Value - other.Value
	return NewPercentage(diff)
}

// Multiply Percentage 값을 주어진 배수로 곱합니다.
func (p Percentage) Multiply(multiplier float64) (Percentage, error) {
	result := p.Value * multiplier
	return NewPercentage(result)
}

// IsZero Percentage 값이 0인지 확인합니다.
func (p Percentage) IsZero() bool {
	return p.Value == 0
}

// IsComplete Percentage 값이 100%인지 확인합니다.
func (p Percentage) IsComplete() bool {
	return p.Value == 100
}

// ToDecimal Percentage 값을 소수로 변환합니다.
func (p Percentage) ToDecimal() float64 {
	return p.Value / 100
}

// FromDecimal 소수를 Percentage로 변환합니다.
func FromDecimal(decimal float64) (Percentage, error) {
	return NewPercentage(decimal * 100)
}
