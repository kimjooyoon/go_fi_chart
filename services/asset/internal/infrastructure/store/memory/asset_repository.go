package memory

import (
	"context"
	"sync"
	
	"github.com/aske/go_fi_chart/internal/common/repository"
	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	assetRepo "github.com/aske/go_fi_chart/services/asset/internal/domain/asset"
)

// AssetRepository는 인메모리 자산 저장소 구현체입니다.
type AssetRepository struct {
	assets   map[string]*domain.Asset
	mu       sync.RWMutex
	eventBus events.EventBus
}

// NewAssetRepository는 새로운 인메모리 자산 저장소를 생성합니다.
func NewAssetRepository(eventBus events.EventBus) *AssetRepository {
	return &AssetRepository{
		assets:   make(map[string]*domain.Asset),
		eventBus: eventBus,
	}
}

// FindByID는 ID로 자산을 조회합니다.
func (r *AssetRepository) FindByID(_ context.Context, id string) (*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	asset, exists := r.assets[id]
	if !exists {
		return nil, repository.NewError(
			"FindByID",
			"Asset",
			"asset not found",
			repository.ErrEntityNotFound,
		)
	}
	
	if asset.IsDeleted {
		return nil, repository.NewError(
			"FindByID",
			"Asset",
			"asset is deleted",
			repository.ErrEntityNotFound,
		)
	}
	
	return asset, nil
}

// FindAll은 모든 자산을 조회합니다. 옵션을 통해 필터링, 정렬, 페이지네이션을 적용할 수 있습니다.
func (r *AssetRepository) FindAll(_ context.Context, opts ...repository.FindOption) ([]*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// 옵션 적용
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}
	
	// 결과 배열 생성
	var result []*domain.Asset
	for _, asset := range r.assets {
		if asset.IsDeleted {
			continue
		}
		
		// 필터링 적용
		if r.matchesFilters(asset, options.Filters) {
			result = append(result, asset)
		}
	}
	
	// 정렬 적용 (미구현 상태)
	
	// 페이지네이션 적용
	if options.Limit > 0 {
		offset := int64(options.Offset)
		limit := int64(options.Limit)
		
		if offset >= int64(len(result)) {
			return []*domain.Asset{}, nil
		}
		
		end := offset + limit
		if end > int64(len(result)) {
			end = int64(len(result))
		}
		
		return result[offset:end], nil
	}
	
	return result, nil
}

// 필터링 조건과 자산이 일치하는지 확인하는 헬퍼 메서드
func (r *AssetRepository) matchesFilters(asset *domain.Asset, filters map[string]interface{}) bool {
	for field, value := range filters {
		switch field {
		case "user_id":
			if asset.UserID != value.(string) {
				return false
			}
		case "type":
			if asset.Type != value.(domain.AssetType) {
				return false
			}
		}
	}
	return true
}

// Save는 자산을 저장합니다.
func (r *AssetRepository) Save(ctx context.Context, asset *domain.Asset) error {
	events := asset.Events()
	
	r.mu.Lock()
	if _, exists := r.assets[asset.ID]; exists {
		r.mu.Unlock()
		return repository.NewError(
			"Save",
			"Asset",
			"asset already exists",
			repository.ErrDuplicateEntity,
		)
	}
	
	r.assets[asset.ID] = asset
	r.mu.Unlock()
	
	// 락 해제 후 이벤트 발행
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Save",
				"Asset",
				"failed to publish event",
				repository.ErrRepositoryError,
			)
		}
	}
	asset.ClearEvents()
	
	return nil
}

// Update는 자산을 업데이트합니다.
func (r *AssetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	events := asset.Events()
	
	r.mu.Lock()
	if _, exists := r.assets[asset.ID]; !exists {
		r.mu.Unlock()
		return repository.NewError(
			"Update",
			"Asset",
			"asset not found",
			repository.ErrEntityNotFound,
		)
	}
	
	if r.assets[asset.ID].IsDeleted {
		r.mu.Unlock()
		return repository.NewError(
			"Update",
			"Asset",
			"asset is deleted",
			repository.ErrEntityNotFound,
		)
	}
	
	r.assets[asset.ID] = asset
	r.mu.Unlock()
	
	// 락 해제 후 이벤트 발행
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Update",
				"Asset",
				"failed to publish event",
				repository.ErrRepositoryError,
			)
		}
	}
	asset.ClearEvents()
	
	return nil
}

// Delete는 자산을 삭제합니다.
func (r *AssetRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	
	asset, exists := r.assets[id]
	if !exists {
		r.mu.Unlock()
		return repository.NewError(
			"Delete",
			"Asset",
			"asset not found",
			repository.ErrEntityNotFound,
		)
	}
	
	if asset.IsDeleted {
		r.mu.Unlock()
		return repository.NewError(
			"Delete",
			"Asset",
			"asset is already deleted",
			repository.ErrEntityNotFound,
		)
	}
	
	asset.MarkAsDeleted()
	r.mu.Unlock()
	
	// 이벤트 발행
	for _, event := range asset.Events() {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return repository.NewError(
				"Delete",
				"Asset",
				"failed to publish event",
				repository.ErrRepositoryError,
			)
		}
	}
	
	return nil
}

// Count는 조건에 맞는 자산의 총 개수를 반환합니다.
func (r *AssetRepository) Count(_ context.Context, opts ...repository.FindOption) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	// 옵션 적용
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}
	
	// 필터링 적용
	var count int64
	for _, asset := range r.assets {
		if asset.IsDeleted {
			continue
		}
		
		if r.matchesFilters(asset, options.Filters) {
			count++
		}
	}
	
	return count, nil
}

// FindByUserID는 사용자 ID로 자산 목록을 조회합니다.
func (r *AssetRepository) FindByUserID(ctx context.Context, userID string, opts ...repository.FindOption) ([]*domain.Asset, error) {
	// 필터에 사용자 ID 추가
	filterOpt := repository.WithFilter("user_id", userID)
	
	// 기존 옵션에 사용자 ID 필터 추가
	opts = append(opts, filterOpt)
	
	// FindAll 메서드 호출
	return r.FindAll(ctx, opts...)
}

// FindByType은 자산 유형으로 자산 목록을 조회합니다.
func (r *AssetRepository) FindByType(ctx context.Context, assetType domain.AssetType, opts ...repository.FindOption) ([]*domain.Asset, error) {
	// 필터에 자산 유형 추가
	filterOpt := repository.WithFilter("type", assetType)
	
	// 기존 옵션에 자산 유형 필터 추가
	opts = append(opts, filterOpt)
	
	// FindAll 메서드 호출
	return r.FindAll(ctx, opts...)
}

// CountByUserID는 사용자 ID에 해당하는 자산 개수를 반환합니다.
func (r *AssetRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	// 필터에 사용자 ID 추가
	filterOpt := repository.WithFilter("user_id", userID)
	
	// Count 메서드 호출
	return r.Count(ctx, filterOpt)
}

// CountByType은 자산 유형에 해당하는 자산 개수를 반환합니다.
func (r *AssetRepository) CountByType(ctx context.Context, assetType domain.AssetType) (int64, error) {
	// 필터에 자산 유형 추가
	filterOpt := repository.WithFilter("type", assetType)
	
	// Count 메서드 호출
	return r.Count(ctx, filterOpt)
}