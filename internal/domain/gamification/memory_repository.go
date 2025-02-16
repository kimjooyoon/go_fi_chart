package gamification

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
)

// MemoryRepository 게임화 프로필의 인메모리 저장소 구현체입니다.
type MemoryRepository struct {
	data  map[string]*Profile
	mutex sync.RWMutex
}

// NewMemoryRepository 새로운 인메모리 저장소를 생성합니다.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		data: make(map[string]*Profile),
	}
}

// Save 프로필을 저장합니다.
func (r *MemoryRepository) Save(_ context.Context, profile *Profile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[profile.ID]; exists {
		return domain.NewError("gamification", domain.ErrCodeAlreadyExists, fmt.Sprintf("profile with ID %s already exists", profile.ID))
	}

	r.data[profile.ID] = profile
	return nil
}

// FindByID ID로 프로필을 조회합니다.
func (r *MemoryRepository) FindByID(_ context.Context, id string) (*Profile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if profile, exists := r.data[id]; exists {
		return profile, nil
	}

	return nil, domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
}

// Update 프로필을 업데이트합니다.
func (r *MemoryRepository) Update(_ context.Context, profile *Profile) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[profile.ID]; !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", profile.ID))
	}

	profile.UpdatedAt = time.Now()
	r.data[profile.ID] = profile
	return nil
}

// Delete ID로 프로필을 삭제합니다.
func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.data[id]; !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
	}

	delete(r.data, id)
	return nil
}

// FindAll 검색 조건에 맞는 모든 프로필을 조회합니다.
func (r *MemoryRepository) FindAll(_ context.Context, _ domain.SearchCriteria) ([]*Profile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	profiles := make([]*Profile, 0, len(r.data))
	for _, profile := range r.data {
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// FindOne 검색 조건에 맞는 하나의 프로필을 조회합니다.
func (r *MemoryRepository) FindOne(_ context.Context, _ domain.SearchCriteria) (*Profile, error) {
	return nil, domain.NewError("gamification", domain.ErrCodeNotImplemented, "FindOne not implemented")
}

// WithTransaction 트랜잭션을 실행합니다.
func (r *MemoryRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// FindByUserID 사용자 ID로 프로필을 조회합니다.
func (r *MemoryRepository) FindByUserID(_ context.Context, userID string) (*Profile, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, profile := range r.data {
		if profile.UserID == userID {
			return profile, nil
		}
	}

	return nil, domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile for user %s not found", userID))
}

// UpdateExperience 사용자의 경험치를 업데이트합니다.
func (r *MemoryRepository) UpdateExperience(_ context.Context, id string, exp int) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	profile, exists := r.data[id]
	if !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
	}

	leveledUp := profile.AddExperience(exp)
	if leveledUp {
		profile.UpdatedAt = time.Now()
	}

	r.data[id] = profile
	return nil
}

// UpdateStats 사용자의 통계를 업데이트합니다.
func (r *MemoryRepository) UpdateStats(_ context.Context, id string, stats Statistics) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	profile, exists := r.data[id]
	if !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
	}

	profile.Stats = stats
	profile.UpdatedAt = time.Now()
	r.data[id] = profile
	return nil
}

// AddBadge 사용자에게 뱃지를 추가합니다.
func (r *MemoryRepository) AddBadge(_ context.Context, id string, badge Badge) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	profile, exists := r.data[id]
	if !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
	}

	profile.Badges = append(profile.Badges, badge)
	profile.Stats.BadgesEarned++
	profile.UpdatedAt = time.Now()
	r.data[id] = profile
	return nil
}

// UpdateStreak 사용자의 연속 달성을 업데이트합니다.
func (r *MemoryRepository) UpdateStreak(_ context.Context, id string, streakType StreakType) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	profile, exists := r.data[id]
	if !exists {
		return domain.NewError("gamification", domain.ErrCodeNotFound, fmt.Sprintf("profile with ID %s not found", id))
	}

	profile.UpdateStreak(streakType)
	r.data[id] = profile
	return nil
}
