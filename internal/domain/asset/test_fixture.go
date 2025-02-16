package asset

import (
	"time"
)

// TestFixture 테스트에 사용할 데이터 세트
type TestFixture struct {
	Assets       []*Asset
	Transactions []*Transaction
	Portfolios   []*Portfolio
}

// NewTestFixture 새로운 테스트 픽스처를 생성합니다.
func NewTestFixture() *TestFixture {
	assets := []*Asset{
		{
			ID:        "asset-1",
			UserID:    "test-user-1",
			Type:      Cash,
			Name:      "현금 자산",
			Amount:    Money{Amount: 1000000.0, Currency: "KRW"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "asset-2",
			UserID:    "test-user-1",
			Type:      Stock,
			Name:      "주식 자산",
			Amount:    Money{Amount: 2000000.0, Currency: "KRW"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "asset-3",
			UserID:    "test-user-1",
			Type:      RealEstate,
			Name:      "부동산 자산",
			Amount:    Money{Amount: 300000000.0, Currency: "KRW"},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	transactions := []*Transaction{
		{
			ID:          "tx-1",
			AssetID:     assets[0].ID,
			Type:        Income,
			Amount:      Money{Amount: 500000.0, Currency: "KRW"},
			Category:    "급여",
			Description: "3월 급여",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		},
		{
			ID:          "tx-2",
			AssetID:     assets[0].ID,
			Type:        Expense,
			Amount:      Money{Amount: 100000.0, Currency: "KRW"},
			Category:    "식비",
			Description: "3월 식비",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		},
		{
			ID:          "tx-3",
			AssetID:     assets[0].ID,
			Type:        Transfer,
			Amount:      Money{Amount: 1000000.0, Currency: "KRW"},
			Category:    "이체",
			Description: "주식 계좌로 이체",
			Date:        time.Now(),
			CreatedAt:   time.Now(),
		},
	}

	// 포트폴리오 데이터 생성
	portfolios := []*Portfolio{
		{
			ID:     "portfolio-1",
			UserID: "test-user-1",
			Assets: []PortfolioAsset{
				{
					AssetID: "asset-1",
					Weight:  0.2, // 20%
				},
				{
					AssetID: "asset-2",
					Weight:  0.5, // 50%
				},
				{
					AssetID: "asset-3",
					Weight:  0.3, // 30%
				},
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	return &TestFixture{
		Assets:       assets,
		Transactions: transactions,
		Portfolios:   portfolios,
	}
}

// GetAssetByID 테스트 픽스처에서 ID로 Asset을 찾습니다.
func (f *TestFixture) GetAssetByID(id string) *Asset {
	for _, asset := range f.Assets {
		if asset.ID == id {
			return asset
		}
	}
	return nil
}

// GetTransactionByID 테스트 픽스처에서 ID로 Transaction을 찾습니다.
func (f *TestFixture) GetTransactionByID(id string) *Transaction {
	for _, tx := range f.Transactions {
		if tx.ID == id {
			return tx
		}
	}
	return nil
}

// GetPortfolioByID 테스트 픽스처에서 ID로 Portfolio를 찾습니다.
func (f *TestFixture) GetPortfolioByID(id string) *Portfolio {
	for _, portfolio := range f.Portfolios {
		if portfolio.ID == id {
			return portfolio
		}
	}
	return nil
}

// GetAssetsByUserID 테스트 픽스처에서 UserID로 Asset 목록을 찾습니다.
func (f *TestFixture) GetAssetsByUserID(userID string) []*Asset {
	var result []*Asset
	for _, asset := range f.Assets {
		if asset.UserID == userID {
			result = append(result, asset)
		}
	}
	return result
}

// GetTransactionsByAssetID 테스트 픽스처에서 AssetID로 Transaction 목록을 찾습니다.
func (f *TestFixture) GetTransactionsByAssetID(assetID string) []*Transaction {
	var result []*Transaction
	for _, tx := range f.Transactions {
		if tx.AssetID == assetID {
			result = append(result, tx)
		}
	}
	return result
}

// GetPortfolioByUserID 테스트 픽스처에서 UserID로 Portfolio를 찾습니다.
func (f *TestFixture) GetPortfolioByUserID(userID string) *Portfolio {
	for _, portfolio := range f.Portfolios {
		if portfolio.UserID == userID {
			return portfolio
		}
	}
	return nil
}

// CreateTestAsset 테스트용 자산을 생성합니다.
func CreateTestAsset() *Asset {
	return &Asset{
		ID:        "test-asset-1",
		UserID:    "test-user-1",
		Type:      Stock,
		Name:      "삼성전자",
		Amount:    NewMoney(500000, "KRW"),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestTransaction 테스트용 거래 내역을 생성합니다.
func CreateTestTransaction() *Transaction {
	tx, _ := NewTransaction(
		"test-asset-1",
		Income,
		NewMoney(100000, "KRW"),
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
