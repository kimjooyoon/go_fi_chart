package valueobjects

import (
	"fmt"
	"math"
)

// RoundingMode는 반올림 정책을 정의합니다.
type RoundingMode int

const (
	// RoundDown 내림
	RoundDown RoundingMode = iota
	// RoundUp 올림
	RoundUp
	// RoundHalfUp 반올림 (0.5 이상이면 올림)
	RoundHalfUp
)

// Money 화폐 값을 나타냅니다.
type Money struct {
	Amount   float64
	Currency string
}

// NewMoney Money 값 객체를 생성합니다.
func NewMoney(amount float64, currency string) (Money, error) {
	if amount < 0 {
		return Money{}, fmt.Errorf("금액은 음수가 될 수 없습니다: %f", amount)
	}
	if currency == "" {
		return Money{}, fmt.Errorf("통화는 비어있을 수 없습니다")
	}
	return Money{
		Amount:   amount,
		Currency: currency,
	}, nil
}

// Round 주어진 정밀도와 반올림 정책으로 금액을 반올림합니다.
func (m Money) Round(precision int, mode RoundingMode) Money {
	multiplier := math.Pow10(precision)
	var rounded float64

	switch mode {
	case RoundDown:
		rounded = math.Floor(m.Amount*multiplier) / multiplier
	case RoundUp:
		rounded = math.Ceil(m.Amount*multiplier) / multiplier
	case RoundHalfUp:
		rounded = math.Round(m.Amount*multiplier) / multiplier
	}

	return Money{
		Amount:   rounded,
		Currency: m.Currency,
	}
}

// Add 두 Money 값을 더합니다.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	return NewMoney(m.Amount+other.Amount, m.Currency)
}

// Subtract 두 Money 값을 뺍니다.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	result := m.Amount - other.Amount
	return NewMoney(result, m.Currency)
}

// Multiply Money 값을 주어진 배수로 곱합니다.
func (m Money) Multiply(multiplier float64) (Money, error) {
	if multiplier < 0 {
		return Money{}, fmt.Errorf("음수 배수는 사용할 수 없습니다: %f", multiplier)
	}
	return NewMoney(m.Amount*multiplier, m.Currency)
}

// Divide Money 값을 주어진 제수로 나눕니다.
func (m Money) Divide(divisor float64) (Money, error) {
	if divisor == 0 {
		return Money{}, fmt.Errorf("0으로 나눌 수 없습니다")
	}
	if divisor < 0 {
		return Money{}, fmt.Errorf("음수로 나눌 수 없습니다: %f", divisor)
	}
	return NewMoney(m.Amount/divisor, m.Currency)
}

// IsZero Money 값이 0인지 확인합니다.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsNegative Money 값이 음수인지 확인합니다.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// IsPositive Money 값이 양수인지 확인합니다.
func (m Money) IsPositive() bool {
	return m.Amount > 0
}

// Equals 두 Money 값이 같은지 확인합니다.
func (m Money) Equals(other Money) bool {
	return m.Amount == other.Amount && m.Currency == other.Currency
}

// GreaterThan 현재 Money 값이 다른 Money 값보다 큰지 확인합니다.
func (m Money) GreaterThan(other Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	return m.Amount > other.Amount, nil
}

// LessThan 현재 Money 값이 다른 Money 값보다 작은지 확인합니다.
func (m Money) LessThan(other Money) (bool, error) {
	if m.Currency != other.Currency {
		return false, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	return m.Amount < other.Amount, nil
}

// String Money 값을 문자열로 변환합니다.
func (m Money) String() string {
	return fmt.Sprintf("%.2f %s", m.Amount, m.Currency)
}
