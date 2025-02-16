package gamification

import (
	"context"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, profile *Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*Profile, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Profile, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockRepository) FindByUserID(ctx context.Context, userID string) (*Profile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, profile *Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) AddBadge(ctx context.Context, profileID string, badge Badge) error {
	args := m.Called(ctx, profileID, badge)
	return args.Error(0)
}

func (m *MockRepository) UpdateExperience(ctx context.Context, profileID string, experience int) error {
	args := m.Called(ctx, profileID, experience)
	return args.Error(0)
}

func (m *MockRepository) UpdateStats(ctx context.Context, profileID string, stats Statistics) error {
	args := m.Called(ctx, profileID, stats)
	return args.Error(0)
}

func (m *MockRepository) UpdateStreak(ctx context.Context, profileID string, streakType StreakType) error {
	args := m.Called(ctx, profileID, streakType)
	return args.Error(0)
}

func (m *MockRepository) FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Profile, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Profile), args.Error(1)
}

func (m *MockRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
