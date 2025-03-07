package domain

import (
	"context"
	"sync"

	"github.com/aske/go_fi_chart/internal/common/repository"
)

// MemoryAssetRepository 인메모리 자산 저장소 구현체입니다.
type MemoryAssetRepository struct {
	assets      map[string]*Asset
	userIDIndex map[string]map[string]*Asset    // userID -> assetID -> Asset
	typeIndex   map[AssetType]map[string]*Asset // AssetType -> assetID -> Asset
	mutex       sync.RWMutex
}

// NewMemoryAssetRepository 새로운 인메모리 자산 저장소를 생성합니다.
func NewMemoryAssetRepository() *MemoryAssetRepository {
	return &MemoryAssetRepository{
		assets:      make(map[string]*Asset),
		userIDIndex: make(map[string]map[string]*Asset),
		typeIndex:   make(map[AssetType]map[string]*Asset),
	}
}

// 인덱스를 업데이트하는 내부 메소드
func (r *MemoryAssetRepository) updateIndices(asset *Asset) {
	// 사용자 ID 인덱스 업데이트
	if _, exists := r.userIDIndex[asset.UserID]; !exists {
		r.userIDIndex[asset.UserID] = make(map[string]*Asset)
	}
	r.userIDIndex[asset.UserID][asset.ID] = asset

	// 자산 유형 인덱스 업데이트
	if _, exists := r.typeIndex[asset.Type]; !exists {
		r.typeIndex[asset.Type] = make(map[string]*Asset)
	}
	r.typeIndex[asset.Type][asset.ID] = asset
}

// 인덱스에서 자산을 제거하는 내부 메소드
func (r *MemoryAssetRepository) removeFromIndices(asset *Asset) {
	// 사용자 ID 인덱스에서 제거
	if userAssets, exists := r.userIDIndex[asset.UserID]; exists {
		delete(userAssets, asset.ID)

		// 사용자가 더 이상 자산을 가지고 있지 않으면 사용자 인덱스 항목 제거
		if len(userAssets) == 0 {
			delete(r.userIDIndex, asset.UserID)
		}
	}

	// 자산 유형 인덱스에서 제거
	if typeAssets, exists := r.typeIndex[asset.Type]; exists {
		delete(typeAssets, asset.ID)

		// 특정 유형의 자산이 더 이상 없으면 유형 인덱스 항목 제거
		if len(typeAssets) == 0 {
			delete(r.typeIndex, asset.Type)
		}
	}
}

// Save 자산을 저장합니다.
func (r *MemoryAssetRepository) Save(_ context.Context, asset *Asset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[asset.ID]; exists {
		return NewAssetAlreadyExistsError(asset.ID)
	}

	r.assets[asset.ID] = asset
	r.updateIndices(asset)
	return nil
}

// FindByID ID로 자산을 조회합니다.
func (r *MemoryAssetRepository) FindByID(_ context.Context, id string) (*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	asset, exists := r.assets[id]
	if !exists {
		return nil, NewAssetNotFoundError(id)
	}

	return asset, nil
}

// Update 자산을 업데이트합니다.
func (r *MemoryAssetRepository) Update(_ context.Context, asset *Asset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	oldAsset, exists := r.assets[asset.ID]
	if !exists {
		return NewAssetNotFoundError(asset.ID)
	}

	// 인덱스에서 이전 자산 정보 제거
	r.removeFromIndices(oldAsset)

	// 새 자산 정보로 업데이트
	r.assets[asset.ID] = asset
	r.updateIndices(asset)
	return nil
}

// Delete ID로 자산을 삭제합니다.
func (r *MemoryAssetRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	asset, exists := r.assets[id]
	if !exists {
		return NewAssetNotFoundError(id)
	}

	// 인덱스에서 제거
	r.removeFromIndices(asset)

	// 자산 맵에서 제거
	delete(r.assets, id)
	return nil
}

// applyPagination 자산 목록에 페이지네이션을 적용하는 내부 도우미 함수입니다.
func applyPagination(assets []*Asset, options *repository.FindOptions) []*Asset {
	if options.Limit <= 0 {
		return assets
	}

	offset := int64(options.Offset)
	limit := int64(options.Limit)

	if offset >= int64(len(assets)) {
		return []*Asset{}
	}

	end := offset + limit
	if end > int64(len(assets)) {
		end = int64(len(assets))
	}

	return assets[offset:end]
}

// FindAll 모든 자산을 조회합니다. 옵션을 통해 필터링, 정렬, 페이지네이션을 적용할 수 있습니다.
func (r *MemoryAssetRepository) FindAll(_ context.Context, opts ...repository.FindOption) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 옵션 적용
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}

	// 결과 배열 생성
	var result []*Asset
	for _, asset := range r.assets {
		// 필터링 로직은 향후 구현
		result = append(result, asset)
	}

	// 페이지네이션 적용
	return applyPagination(result, options), nil
}

// Count 조건에 맞는 자산의 총 개수를 반환합니다.
func (r *MemoryAssetRepository) Count(_ context.Context, opts ...repository.FindOption) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 필터링 로직은 향후 구현
	// 지금은 모든 자산의 개수를 반환
	return int64(len(r.assets)), nil
}

// FindByUserID 사용자 ID로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(_ context.Context, userID string, opts ...repository.FindOption) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 옵션 적용
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}

	// 인덱스를 활용하여 사용자의 자산 가져오기
	var assets []*Asset
	if userAssets, exists := r.userIDIndex[userID]; exists {
		for _, asset := range userAssets {
			assets = append(assets, asset)
		}
	}

	// 페이지네이션 적용
	return applyPagination(assets, options), nil
}

// FindByType 자산 유형으로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByType(_ context.Context, assetType AssetType, opts ...repository.FindOption) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 옵션 적용
	options := repository.NewFindOptions()
	for _, opt := range opts {
		opt.Apply(options)
	}

	// 인덱스를 활용하여 특정 유형의 자산 가져오기
	var assets []*Asset
	if typeAssets, exists := r.typeIndex[assetType]; exists {
		for _, asset := range typeAssets {
			assets = append(assets, asset)
		}
	}

	// 페이지네이션 적용
	return applyPagination(assets, options), nil
}

// CountByUserID 사용자 ID에 해당하는 자산 개수를 반환합니다.
func (r *MemoryAssetRepository) CountByUserID(_ context.Context, userID string) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 인덱스를 활용하여 사용자의 자산 개수 계산
	if userAssets, exists := r.userIDIndex[userID]; exists {
		return int64(len(userAssets)), nil
	}
	return 0, nil
}

// CountByType 자산 유형에 해당하는 자산 개수를 반환합니다.
func (r *MemoryAssetRepository) CountByType(_ context.Context, assetType AssetType) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// 인덱스를 활용하여 특정 유형의 자산 개수 계산
	if typeAssets, exists := r.typeIndex[assetType]; exists {
		return int64(len(typeAssets)), nil
	}
	return 0, nil
}
