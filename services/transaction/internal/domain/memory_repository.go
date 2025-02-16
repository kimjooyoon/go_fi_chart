package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var ErrTransactionNotFound = errors.New("transaction not found")

// MemoryTransactionRepository는 인메모리 거래 저장소 구현입니다
type MemoryTransactionRepository struct {
	transactions map[uuid.UUID]*Transaction
	mu           sync.RWMutex
}

// NewMemoryTransactionRepository는 새로운 인메모리 거래 저장소를 생성합니다
func NewMemoryTransactionRepository() *MemoryTransactionRepository {
	return &MemoryTransactionRepository{
		transactions: make(map[uuid.UUID]*Transaction),
	}
}

// Save는 새로운 거래를 저장합니다
func (r *MemoryTransactionRepository) Save(_ context.Context, transaction *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID]; exists {
		return fmt.Errorf("transaction with ID %s already exists", transaction.ID)
	}

	r.transactions[transaction.ID] = transaction
	return nil
}

// FindByID는 ID로 거래를 조회합니다
func (r *MemoryTransactionRepository) FindByID(_ context.Context, id uuid.UUID) (*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	transaction, exists := r.transactions[id]
	if !exists {
		return nil, fmt.Errorf("transaction with ID %s not found", id)
	}

	return transaction, nil
}

// FindByUserID는 사용자 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByUserID(_ context.Context, userID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, t := range r.transactions {
		if t.UserID == userID {
			transactions = append(transactions, t)
		}
	}

	return transactions, nil
}

// FindByPortfolioID는 포트폴리오 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByPortfolioID(_ context.Context, portfolioID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, t := range r.transactions {
		if t.PortfolioID == portfolioID {
			transactions = append(transactions, t)
		}
	}

	return transactions, nil
}

// FindByAssetID는 자산 ID로 거래 목록을 조회합니다
func (r *MemoryTransactionRepository) FindByAssetID(_ context.Context, assetID uuid.UUID) ([]*Transaction, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var transactions []*Transaction
	for _, t := range r.transactions {
		if t.AssetID == assetID {
			transactions = append(transactions, t)
		}
	}

	return transactions, nil
}

// Update는 기존 거래를 업데이트합니다
func (r *MemoryTransactionRepository) Update(_ context.Context, transaction *Transaction) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[transaction.ID]; !exists {
		return fmt.Errorf("transaction with ID %s not found", transaction.ID)
	}

	r.transactions[transaction.ID] = transaction
	return nil
}

// Delete는 거래를 삭제합니다
func (r *MemoryTransactionRepository) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.transactions[id]; !exists {
		return ErrTransactionNotFound
	}

	delete(r.transactions, id)
	return nil
}
