package domain

import (
	"context"

	"github.com/aske/go_fi_chart/internal/common/repository"
)

// AssetRepository는 자산 저장소 인터페이스를 정의합니다.
type AssetRepository interface {
	// 기본 CRUD 작업
	FindByID(ctx context.Context, id string) (*Asset, error)
	FindAll(ctx context.Context, opts ...repository.FindOption) ([]*Asset, error)
	Save(ctx context.Context, entity *Asset) error
	Update(ctx context.Context, entity *Asset) error
	Delete(ctx context.Context, id string) error

	// 추가 조회 기능
	Count(ctx context.Context, opts ...repository.FindOption) (int64, error)

	// 특화된 조회 기능
	FindByUserID(ctx context.Context, userID string, opts ...repository.FindOption) ([]*Asset, error)
	FindByType(ctx context.Context, assetType AssetType, opts ...repository.FindOption) ([]*Asset, error)
	CountByUserID(ctx context.Context, userID string) (int64, error)
	CountByType(ctx context.Context, assetType AssetType) (int64, error)
}
