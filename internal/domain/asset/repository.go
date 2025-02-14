package asset

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// Repository Asset 도메인의 저장소 인터페이스
type Repository interface {
	domain.Repository[*Asset, string]

	// 도메인 특화 메서드
	FindByUserID(ctx context.Context, userID string) ([]*Asset, error)
	FindByType(ctx context.Context, assetType Type) ([]*Asset, error)
	UpdateAmount(ctx context.Context, id string, amount float64) error
}

// TransactionRepository Transaction 도메인의 저장소 인터페이스
type TransactionRepository interface {
	domain.Repository[*Transaction, string]

	// 도메인 특화 메서드
	FindByAssetID(ctx context.Context, assetID string) ([]*Transaction, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*Transaction, error)
	GetTotalAmount(ctx context.Context, assetID string) (float64, error)
}

// PortfolioRepository Portfolio 도메인의 저장소 인터페이스
type PortfolioRepository interface {
	domain.Repository[*Portfolio, string]

	// 도메인 특화 메서드
	FindByUserID(ctx context.Context, userID string) (*Portfolio, error)
	UpdateAssets(ctx context.Context, id string, assets []PortfolioAsset) error
}
