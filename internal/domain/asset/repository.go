package asset

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// Repository Asset 도메인의 저장소 인터페이스
type Repository interface {
	domain.Repository[*Asset, string]

	FindByUserID(ctx context.Context, userID string) ([]*Asset, error)
	FindByType(ctx context.Context, assetType Type) ([]*Asset, error)
	UpdateAmount(ctx context.Context, id string, amount float64) error
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Asset, error)
}

// TransactionRepository Transaction 도메인의 저장소 인터페이스
type TransactionRepository interface {
	domain.Repository[*Transaction, string]

	FindByAssetID(ctx context.Context, assetID string) ([]*Transaction, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*Transaction, error)
	GetTotalAmount(ctx context.Context, assetID string) (Money, error)
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Transaction, error)
}

// PortfolioRepository Portfolio 도메인의 저장소 인터페이스
type PortfolioRepository interface {
	domain.Repository[*Portfolio, string]

	FindByUserID(ctx context.Context, userID string) (*Portfolio, error)
	UpdateAssets(ctx context.Context, id string, assets []PortfolioAsset) error
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Portfolio, error)
}
