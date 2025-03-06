package asset

import (
	"context"

	"github.com/aske/go_fi_chart/internal/common/repository"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
)

// Repository는 자산 저장소 인터페이스를 정의합니다.
// 제네릭 기반 Repository 인터페이스를 확장하여 Asset 도메인에 특화된 메서드를 추가합니다.
type Repository interface {
	// 기본 CRUD 작업을 위한 제네릭 Repository 인터페이스 상속
	repository.Repository[*domain.Asset, string]

	// ReadRepository는 읽기 전용 작업을 위한 인터페이스입니다.
	repository.ReadRepository[*domain.Asset, string]

	// WriteRepository는 쓰기 전용 작업을 위한 인터페이스입니다.
	repository.WriteRepository[*domain.Asset, string]

	// FindByUserID는 사용자 ID로 자산 목록을 조회합니다.
	FindByUserID(ctx context.Context, userID string, opts ...repository.FindOption) ([]*domain.Asset, error)

	// FindByType는 자산 유형으로 자산 목록을 조회합니다.
	FindByType(ctx context.Context, assetType domain.AssetType, opts ...repository.FindOption) ([]*domain.Asset, error)

	// CountByUserID는 사용자 ID에 해당하는 자산 개수를 반환합니다.
	CountByUserID(ctx context.Context, userID string) (int64, error)

	// CountByType은 자산 유형에 해당하는 자산 개수를 반환합니다.
	CountByType(ctx context.Context, assetType domain.AssetType) (int64, error)
}
