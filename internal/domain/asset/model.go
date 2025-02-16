package asset

import (
	"fmt"
	"time"

	"github.com/aske/go_fi_chart/internal/domain/event"
	"github.com/google/uuid"
)

// Money 화폐 값을 나타냅니다.
type Money struct {
	Amount   float64
	Currency string
}

// NewMoney Money 값 객체를 생성합니다.
func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

// Add 두 Money 값을 더합니다.
func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}, nil
}

// Subtract 두 Money 값을 뺍니다.
func (m Money) Subtract(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, fmt.Errorf("통화가 일치하지 않습니다: %s != %s", m.Currency, other.Currency)
	}
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}, nil
}

// Multiply Money 값을 주어진 배수로 곱합니다.
func (m Money) Multiply(multiplier float64) Money {
	return Money{
		Amount:   m.Amount * multiplier,
		Currency: m.Currency,
	}
}

// IsZero Money 값이 0인지 확인합니다.
func (m Money) IsZero() bool {
	return m.Amount == 0
}

// IsNegative Money 값이 음수인지 확인합니다.
func (m Money) IsNegative() bool {
	return m.Amount < 0
}

// Performance 자산의 성과를 나타냅니다.
type Performance struct {
	StartValue     Money
	CurrentValue   Money
	GrowthRate     float64
	RiskScore      float64
	LastUpdateTime time.Time
}

