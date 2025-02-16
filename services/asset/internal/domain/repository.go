package domain

import (
	"context"
	"fmt"
	"sync"
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

// FindByUserID 사용자 ID로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(_ context.Context, userID string) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var assets []*Asset
	for _, asset := range r.assets {
		if asset.UserID == userID {
			assets = append(assets, asset)
		}
	}

	return assets, nil
}

// FindByType 자산 유형으로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByType(_ context.Context, assetType AssetType) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var assets []*Asset
	for _, asset := range r.assets {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}

	return assets, nil
}
