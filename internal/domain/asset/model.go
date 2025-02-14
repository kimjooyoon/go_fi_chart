package asset

import (
	"time"
)

// Asset 자산을 나타내는 루트 엔티티
type Asset struct {
	ID           string
	UserID       string
	Type         Type
	Name         string
	Amount       float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Transactions []Transaction
}

// Entity 인터페이스 구현
func (a *Asset) GetID() string {
	return a.ID
}

func (a *Asset) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *Asset) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// Type 자산 유형
type Type string

const (
	Cash       Type = "CASH"
	Stock      Type = "STOCK"
	Bond       Type = "BOND"
	RealEstate Type = "REAL_ESTATE"
	Crypto     Type = "CRYPTO"
)

// Transaction 거래 내역
type Transaction struct {
	ID          string
	AssetID     string
	Type        TransactionType
	Amount      float64
	Category    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
}

// Entity 인터페이스 구현
func (t *Transaction) GetID() string {
	return t.ID
}

func (t *Transaction) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Transaction) GetUpdatedAt() time.Time {
	return t.Date // Transaction은 수정되지 않으므로 Date를 UpdatedAt으로 사용
}

// TransactionType 거래 유형
type TransactionType string

const (
	Income   TransactionType = "INCOME"
	Expense  TransactionType = "EXPENSE"
	Transfer TransactionType = "TRANSFER"
)

// Portfolio 포트폴리오 구성
type Portfolio struct {
	ID        string
	UserID    string
	Assets    []PortfolioAsset
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Entity 인터페이스 구현
func (p *Portfolio) GetID() string {
	return p.ID
}

func (p *Portfolio) GetCreatedAt() time.Time {
	return p.CreatedAt
}

func (p *Portfolio) GetUpdatedAt() time.Time {
	return p.UpdatedAt
}

// PortfolioAsset 포트폴리오 내 자산
type PortfolioAsset struct {
	AssetID string
	Weight  float64
}

// NewAsset 새로운 Asset 생성
func NewAsset(userID string, assetType Type, name string, amount float64) *Asset {
	now := time.Now()
	return &Asset{
		ID:        generateID(),
		UserID:    userID,
		Type:      assetType,
		Name:      name,
		Amount:    amount,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewTransaction 새로운 Transaction 생성
func NewTransaction(assetID string, transactionType TransactionType, amount float64, category string, description string) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:          generateID(),
		AssetID:     assetID,
		Type:        transactionType,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        now,
		CreatedAt:   now,
	}
}

// NewPortfolio 새로운 Portfolio 생성
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