// Goal 재무 목표를 나타냅니다.
type Goal struct {
	ID        string
	Type      GoalType
	Target    Money
	Deadline  time.Time
	Progress  float64
	Rewards   []Reward
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GoalType 목표의 유형을 나타냅니다.
type GoalType string

const (
	GoalTypeSaving    GoalType = "SAVING"
	GoalTypeInvesting GoalType = "INVESTING"
	GoalTypeDebtFree  GoalType = "DEBT_FREE"
)

// Achievement 업적을 나타냅니다.
type Achievement struct {
	ID         string
	Type       AchievementType
	Progress   float64
	Conditions []Condition
	Rewards    []Reward
	UnlockedAt *time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// AchievementType 업적의 유형을 나타냅니다.
type AchievementType string

const (
	AchievementTypeSaving    AchievementType = "SAVING_MASTER"
	AchievementTypeInvesting AchievementType = "INVESTING_GURU"
	AchievementTypeCommunity AchievementType = "COMMUNITY_STAR"
)

// Condition 업적 달성 조건을 나타냅니다.
type Condition struct {
	Type      ConditionType
	Target    interface{}
	Current   interface{}
	Completed bool
}

// ConditionType 조건의 유형을 나타냅니다.
type ConditionType string

const (
	ConditionTypeAmount      ConditionType = "AMOUNT"
	ConditionTypeDuration    ConditionType = "DURATION"
	ConditionTypeStreak      ConditionType = "STREAK"
	ConditionTypeInteraction ConditionType = "INTERACTION"
)

// Reward 보상을 나타냅니다.
type Reward struct {
	Type    RewardType
	Value   interface{}
	Claimed bool
	ClaimBy time.Time
}

// RewardType 보상의 유형을 나타냅니다.
type RewardType string

const (
	RewardTypeBadge   RewardType = "BADGE"
	RewardTypeTitle   RewardType = "TITLE"
	RewardTypeFeature RewardType = "FEATURE"
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

// Asset 자산을 나타냅니다.
type Asset struct {
	ID           string
	UserID       string
	Type         Type
	Name         string
	Amount       Money
	Performance  Performance
	Goals        []Goal
	Achievements []Achievement
	CreatedAt    time.Time
	UpdatedAt    time.Time
	events       []event.Event // 미발행 이벤트 저장
}

func (a *Asset) GetID() string {
	return a.ID
}

func (a *Asset) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *Asset) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// Type 자산의 유형을 나타냅니다.
type Type string

const (
	Cash       Type = "CASH"
	Stock      Type = "STOCK"
	Bond       Type = "BOND"
	RealEstate Type = "REAL_ESTATE"
	Crypto     Type = "CRYPTO"
)

// NewAsset 새로운 자산을 생성합니다.
func NewAsset(userID string, assetType Type, name string, amount float64, currency string) *Asset {
	now := time.Now()
	return &Asset{
		ID:     generateID(),
		UserID: userID,
		Type:   assetType,
		Name:   name,
		Amount: Money{
			Amount:   amount,
			Currency: currency,
		},
		Performance: Performance{
			StartValue: Money{
				Amount:   amount,
				Currency: currency,
			},
			CurrentValue: Money{
				Amount:   amount,
				Currency: currency,
			},
			LastUpdateTime: now,
		},
		Goals:        make([]Goal, 0),
		Achievements: make([]Achievement, 0),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddGoal 자산에 새로운 목표를 추가합니다.
func (a *Asset) AddGoal(goalType GoalType, target Money, deadline time.Time) *Goal {
	now := time.Now()
	goal := &Goal{
		ID:        generateID(),
		Type:      goalType,
		Target:    target,
		Deadline:  deadline,
		Progress:  0,
		Rewards:   make([]Reward, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
	a.Goals = append(a.Goals, *goal)
	return goal
}

// UpdateProgress 목표의 진행률을 업데이트합니다.
func (g *Goal) UpdateProgress(current Money) {
	g.Progress = (current.Amount / g.Target.Amount) * 100
	g.UpdatedAt = time.Now()
}

// IsAchieved 목표가 달성되었는지 확인합니다.
func (g *Goal) IsAchieved() bool {
	return g.Progress >= 100
}

// AddAchievement 자산에 새로운 업적을 추가합니다.
func (a *Asset) AddAchievement(achievementType AchievementType, conditions []Condition) *Achievement {
	now := time.Now()
	achievement := &Achievement{
		ID:         generateID(),
		Type:       achievementType,
		Progress:   0,
		Conditions: conditions,
		Rewards:    make([]Reward, 0),
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	a.Achievements = append(a.Achievements, *achievement)
	return achievement
}

// UpdateAchievementProgress 업적의 진행률을 업데이트합니다.
func (a *Achievement) UpdateProgress() {
	var completed int
	for _, condition := range a.Conditions {
		if condition.Completed {
			completed++
		}
	}
	a.Progress = float64(completed) / float64(len(a.Conditions)) * 100

	if a.Progress >= 100 && a.UnlockedAt == nil {
		now := time.Now()
		a.UnlockedAt = &now
	}

	a.UpdatedAt = time.Now()
}

// IsUnlocked 업적이 해금되었는지 확인합니다.
func (a *Achievement) IsUnlocked() bool {
	return a.UnlockedAt != nil
}

// TransactionType 거래의 유형을 나타냅니다.
type TransactionType string

const (
	Income   TransactionType = "INCOME"
	Expense  TransactionType = "EXPENSE"
	Transfer TransactionType = "TRANSFER"
)

// Transaction 거래 내역을 나타냅니다.
type Transaction struct {
	ID          string
	AssetID     string
	Type        TransactionType
	Amount      Money
	Category    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

// NewTransaction 새로운 Transaction 값 객체를 생성합니다.
func NewTransaction(assetID string, transactionType TransactionType, amount Money, category string, description string) (*Transaction, error) {
	if amount.IsZero() {
		return nil, fmt.Errorf("거래 금액은 0이 될 수 없습니다")
	}
	if amount.IsNegative() {
		return nil, fmt.Errorf("거래 금액은 음수가 될 수 없습니다")
	}

	now := time.Now()
	return &Transaction{
		ID:          uuid.New().String(),
		AssetID:     assetID,
		Type:        transactionType,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        now,
		CreatedAt:   now,
	}, nil
}

// GetID 거래의 ID를 반환합니다.
func (t *Transaction) GetID() string {
	return t.ID
}

// GetAmount 거래 금액을 반환합니다.
func (t *Transaction) GetAmount() Money {
	return t.Amount
}

// GetDate 거래 일자를 반환합니다.
func (t *Transaction) GetDate() time.Time {
	return t.Date
}

// GetCreatedAt 거래 생성 일자를 반환합니다.
func (t *Transaction) GetCreatedAt() time.Time {
	return t.CreatedAt
}

// GetUpdatedAt 거래의 업데이트 일자를 반환합니다.
func (t *Transaction) GetUpdatedAt() time.Time {
	return t.Date
}

// Portfolio 포트폴리오를 나타냅니다.
type Portfolio struct {
	ID        string
	UserID    string
	Assets    []PortfolioAsset
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GetID 포트폴리오의 ID를 반환합니다.
func (p *Portfolio) GetID() string {
	return p.ID
}

// GetCreatedAt 포트폴리오의 생성 일자를 반환합니다.
func (p *Portfolio) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetUpdatedAt 포트폴리오의 업데이트 일자를 반환합니다.
func (p *Portfolio) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// PortfolioAsset 포트폴리오의 자산 구성을 나타냅니다.
type PortfolioAsset struct {
	AssetID string
	Weight  float64
}

func NewPortfolio(userID string, assets []PortfolioAsset) *Portfolio {
	now := time.Now()
	return &Portfolio{
		ID:        generateID(),
		UserID:    userID,
		Assets:    assets,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddEvent 이벤트를 추가합니다.
func (a *Asset) AddEvent(evt event.Event) {
	if a.events == nil {
		a.events = make([]event.Event, 0)
	}
	a.events = append(a.events, evt)
}

// GetUncommittedEvents 미발행 이벤트를 반환합니다.
func (a *Asset) GetUncommittedEvents() []event.Event {
	return a.events
}

// ClearEvents 이벤트를 초기화합니다.
func (a *Asset) ClearEvents() {
	a.events = make([]event.Event, 0)
}

// ProcessTransaction 거래를 처리합니다.
func (a *Asset) ProcessTransaction(tx *Transaction) error {
	if err := a.ValidateTransaction(tx); err != nil {
		return err
	}

	switch tx.Type {
	case Income:
		result, err := a.Amount.Add(tx.Amount)
		if err != nil {
			return err
		}
		a.Amount = result
	case Expense:
		result, err := a.Amount.Subtract(tx.Amount)
		if err != nil {
			return err
		}
		a.Amount = result
	case Transfer:
		result, err := a.Amount.Subtract(tx.Amount)
		if err != nil {
			return err
		}
		a.Amount = result
	}

	a.UpdatedAt = time.Now()

	// 거래 처리 이벤트 발행
	a.AddEvent(event.NewEvent(
		event.TypeTransactionRecorded,
		a.ID,
		"asset",
		map[string]interface{}{
			"transactionID": tx.ID,
			"type":          tx.Type,
			"amount":        tx.Amount,
		},
		map[string]string{
			"userID": a.UserID,
		},
		1,
	))

	return nil
}

// ValidateTransaction 거래가 유효한지 검증합니다.
func (a *Asset) ValidateTransaction(tx *Transaction) error {
	if tx.Amount.IsZero() {
		return fmt.Errorf("거래 금액은 0이 될 수 없습니다")
	}

	if tx.Amount.IsNegative() {
		return fmt.Errorf("거래 금액은 음수가 될 수 없습니다")
	}

	if tx.Type == Expense && a.Amount.Amount < tx.Amount.Amount {
		return fmt.Errorf("잔액이 부족합니다")
	}

	return nil
}

// TimeRange 시간 범위를 나타냅니다.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewTimeRange TimeRange 값 객체를 생성합니다.
func NewTimeRange(start, end time.Time) (TimeRange, error) {
	if end.Before(start) {
		return TimeRange{}, fmt.Errorf("종료 시간은 시작 시간보다 이후여야 합니다: %v > %v", start, end)
	}
	return TimeRange{
		Start: start,
		End:   end,
	}, nil
}

// Duration 기간을 반환합니다.
func (tr TimeRange) Duration() time.Duration {
	return tr.End.Sub(tr.Start)
}

// Contains 주어진 시간이 범위 내에 있는지 확인합니다.
func (tr TimeRange) Contains(t time.Time) bool {
	return (t.Equal(tr.Start) || t.After(tr.Start)) && (t.Equal(tr.End) || t.Before(tr.End))
}

// Overlaps 다른 TimeRange와 겹치는지 확인합니다.
func (tr TimeRange) Overlaps(other TimeRange) bool {
	return tr.Contains(other.Start) || tr.Contains(other.End) ||
		other.Contains(tr.Start) || other.Contains(tr.End)
}

// IsZero TimeRange가 zero value인지 확인합니다.
func (tr TimeRange) IsZero() bool {
	return tr.Start.IsZero() && tr.End.IsZero()
}

// Split TimeRange를 주어진 간격으로 분할합니다.
func (tr TimeRange) Split(interval time.Duration) []TimeRange {
	if interval <= 0 {
		return []TimeRange{tr}
	}

	var ranges []TimeRange
	current := tr.Start
	for current.Before(tr.End) {
		next := current.Add(interval)
		if next.After(tr.End) {
			next = tr.End
		}
		if r, err := NewTimeRange(current, next); err == nil {
			ranges = append(ranges, r)
		}
		current = next
	}
	return ranges
}

// Extend TimeRange를 주어진 기간만큼 확장합니다.
func (tr TimeRange) Extend(d time.Duration) (TimeRange, error) {
	return NewTimeRange(tr.Start, tr.End.Add(d))
}

// Shift TimeRange를 주어진 기간만큼 이동합니다.
func (tr TimeRange) Shift(d time.Duration) TimeRange {
	return TimeRange{
		Start: tr.Start.Add(d),
		End:   tr.End.Add(d),
	}
}
