package asset

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// MemoryRepository 인메모리 저장소 구현체
type MemoryRepository[T domain.Entity] struct {
	data  map[string]T
	mutex sync.RWMutex
}

// NewMemoryRepository 새로운 인메모리 저장소를 생성합니다.
func NewMemoryRepository[T domain.Entity]() *MemoryRepository[T] {
	return &MemoryRepository[T]{
		data: make(map[string]T),
	}
}

// Save 엔티티를 저장합니다.
func (r *MemoryRepository[T]) Save(_ context.Context, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[entity.GetID()]; exists {
		return domain.NewRepositoryError("Save", fmt.Errorf("entity with ID %s already exists", entity.GetID()))
	}

	r.data[entity.GetID()] = entity
	return nil
}

// FindByID ID로 엔티티를 조회합니다.
func (r *MemoryRepository[T]) FindByID(_ context.Context, id string) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if entity, exists := r.data[id]; exists {
		return entity, nil
	}

	var zero T
	return zero, domain.NewRepositoryError("FindByID", fmt.Errorf("entity with ID %s not found", id))
}

// Update 엔티티를 업데이트합니다.
func (r *MemoryRepository[T]) Update(_ context.Context, entity T) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[entity.GetID()]; !exists {
		return domain.NewRepositoryError("Update", fmt.Errorf("entity with ID %s not found", entity.GetID()))
	}

	r.data[entity.GetID()] = entity
	return nil
}

// Delete ID로 엔티티를 삭제합니다.
func (r *MemoryRepository[T]) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[id]; !exists {
		return domain.NewRepositoryError("Delete", fmt.Errorf("entity with ID %s not found", id))
	}

	delete(r.data, id)
	return nil
}

// FindAll 검색 조건에 맞는 모든 엔티티를 조회합니다.
func (r *MemoryRepository[T]) FindAll(_ context.Context, _ domain.SearchCriteria) ([]T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []T
	for _, entity := range r.data {
		result = append(result, entity)
	}
	return result, nil
}

// FindOne 검색 조건에 맞는 하나의 엔티티를 조회합니다.
func (r *MemoryRepository[T]) FindOne(_ context.Context, _ domain.SearchCriteria) (T, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var zero T
	return zero, domain.NewRepositoryError("FindOne", fmt.Errorf("not implemented"))
}

// WithTransaction 트랜잭션을 실행합니다.
func (r *MemoryRepository[T]) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// MemoryAssetRepository Asset 도메인의 인메모리 저장소
type MemoryAssetRepository struct {
	*MemoryRepository[*Asset]
}

// NewMemoryAssetRepository 새로운 Asset 인메모리 저장소를 생성합니다.
func NewMemoryAssetRepository() *MemoryAssetRepository {
	return &MemoryAssetRepository{
		MemoryRepository: NewMemoryRepository[*Asset](),
	}
}

// FindByUserID 사용자 ID로 Asset 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(_ context.Context, userID string) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Asset
	for _, asset := range r.data {
		if asset.UserID == userID {
			result = append(result, asset)
		}
	}
	return result, nil
}

// FindByType Asset 유형으로 Asset 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByType(_ context.Context, assetType Type) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Asset
	for _, asset := range r.data {
		if asset.Type == assetType {
			result = append(result, asset)
		}
	}
	return result, nil
}

// UpdateAmount Asset의 금액을 업데이트합니다.
func (r *MemoryAssetRepository) UpdateAmount(_ context.Context, id string, amount float64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	asset, ok := r.data[id]
	if !ok {
		return domain.NewRepositoryError("UpdateAmount", fmt.Errorf("asset with ID %s not found", id))
	}

	asset.Amount = Money{Amount: amount, Currency: asset.Amount.Currency}
	asset.UpdatedAt = time.Now()
	r.data[id] = asset

	return nil
}

// MemoryTransactionRepository Transaction 도메인의 인메모리 저장소
type MemoryTransactionRepository struct {
	*MemoryRepository[*Transaction]
}

// NewMemoryTransactionRepository 새로운 Transaction 인메모리 저장소를 생성합니다.
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		MemoryRepository: NewMemoryRepository[*Transaction](),
	}
}

// FindByAssetID 자산 ID로 Transaction 목록을 조회합니다.
func (r *MemoryTransactionRepository) FindByAssetID(_ context.Context, assetID string) ([]*Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.data {
		if tx.AssetID == assetID {
			result = append(result, tx)
		}
	}
	return result, nil
}

// FindByDateRange 날짜 범위로 Transaction 목록을 조회합니다.
func (r *MemoryTransactionRepository) FindByDateRange(_ context.Context, start, end time.Time) ([]*Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.data {
		if (tx.Date.Equal(start) || tx.Date.After(start)) && (tx.Date.Equal(end) || tx.Date.Before(end)) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// GetTotalAmount 자산 ID에 대한 총 거래 금액을 계산합니다.
func (r *MemoryTransactionRepository) GetTotalAmount(_ context.Context, assetID string) (float64, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var total float64
	for _, tx := range r.data {
		if tx.AssetID == assetID {
			switch tx.Type {
			case Income:
				total += tx.Amount
			case Expense:
				total -= tx.Amount
			case Transfer:
				// Transfer는 별도 처리 필요
			}
		}
	}
	return total, nil
}

// MemoryPortfolioRepository Portfolio 도메인의 인메모리 저장소
type MemoryPortfolioRepository struct {
	*MemoryRepository[*Portfolio]
}

// NewMemoryPortfolioRepository 새로운 Portfolio 인메모리 저장소를 생성합니다.
func NewMemoryPortfolioRepository() *MemoryPortfolioRepository {
	return &MemoryPortfolioRepository{
		MemoryRepository: NewMemoryRepository[*Portfolio](),
	}
}

// FindByUserID 사용자 ID로 Portfolio를 조회합니다.
func (r *MemoryPortfolioRepository) FindByUserID(_ context.Context, userID string) (*Portfolio, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, portfolio := range r.data {
		if portfolio.UserID == userID {
			return portfolio, nil
		}
	}
	return nil, domain.NewRepositoryError("FindByUserID", fmt.Errorf("portfolio for user %s not found", userID))
}

// UpdateAssets Portfolio의 자산 구성을 업데이트합니다.
func (r *MemoryPortfolioRepository) UpdateAssets(_ context.Context, id string, assets []PortfolioAsset) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	portfolio, exists := r.data[id]
	if !exists {
		return domain.NewRepositoryError("UpdateAssets", fmt.Errorf("portfolio with ID %s not found", id))
	}

	portfolio.Assets = assets
	portfolio.UpdatedAt = time.Now()
	r.data[id] = portfolio
	return nil
}
