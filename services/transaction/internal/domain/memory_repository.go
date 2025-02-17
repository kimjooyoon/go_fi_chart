package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/google/uuid"
)

var ErrTransactionNotFound = errors.New("transaction not found")

// MemoryTransactionRepository는 인메모리 거래 저장소 구현입니다
type MemoryTransactionRepository struct {
	transactions map[string]*Transaction
	eventBus     events.EventBus
	mu           sync.RWMutex
}

// NewMemoryTransactionRepository는 새로운 인메모리 거래 저장소를 생성합니다
func NewMemoryTransactionRepository(eventBus events.EventBus) *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[string]*Transaction),
		eventBus:     eventBus,
	}
}

// Save는 새로운 거래를 저장합니다
func (r *MemoryTransactionRepository) Save(ctx context.Context, transaction *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID.String()]; exists {
		return fmt.Errorf("transaction already exists: %s", transaction.ID)
	}

	r.transactions[transaction.ID.String()] = transaction
	return nil
}

// FindByID는 ID로 거래를 조회합니다
func (r *MemoryTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if transaction, exists := r.transactions[id.String()]; exists {
		return transaction, nil
	}

	return nil, fmt.Errorf("transaction not found: %s", id)
}

// FindByUserID는 사용자 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, transaction := range r.transactions {
		if transaction.UserID.String() == userID.String() {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

// FindByPortfolioID는 포트폴리오 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, transaction := range r.transactions {
		if transaction.PortfolioID.String() == portfolioID.String() {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

// FindByAssetID는 자산 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByAssetID(ctx context.Context, assetID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, transaction := range r.transactions {
		if transaction.AssetID.String() == assetID.String() {
			transactions = append(transactions, transaction)
		}
	}

	return transactions, nil
}

// Update는 기존 거래를 업데이트합니다
func (r *MemoryTransactionRepository) Update(ctx context.Context, transaction *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID.String()]; !exists {
		return fmt.Errorf("transaction not found: %s", transaction.ID)
	}

	r.transactions[transaction.ID.String()] = transaction
	return nil
}

// Delete는 거래를 삭제합니다
func (r *MemoryTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[id.String()]; !exists {
		return fmt.Errorf("transaction not found: %s", id)
	}

	delete(r.transactions, id.String())
	return nil
}
