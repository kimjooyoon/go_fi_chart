package asset

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// Repository 자산 저장소 인터페이스입니다.
type Repository interface {
	Save(ctx context.Context, asset *Asset) error
	FindByID(ctx context.Context, id string) (*Asset, error)
	Update(ctx context.Context, asset *Asset) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) ([]*Asset, error)
	FindByType(ctx context.Context, assetType Type) ([]*Asset, error)
	UpdateAmount(ctx context.Context, id string, amount Money) error
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Asset, error)
	FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Asset, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// TransactionRepository 거래 내역 저장소 인터페이스입니다.
type TransactionRepository interface {
	Save(ctx context.Context, tx *Transaction) error
	FindByID(ctx context.Context, id string) (*Transaction, error)
	Update(ctx context.Context, tx *Transaction) error
	Delete(ctx context.Context, id string) error
	FindByAssetID(ctx context.Context, assetID string) ([]*Transaction, error)
	FindByDateRange(ctx context.Context, start, end time.Time) ([]*Transaction, error)
	GetTotalAmount(ctx context.Context, assetID string) (Money, error)
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Transaction, error)
	FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Transaction, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// PortfolioRepository 포트폴리오 저장소 인터페이스입니다.
type PortfolioRepository interface {
	Save(ctx context.Context, portfolio *Portfolio) error
	FindByID(ctx context.Context, id string) (*Portfolio, error)
	Update(ctx context.Context, portfolio *Portfolio) error
	Delete(ctx context.Context, id string) error
	FindByUserID(ctx context.Context, userID string) (*Portfolio, error)
	UpdateAssets(ctx context.Context, id string, assets []PortfolioAsset) error
	FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Portfolio, error)
	FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Portfolio, error)
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
