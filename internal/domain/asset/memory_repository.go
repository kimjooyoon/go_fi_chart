package asset

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// MemoryRepository 인메모리 저장소 구현체
type MemoryRepository[T any] struct {
	data  map[string]T
	mutex sync.RWMutex
}

// NewMemoryRepository 새로운 인메모리 저장소를 생성합니다.
func NewMemoryRepository[T any]() *MemoryRepository[T] {
	return &MemoryRepository[T]{
		data: make(map[string]T),
	}
}

// Save 엔티티를 저장합니다.
func (r *MemoryRepository[T]) Save(_ context.Context, entity interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if e, ok := entity.(domain.Entity); ok {
		if _, exists := r.data[e.GetID()]; exists {
			return domain.NewRepositoryError("Save", fmt.Errorf("entity with ID %s already exists", e.GetID()))
		}

		if t, ok := entity.(T); ok {
			r.data[e.GetID()] = t
			return nil
		}
	}
	return domain.NewRepositoryError("Save", fmt.Errorf("invalid entity type"))
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
func (r *MemoryRepository[T]) Update(_ context.Context, entity interface{}) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if e, ok := entity.(domain.Entity); ok {
		if _, exists := r.data[e.GetID()]; !exists {
			return domain.NewRepositoryError("Update", fmt.Errorf("entity with ID %s not found", e.GetID()))
		}

		if t, ok := entity.(T); ok {
			r.data[e.GetID()] = t
			return nil
		}
	}
	return domain.NewRepositoryError("Update", fmt.Errorf("invalid entity type"))
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

// MemoryAssetRepository Asset 도메인의 인메모리 저장소 구현체입니다.
type MemoryAssetRepository struct {
	repo         *MemoryRepository[*Asset]
	transactions map[string]*Transaction
	mutex        sync.RWMutex
}

// NewMemoryAssetRepository 새로운 MemoryAssetRepository를 생성합니다.
func NewMemoryAssetRepository() *MemoryAssetRepository {
	return &MemoryAssetRepository{
		repo:         NewMemoryRepository[*Asset](),
		transactions: make(map[string]*Transaction),
	}
}

// Save Asset을 저장합니다.
func (r *MemoryAssetRepository) Save(ctx context.Context, asset *Asset) error {
	return r.repo.Save(ctx, asset)
}

// FindByID ID로 Asset을 조회합니다.
func (r *MemoryAssetRepository) FindByID(ctx context.Context, id string) (*Asset, error) {
	return r.repo.FindByID(ctx, id)
}

// Update Asset을 업데이트합니다.
func (r *MemoryAssetRepository) Update(ctx context.Context, asset *Asset) error {
	return r.repo.Update(ctx, asset)
}

// Delete ID로 Asset을 삭제합니다.
func (r *MemoryAssetRepository) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

// FindByUserID 사용자 ID로 Asset 목록을 조회합니다.
func (r *MemoryAssetRepository) FindByUserID(_ context.Context, userID string) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Asset
	for _, asset := range r.repo.data {
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
	for _, asset := range r.repo.data {
		if asset.Type == assetType {
			result = append(result, asset)
		}
	}
	return result, nil
}

// UpdateAmount 자산의 금액을 업데이트합니다.
func (r *MemoryAssetRepository) UpdateAmount(_ context.Context, id string, amount Money) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	asset, exists := r.repo.data[id]
	if !exists {
		return domain.NewRepositoryError("UpdateAmount", fmt.Errorf("asset with ID %s not found", id))
	}

	asset.Amount = amount
	asset.UpdatedAt = time.Now()
	r.repo.data[id] = asset

	return nil
}

// FindAll 검색 조건에 맞는 모든 Asset을 조회합니다.
func (r *MemoryAssetRepository) FindAll(_ context.Context, _ domain.SearchCriteria) ([]*Asset, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Asset
	for _, asset := range r.repo.data {
		result = append(result, asset)
	}
	return result, nil
}

// FindOne 검색 조건에 맞는 하나의 Asset을 조회합니다.
func (r *MemoryAssetRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Asset, error) {
	return r.repo.FindOne(ctx, criteria)
}

