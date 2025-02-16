package domain

import (
	"context"
	"fmt"
	"sync"
)

// MemoryPortfolioRepository 인메모리 포트폴리오 저장소 구현체입니다.
type MemoryPortfolioRepository struct {
	portfolios map[string]*Portfolio
	mutex      sync.RWMutex
}

// NewMemoryPortfolioRepository 새로운 인메모리 포트폴리오 저장소를 생성합니다.
func NewMemoryPortfolioRepository() *MemoryPortfolioRepository {
	return &MemoryPortfolioRepository{
		portfolios: make(map[string]*Portfolio),
	}
}

// Save 포트폴리오를 저장합니다.
func (r *MemoryPortfolioRepository) Save(_ context.Context, portfolio *Portfolio) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.portfolios[portfolio.ID]; exists {
		return fmt.Errorf("포트폴리오가 이미 존재합니다: %s", portfolio.ID)
	}

	r.portfolios[portfolio.ID] = portfolio
	return nil
}

// FindByID ID로 포트폴리오를 조회합니다.
func (r *MemoryPortfolioRepository) FindByID(_ context.Context, id string) (*Portfolio, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	portfolio, exists := r.portfolios[id]
	if !exists {
		return nil, fmt.Errorf("포트폴리오를 찾을 수 없습니다: %s", id)
	}

	return portfolio, nil
}

// Update 포트폴리오를 업데이트합니다.
func (r *MemoryPortfolioRepository) Update(_ context.Context, portfolio *Portfolio) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.portfolios[portfolio.ID]; !exists {
		return fmt.Errorf("포트폴리오를 찾을 수 없습니다: %s", portfolio.ID)
	}

	r.portfolios[portfolio.ID] = portfolio
	return nil
}

// Delete ID로 포트폴리오를 삭제합니다.
func (r *MemoryPortfolioRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.portfolios[id]; !exists {
		return fmt.Errorf("포트폴리오를 찾을 수 없습니다: %s", id)
	}

	delete(r.portfolios, id)
	return nil
}

// FindByUserID 사용자 ID로 포트폴리오 목록을 조회합니다.
func (r *MemoryPortfolioRepository) FindByUserID(_ context.Context, userID string) ([]*Portfolio, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var portfolios []*Portfolio
	for _, portfolio := range r.portfolios {
		if portfolio.UserID == userID {
			portfolios = append(portfolios, portfolio)
		}
	}

	return portfolios, nil
}
