package asset

import (
	"time"
)

// TestFixture 테스트에 사용할 데이터 세트
type TestFixture struct {
	assets       map[string]*Asset
	transactions map[string]*Transaction
	portfolios   map[string]*Portfolio
}

// NewTestFixture 새로운 테스트 픽스처를 생성합니다.
func NewTestFixture() *TestFixture {
	return &TestFixture{
		assets:       make(map[string]*Asset),
		transactions: make(map[string]*Transaction),
		portfolios:   make(map[string]*Portfolio),
	}
}

// GetAssetByID 테스트 픽스처에서 ID로 Asset을 찾습니다.
func (f *TestFixture) GetAssetByID(id string) *Asset {
	return f.assets[id]
}

// GetTransactionByID 테스트 픽스처에서 ID로 Transaction을 찾습니다.
func (f *TestFixture) GetTransactionByID(id string) *Transaction {
	return f.transactions[id]
}

// GetPortfolioByID 테스트 픽스처에서 ID로 Portfolio를 찾습니다.
func (f *TestFixture) GetPortfolioByID(id string) *Portfolio {
	return f.portfolios[id]
}

// GetAssetsByUserID 테스트 픽스처에서 UserID로 Asset 목록을 찾습니다.
func (f *TestFixture) GetAssetsByUserID(userID string) []*Asset {
	var result []*Asset
	for _, asset := range f.assets {
		if asset.UserID == userID {
			result = append(result, asset)
		}
	}
	return result
}

// GetTransactionsByAssetID 테스트 픽스처에서 AssetID로 Transaction 목록을 찾습니다.
func (f *TestFixture) GetTransactionsByAssetID(assetID string) []*Transaction {
	var result []*Transaction
	for _, tx := range f.transactions {
		if tx.AssetID == assetID {
			result = append(result, tx)
		}
	}
	return result
}

// GetPortfolioByUserID 테스트 픽스처에서 UserID로 Portfolio를 찾습니다.
func (f *TestFixture) GetPortfolioByUserID(userID string) *Portfolio {
	for _, portfolio := range f.portfolios {
		if portfolio.UserID == userID {
			return portfolio
		}
	}
	return nil
}

// CreateTestAsset 테스트용 자산을 생성합니다.
func CreateTestAsset() *Asset {
	money := NewTestMoney(500000, "KRW")
	asset := &Asset{
		ID:        "test-asset-1",
		UserID:    "test-user-1",
		Type:      Stock,
		Name:      "삼성전자",
		Amount:    money,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return asset
}

// CreateTestTransaction 테스트용 거래 내역을 생성합니다.
func CreateTestTransaction() *Transaction {
	money := NewTestMoney(100000, "KRW")
	tx, _ := NewTransaction(
		"test-asset-1",
		Income,
		money,
		"급여",
		"2월 급여",
	)
	return tx
}

// CreateTestPortfolio 테스트용 포트폴리오를 생성합니다.
func CreateTestPortfolio() *Portfolio {
	return &Portfolio{
		ID:        "test-portfolio-1",
		UserID:    "test-user-1",
		Assets:    []PortfolioAsset{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewTestMoney 테스트용 Money 값을 생성합니다.
func NewTestMoney(amount float64, currency string) Money {
	money, err := NewMoney(amount, currency)
	if err != nil {
		panic(err)
	}
	return money
}

// NewTestAsset 테스트용 자산을 생성합니다.
func NewTestAsset() *Asset {
	asset, err := NewAsset("test-user", Cash, "Test Asset", 1000000, "KRW")
	if err != nil {
		panic(err)
	}
	return asset
}

// NewTestTransaction 테스트용 거래를 생성합니다.
func NewTestTransaction() *Transaction {
	money := NewTestMoney(500000, "KRW")
	tx, err := NewTransaction("test-asset", Income, money, "Test", "Test Transaction")
	if err != nil {
		panic(err)
	}
	return tx
}

// NewTestPortfolio 테스트용 포트폴리오를 생성합니다.
func NewTestPortfolio() *Portfolio {
	return NewPortfolio("test-user", []PortfolioAsset{
		{
			AssetID: "test-asset-1",
			Weight:  0.6,
		},
		{
			AssetID: "test-asset-2",
			Weight:  0.4,
		},
	})
}

// CreateFixture 테스트 데이터를 생성합니다.
func CreateFixture() *TestFixture {
	fixture := NewTestFixture()

	// 자산 생성
	asset1, err := NewAsset("user-1", Cash, "현금 자산", 1000000, "KRW")
	if err != nil {
		panic(err)
	}
	asset2, err := NewAsset("user-1", Stock, "주식 자산", 2000000, "KRW")
	if err != nil {
		panic(err)
	}
	fixture.assets[asset1.ID] = asset1
	fixture.assets[asset2.ID] = asset2

	// 거래 생성
	money1 := NewTestMoney(500000, "KRW")
	tx1, err := NewTransaction(asset1.ID, Income, money1, "급여", "2월 급여")
	if err != nil {
		panic(err)
	}
	money2 := NewTestMoney(300000, "KRW")
	tx2, err := NewTransaction(asset1.ID, Expense, money2, "식비", "2월 식비")
	if err != nil {
		panic(err)
	}
	money3 := NewTestMoney(200000, "KRW")
	tx3, err := NewTransaction(asset2.ID, Income, money3, "배당금", "2월 배당금")
	if err != nil {
		panic(err)
	}
	fixture.transactions[tx1.ID] = tx1
	fixture.transactions[tx2.ID] = tx2
	fixture.transactions[tx3.ID] = tx3

	// 포트폴리오 생성
	portfolio := NewPortfolio("user-1", []PortfolioAsset{
		{AssetID: asset1.ID, Weight: 0.6},
		{AssetID: asset2.ID, Weight: 0.4},
	})
	fixture.portfolios[portfolio.ID] = portfolio

	return fixture
}
