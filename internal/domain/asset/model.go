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

// AddEvent 새로운 이벤트를 추가합니다.
func (a *Asset) AddEvent(event event.Event) {
	a.events = append(a.events, event)
}

// GetUncommittedEvents 미발행 이벤트 목록을 반환합니다.
func (a *Asset) GetUncommittedEvents() []event.Event {
	return a.events
}

// ClearEvents 미발행 이벤트를 모두 제거합니다.
func (a *Asset) ClearEvents() {
	a.events = make([]event.Event, 0)
}

// ProcessTransaction 거래를 처리하고 이벤트를 발행합니다.
func (a *Asset) ProcessTransaction(tx *Transaction) error {
	if err := a.ValidateTransaction(tx); err != nil {
		return err
	}

	switch tx.Type {
	case Income, Transfer:
		a.Amount, _ = a.Amount.Add(tx.Amount)
	case Expense:
		a.Amount, _ = a.Amount.Subtract(tx.Amount)
	}

	a.UpdatedAt = time.Now()

	// 거래 처리 이벤트 발행
	a.AddEvent(event.NewEvent(
		"asset.transaction.processed",
		"asset",
		map[string]interface{}{
			"assetId":       a.ID,
			"transactionId": tx.ID,
			"type":          tx.Type,
			"amount":        tx.Amount.Amount,
			"newBalance":    a.Amount.Amount,
		},
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
