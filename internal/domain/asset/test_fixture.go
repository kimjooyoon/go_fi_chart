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
	now := time.Now()
	userID := "test-user-1"

	// 자산 데이터 생성
	assets := []*Asset{
		{
			ID:        "asset-1",
			UserID:    userID,
			Type:      Cash,
			Name:      "현금 자산",
			Amount:    1000000,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "asset-2",
			UserID:    userID,
			Type:      Stock,
			Name:      "주식 투자",
			Amount:    5000000,
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "asset-3",
			UserID:    userID,
			Type:      Bond,
			Name:      "채권 투자",
			Amount:    3000000,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	// 거래 내역 데이터 생성
	transactions := []*Transaction{
		{
			ID:          "tx-1",
			AssetID:     "asset-1",
			Type:        Income,
			Amount:      500000,
			Category:    "급여",
			Description: "3월 급여",
			Date:        now,
			CreatedAt:   now,
		},
		{
			ID:          "tx-2",
			AssetID:     "asset-1",
			Type:        Expense,
			Amount:      100000,
			Category:    "식비",
			Description: "3월 식비",
			Date:        now,
			CreatedAt:   now,
		},
		{
			ID:          "tx-3",
			AssetID:     "asset-2",
			Type:        Transfer,
			Amount:      1000000,
			Category:    "투자",
			Description: "주식 매수",
			Date:        now,
			CreatedAt:   now,
		},
	}

	// 포트폴리오 데이터 생성
	portfolios := []*Portfolio{
		{
			ID:     "portfolio-1",
			UserID: userID,
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
			CreatedAt: now,
			UpdatedAt: now,
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
