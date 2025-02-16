package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
)

// MemoryAssetRepository는 인메모리 자산 저장소입니다.
type MemoryAssetRepository struct {
	assets   map[string]*domain.Asset
	mu       sync.RWMutex
	eventBus events.EventBus
}

// NewMemoryAssetRepository는 새로운 인메모리 자산 저장소를 생성합니다.
func NewMemoryAssetRepository(eventBus events.EventBus) *MemoryAssetRepository {
	return &MemoryAssetRepository{
		assets:   make(map[string]*domain.Asset),
		eventBus: eventBus,
	}
}

// Save는 자산을 저장합니다.
func (r *MemoryAssetRepository) Save(ctx context.Context, asset *domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[asset.ID]; exists {
		return fmt.Errorf("asset already exists: %s", asset.ID)
	}

	r.assets[asset.ID] = asset

	// 이벤트 발행
	for _, event := range asset.Events() {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}
	asset.ClearEvents()

	return nil
}

// FindByID는 ID로 자산을 조회합니다.
func (r *MemoryAssetRepository) FindByID(ctx context.Context, id string) (*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	asset, exists := r.assets[id]
	if !exists {
		return nil, fmt.Errorf("asset not found: %s", id)
	}

	return asset, nil
}

// Update는 자산을 업데이트합니다.
func (r *MemoryAssetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.assets[asset.ID]; !exists {
		return fmt.Errorf("asset not found: %s", asset.ID)
	}

	r.assets[asset.ID] = asset

	// 이벤트 발행
	for _, event := range asset.Events() {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}
	asset.ClearEvents()

	return nil
}

// Delete는 자산을 삭제합니다.
func (r *MemoryAssetRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	asset, exists := r.assets[id]
	if !exists {
		return fmt.Errorf("asset not found: %s", id)
	}

	asset.MarkAsDeleted()

	// 이벤트 발행
	for _, event := range asset.Events() {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}
	asset.ClearEvents()

	delete(r.assets, id)
	return nil
}

// FindByUserID는 사용자 ID로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var assets []*domain.Asset
	for _, asset := range r.assets {
		if asset.UserID == userID {
			assets = append(assets, asset)
		}
	}

	return assets, nil
}

// FindByType는 자산 유형으로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByType(ctx context.Context, assetType domain.AssetType) ([]*domain.Asset, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var assets []*domain.Asset
	for _, asset := range r.assets {
		if asset.Type == assetType {
			assets = append(assets, asset)
		}
	}

	return assets, nil
}