// WithTransaction 트랜잭션을 실행합니다.
func (r *MemoryAssetRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.repo.WithTransaction(ctx, fn)
}

// SaveTransaction 거래를 저장합니다.
func (r *MemoryAssetRepository) SaveTransaction(_ context.Context, tx *Transaction) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.transactions[tx.ID] = tx
	return nil
}

// FindTransactionByID ID로 거래를 찾습니다.
func (r *MemoryAssetRepository) FindTransactionByID(_ context.Context, id string) (*Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tx, exists := r.transactions[id]
	if !exists {
		return nil, domain.NewRepositoryError("FindTransactionByID", fmt.Errorf("transaction with ID %s not found", id))
	}
	return tx, nil
}

// FindTransactionsByDateRange 날짜 범위로 거래를 찾습니다.
func (r *MemoryAssetRepository) FindTransactionsByDateRange(_ context.Context, start, end time.Time) ([]*Transaction, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.transactions {
		if (tx.Date.Equal(start) || tx.Date.After(start)) && (tx.Date.Equal(end) || tx.Date.Before(end)) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// CalculateTotalAmount 특정 통화의 총 금액을 계산합니다.
func (r *MemoryAssetRepository) CalculateTotalAmount(_ context.Context, currency string) (Money, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	total := NewTestMoney(0, currency)
	for _, tx := range r.transactions {
		if tx.Amount.Currency == currency {
			switch tx.Type {
			case Income:
				result, err := total.Add(tx.Amount)
				if err != nil {
					return Money{}, err
				}
				total = result
			case Expense:
				result, err := total.Subtract(tx.Amount)
				if err != nil {
					return Money{}, err
				}
				total = result
			}
		}
	}
	return total, nil
}

// MemoryTransactionRepository Transaction 도메인의 인메모리 저장소 구현체입니다.
type MemoryTransactionRepository struct {
	repo *MemoryRepository[*Transaction]
}

// NewMemoryTransactionRepository 새로운 인메모리 Transaction 저장소를 생성합니다.
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		repo: NewMemoryRepository[*Transaction](),
	}
}

// Save Transaction을 저장합니다.
func (r *MemoryTransactionRepository) Save(ctx context.Context, tx *Transaction) error {
	return r.repo.Save(ctx, tx)
}

// FindByID ID로 Transaction을 조회합니다.
func (r *MemoryTransactionRepository) FindByID(ctx context.Context, id string) (*Transaction, error) {
	return r.repo.FindByID(ctx, id)
}

// Update Transaction을 업데이트합니다.
func (r *MemoryTransactionRepository) Update(ctx context.Context, tx *Transaction) error {
	return r.repo.Update(ctx, tx)
}

// Delete ID로 Transaction을 삭제합니다.
func (r *MemoryTransactionRepository) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

// FindByAssetID AssetID로 Transaction 목록을 조회합니다.
func (r *MemoryTransactionRepository) FindByAssetID(_ context.Context, assetID string) ([]*Transaction, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.repo.data {
		if tx.AssetID == assetID {
			result = append(result, tx)
		}
	}
	return result, nil
}

// FindByDateRange 날짜 범위로 Transaction 목록을 조회합니다.
func (r *MemoryTransactionRepository) FindByDateRange(_ context.Context, start, end time.Time) ([]*Transaction, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.repo.data {
		if (tx.Date.Equal(start) || tx.Date.After(start)) &&
			(tx.Date.Equal(end) || tx.Date.Before(end)) {
			result = append(result, tx)
		}
	}
	return result, nil
}

// GetTotalAmount 특정 기간 동안의 총 거래 금액을 계산합니다.
func (r *MemoryTransactionRepository) GetTotalAmount(_ context.Context, assetID string) (Money, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	total := Money{Amount: 0, Currency: "KRW"}
	for _, tx := range r.repo.data {
		if tx.AssetID == assetID {
			switch tx.Type {
			case Income:
				result, err := total.Add(tx.Amount)
				if err != nil {
					return Money{}, err
				}
				total = result
			case Expense:
				result, err := total.Subtract(tx.Amount)
				if err != nil {
					return Money{}, err
				}
				total = result
			}
		}
	}
	return total, nil
}

