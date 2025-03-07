package infrastructure

import (
	"context"
	"fmt"
	"sync"

	"github.com/aske/go_fi_chart/internal/common/repository"
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
	events := asset.Events()

	r.mu.Lock()
	if _, exists := r.assets[asset.ID]; exists {
		r.mu.Unlock()
		return fmt.Errorf("asset already exists: %s", asset.ID)
	}

	r.assets[asset.ID] = asset
	r.mu.Unlock()

	// 락 해제 후 이벤트 발행
	for _, event := range events {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}
	asset.ClearEvents()

	return nil
}

// FindByID는 ID로 자산을 조회합니다.
func (r *MemoryAssetRepository) FindByID(_ context.Context, id string) (*domain.Asset, error) {
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
	events := asset.Events()

	r.mu.Lock()
	if _, exists := r.assets[asset.ID]; !exists {
		r.mu.Unlock()
		return fmt.Errorf("asset not found: %s", asset.ID)
	}

	r.assets[asset.ID] = asset
	r.mu.Unlock()

	// 락 해제 후 이벤트 발행
	for _, event := range events {
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
		return domain.NewAssetNotFoundError(id)
	}

	asset.MarkAsDeleted()
	delete(r.assets, id)

	// 이벤트 발행
	for _, event := range asset.Events() {
		if err := r.eventBus.Publish(ctx, event); err != nil {
			return fmt.Errorf("이벤트 발행 실패: %w", err)
		}
	}

	return nil
}

// FindAll은 모든 자산을 조회합니다. 옵션을 통해 필터링, 정렬, 페이지네이션을 적용할 수 있습니다.
func (r *MemoryAssetRepository) FindAll(_ context.Context, opts ...repository.FindOption) ([]*domain.Asset, error) {
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
		// 필터링 로직은 향후 구현
		result = append(result, asset)
	}

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

// Count는 조건에 맞는 자산의 총 개수를 반환합니다.
func (r *MemoryAssetRepository) Count(_ context.Context, opts ...repository.FindOption) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// 필터링 로직은 향후 구현
	// 지금은 모든 자산의 개수를 반환
	return int64(len(r.assets)), nil
}

// FindByUserID는 사용자 ID로 자산 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(_ context.Context, userID string) ([]*domain.Asset, error) {
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
func (r *MemoryAssetRepository) FindByType(_ context.Context, assetType domain.AssetType) ([]*domain.Asset, error) {
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
