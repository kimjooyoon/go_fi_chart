package asset

import (
	"context"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, asset *Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Asset), args.Error(1)
}

func (m *MockRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Asset, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Asset), args.Error(1)
}

func (m *MockRepository) FindByUserID(ctx context.Context, userID string) ([]*Asset, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Asset), args.Error(1)
}

func (m *MockRepository) FindByType(ctx context.Context, assetType Type) ([]*Asset, error) {
	args := m.Called(ctx, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Asset), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, asset *Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockRepository) UpdateAmount(ctx context.Context, id string, amount Money) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Asset, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Asset), args.Error(1)
}

func (m *MockRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(ctx context.Context, transaction *Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id string) (*Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByUserID(ctx context.Context, userID string) ([]*Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByAssetID(ctx context.Context, assetID string) ([]*Transaction, error) {
	args := m.Called(ctx, assetID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Transaction, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*Transaction, error) {
	args := m.Called(ctx, startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Transaction, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetTotalAmount(ctx context.Context, userID string) (Money, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return Money{}, args.Error(1)
	}
	return args.Get(0).(Money), args.Error(1)
}

func (m *MockTransactionRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) Save(ctx context.Context, portfolio *Portfolio) error {
	args := m.Called(ctx, portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) FindByID(ctx context.Context, id string) (*Portfolio, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*Portfolio, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) FindByUserID(ctx context.Context, userID string) (*Portfolio, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) Update(ctx context.Context, portfolio *Portfolio) error {
	args := m.Called(ctx, portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPortfolioRepository) FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*Portfolio, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) UpdateAssets(ctx context.Context, id string, assets []PortfolioAsset) error {
	args := m.Called(ctx, id, assets)
	return args.Error(0)
}

func (m *MockPortfolioRepository) WithTransaction(ctx context.Context, fn func(context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}