// FindAll 검색 조건에 맞는 모든 Transaction을 조회합니다.
func (r *MemoryTransactionRepository) FindAll(_ context.Context, _ domain.SearchCriteria) ([]*Transaction, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	var result []*Transaction
	for _, tx := range r.repo.data {
		result = append(result, tx)
	}
	return result, nil
}

// FindOne 검색 조건에 맞는 하나의 Transaction을 조회합니다.
func (r *MemoryTransactionRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Transaction, error) {
	return r.repo.FindOne(ctx, criteria)
}

// WithTransaction 트랜잭션을 실행합니다.
func (r *MemoryTransactionRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.repo.WithTransaction(ctx, fn)
}

// MemoryPortfolioRepository Portfolio 도메인의 인메모리 저장소 구현체입니다.
type MemoryPortfolioRepository struct {
	repo *MemoryRepository[*Portfolio]
}

// NewMemoryPortfolioRepository 새로운 인메모리 Portfolio 저장소를 생성합니다.
func NewMemoryPortfolioRepository() *MemoryPortfolioRepository {
	return &MemoryPortfolioRepository{
		repo: NewMemoryRepository[*Portfolio](),
	}
}

// Save Portfolio를 저장합니다.
func (r *MemoryPortfolioRepository) Save(ctx context.Context, portfolio *Portfolio) error {
	return r.repo.Save(ctx, portfolio)
}

// FindByID ID로 Portfolio를 조회합니다.
func (r *MemoryPortfolioRepository) FindByID(ctx context.Context, id string) (*Portfolio, error) {
	return r.repo.FindByID(ctx, id)
}

// Update Portfolio를 업데이트합니다.
func (r *MemoryPortfolioRepository) Update(ctx context.Context, portfolio *Portfolio) error {
	return r.repo.Update(ctx, portfolio)
}

// Delete ID로 Portfolio를 삭제합니다.
func (r *MemoryPortfolioRepository) Delete(ctx context.Context, id string) error {
	return r.repo.Delete(ctx, id)
}

// FindByUserID 사용자 ID로 Portfolio를 조회합니다.
func (r *MemoryPortfolioRepository) FindByUserID(_ context.Context, userID string) (*Portfolio, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	for _, portfolio := range r.repo.data {
		if portfolio.UserID == userID {
			return portfolio, nil
		}
	}
	return nil, domain.NewRepositoryError("FindByUserID", fmt.Errorf("portfolio for user %s not found", userID))
}

// UpdateAssets Portfolio의 자산 구성을 업데이트합니다.
func (r *MemoryPortfolioRepository) UpdateAssets(_ context.Context, id string, assets []PortfolioAsset) error {
	r.repo.mutex.Lock()
	defer r.repo.mutex.Unlock()

	portfolio, ok := r.repo.data[id]
	if !ok {
		return domain.NewRepositoryError("UpdateAssets", fmt.Errorf("portfolio with ID %s not found", id))
	}

	portfolio.Assets = assets
	portfolio.UpdatedAt = time.Now()
	r.repo.data[id] = portfolio

	return nil
}

// FindAll 검색 조건에 맞는 모든 Portfolio를 조회합니다.
func (r *MemoryPortfolioRepository) FindAll(_ context.Context, _ domain.SearchCriteria) ([]*Portfolio, error) {
	r.repo.mutex.RLock()
	defer r.repo.mutex.RUnlock()

	var result []*Portfolio
	for _, portfolio := range r.repo.data {
		result = append(result, portfolio)
	}
	return result, nil
}

// FindOne 검색 조건에 맞는 하나의 Portfolio를 조회합니다.
func (r *MemoryPortfolioRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Portfolio, error) {
	return r.repo.FindOne(ctx, criteria)
}

// WithTransaction 트랜잭션을 실행합니다.
func (r *MemoryPortfolioRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return r.repo.WithTransaction(ctx, fn)
}
