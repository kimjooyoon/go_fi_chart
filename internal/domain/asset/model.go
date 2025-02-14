package asset

import (
	"time"
)

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

func (a *Asset) GetID() string {
	return a.ID
}

func (a *Asset) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *Asset) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

type Type string

const (
	Cash       Type = "CASH"
	Stock      Type = "STOCK"
	Bond       Type = "BOND"
	RealEstate Type = "REAL_ESTATE"
	Crypto     Type = "CRYPTO"
)

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

func (t *Transaction) GetID() string {
	return t.ID
}

func (t *Transaction) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t *Transaction) GetUpdatedAt() time.Time {
	return t.Date
}

type TransactionType string

const (
	Income   TransactionType = "INCOME"
	Expense  TransactionType = "EXPENSE"
	Transfer TransactionType = "TRANSFER"
)

type Portfolio struct {
	ID        string
	UserID    string
	Assets    []PortfolioAsset
	CreatedAt time.Time
	UpdatedAt time.Time
}

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
