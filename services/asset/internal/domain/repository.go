package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/aske/go_fi_chart/internal/common/repository"
)

// MemoryAssetRepository 인메모리 자산 저장소 구현체입니다.
type MemoryAssetRepository struct {
	assets map[string]*Asset
	mutex  sync.RWMutex
}

// NewMemoryAssetRepository 새로운 인메모리 자산 저장소를 생성합니다.
func NewMemoryAssetRepository() *MemoryAssetRepository {
	return &MemoryAssetRepository{
		assets: make(map[string]*Asset),
	}
}

// Save 자산을 저장합니다.
func (r *MemoryAssetRepository) Save(_ context.Context, asset *Asset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[asset.ID]; exists {
		return fmt.Errorf("자산이 이미 존재합니다: %s", asset.ID)
	}

	r.assets[asset.ID] = asset
	return nil
}

// FindByID ID로 자산을 조회합니다.
func (r *MemoryAssetRepository) FindByID(_ context.Context, id string) (*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	asset, exists := r.assets[id]
	if !exists {
		return nil, fmt.Errorf("자산을 찾을 수 없습니다: %s", id)
	}

	return asset, nil
}

// Update 자산을 업데이트합니다.
func (r *MemoryAssetRepository) Update(_ context.Context, asset *Asset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[asset.ID]; !exists {
		return fmt.Errorf("자산을 찾을 수 없습니다: %s", asset.ID)
	}

	r.assets[asset.ID] = asset
	return nil
}

// Delete ID로 자산을 삭제합니다.
func (r *MemoryAssetRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.assets[id]; !exists {
		return fmt.Errorf("자산을 찾을 수 없습니다: %s", id)
	}

	delete(r.assets, id)
	return nil
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
	if options.Limit > 0 {
		offset := int64(options.Offset)
		limit := int64(options.Limit)

		if offset >= int64(len(result)) {
			return []*Asset{}, nil
		}

		end := offset + limit
		if end > int64(len(result)) {
			end = int64(len(result))
		}

		return result[offset:end], nil
	}

	return result, nil
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

	var assets []*Asset
	for _, asset := range r.assets {
		if asset.UserID == userID {
			assets = append(assets, asset)
		}
	}

	// 페이지네이션 적용
	if options.Limit > 0 {
		offset := int64(options.Offset)
		limit := int64(options.Limit)

		if offset >= int64(len(assets)) {
			return []*Asset{}, nil
		}

		end := offset + limit
		if end > int64(len(assets)) {
			end = int64(len(assets))
		}

		return assets[offset:end], nil
	}

	return assets, nil
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

	var assets []*Asset
	for _, asset := range r.assets {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}

	// 페이지네이션 적용
	if options.Limit > 0 {
		offset := int64(options.Offset)
		limit := int64(options.Limit)

		if offset >= int64(len(assets)) {
			return []*Asset{}, nil
		}

		end := offset + limit
		if end > int64(len(assets)) {
			end = int64(len(assets))
		}

		return assets[offset:end], nil
	}

	return assets, nil
}

// CountByUserID는 사용자 ID에 해당하는 자산 개수를 반환합니다.
func (r *MemoryAssetRepository) CountByUserID(_ context.Context, userID string) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var count int64
	for _, asset := range r.assets {
		if asset.UserID == userID {
			count++
		}
	}

	return count, nil
}

// CountByType은 자산 유형에 해당하는 자산 개수를 반환합니다.
func (r *MemoryAssetRepository) CountByType(_ context.Context, assetType AssetType) (int64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var count int64
	for _, asset := range r.assets {
		if asset.Type == assetType {
			count++
		}
	}

	return count, nil
}
