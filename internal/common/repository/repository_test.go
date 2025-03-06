package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/aske/go_fi_chart/internal/common/repository"
)

// TestEntity 테스트용 엔티티
type TestEntity struct {
	ID   string
	Name string
}

// MockRepository 테스트를 위한 모의 레포지토리 구현
type MockRepository struct {
	mock.Mock
}

// FindByID 메서드 구현
func (m *MockRepository) FindByID(ctx context.Context, id string) (TestEntity, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(TestEntity), args.Error(1)
}

// FindAll 메서드 구현
func (m *MockRepository) FindAll(ctx context.Context, opts ...repository.FindOption) ([]TestEntity, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).([]TestEntity), args.Error(1)
}

// Save 메서드 구현
func (m *MockRepository) Save(ctx context.Context, entity TestEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

// Update 메서드 구현
func (m *MockRepository) Update(ctx context.Context, entity TestEntity) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

// Delete 메서드 구현
func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Count 메서드 구현
func (m *MockRepository) Count(ctx context.Context, opts ...repository.FindOption) (int64, error) {
	args := m.Called(ctx, opts)
	return args.Get(0).(int64), args.Error(1)
}

// TestRepositoryInterface Repository 인터페이스 테스트
func TestRepositoryInterface(t *testing.T) {
	// 타입 확인: MockRepository가 Repository 인터페이스를 구현하는지 확인
	var _ repository.Repository[TestEntity, string] = (*MockRepository)(nil)

	ctx := context.Background()
	mockRepo := new(MockRepository)

	entity := TestEntity{ID: "1", Name: "Test Entity"}

	t.Run("FindByID", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, "1").Return(entity, nil).Once()

		result, err := mockRepo.FindByID(ctx, "1")

		assert.NoError(t, err)
		assert.Equal(t, entity, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindByID_NotFound", func(t *testing.T) {
		mockRepo.On("FindByID", ctx, "999").Return(TestEntity{}, repository.ErrEntityNotFound).Once()

		_, err := mockRepo.FindByID(ctx, "999")

		assert.Error(t, err)
		assert.True(t, errors.Is(err, repository.ErrEntityNotFound))
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindAll", func(t *testing.T) {
		entities := []TestEntity{
			{ID: "1", Name: "Entity 1"},
			{ID: "2", Name: "Entity 2"},
		}

		opts := []repository.FindOption{
			repository.WithLimit(10),
			repository.WithSort("id", repository.SortAscending),
		}

		mockRepo.On("FindAll", ctx, opts).Return(entities, nil).Once()

		results, err := mockRepo.FindAll(ctx, opts...)

		assert.NoError(t, err)
		assert.Equal(t, entities, results)
		assert.Len(t, results, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Save", func(t *testing.T) {
		mockRepo.On("Save", ctx, entity).Return(nil).Once()

		err := mockRepo.Save(ctx, entity)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		mockRepo.On("Update", ctx, entity).Return(nil).Once()

		err := mockRepo.Update(ctx, entity)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		mockRepo.On("Delete", ctx, "1").Return(nil).Once()

		err := mockRepo.Delete(ctx, "1")

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Count", func(t *testing.T) {
		opts := []repository.FindOption{
			repository.WithComplexFilter("name", "eq", "Test"),
		}

		mockRepo.On("Count", ctx, opts).Return(int64(5), nil).Once()

		count, err := mockRepo.Count(ctx, opts...)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		mockRepo.AssertExpectations(t)
	})
}

// TestRepositoryError 레포지토리 에러 테스트
func TestRepositoryError(t *testing.T) {
	// 기본 에러 생성
	err := repository.NewRepositoryError("FindByID", "User", "not found", repository.ErrEntityNotFound)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrEntityNotFound))
	assert.Contains(t, err.Error(), "FindByID")
	assert.Contains(t, err.Error(), "User")
	assert.Contains(t, err.Error(), "not found")

	// Unwrap 테스트
	unwrapped := errors.Unwrap(err)
	assert.Equal(t, repository.ErrEntityNotFound, unwrapped)

	// 에러 체이닝
	wrappedErr := repository.NewRepositoryError("Query", "Database", "execution failed", err)
	assert.True(t, errors.Is(wrappedErr, repository.ErrEntityNotFound))
}

// TestFindOptions 조회 옵션 테스트
func TestFindOptions(t *testing.T) {
	options := repository.NewFindOptions()

	// 기본값 확인
	assert.Equal(t, 100, options.Limit)
	assert.Equal(t, 0, options.Offset)
	assert.Equal(t, repository.SortAscending, options.SortOrder)

	// WithLimit
	limitOpt := repository.WithLimit(20)
	limitOpt.Apply(options)
	assert.Equal(t, 20, options.Limit)

	// WithOffset
	offsetOpt := repository.WithOffset(10)
	offsetOpt.Apply(options)
	assert.Equal(t, 10, options.Offset)

	// WithSort
	sortOpt := repository.WithSort("created_at", repository.SortDescending)
	sortOpt.Apply(options)
	assert.Equal(t, "created_at", options.SortBy)
	assert.Equal(t, repository.SortDescending, options.SortOrder)

	// WithFilter
	filterOpt := repository.WithComplexFilter("status", "eq", "active")
	filterOpt.Apply(options)
	assert.Len(t, options.FilterList, 1)
	assert.Equal(t, "status", options.FilterList[0].Field)
	assert.Equal(t, "eq", options.FilterList[0].Operator)
	assert.Equal(t, "active", options.FilterList[0].Value)

	// WithPagination
	paginationOpt := repository.WithPagination(2, 25)
	paginationOpt.Apply(options)
	assert.NotNil(t, options.Pagination)
	assert.Equal(t, 2, options.Pagination.Page)
	assert.Equal(t, 25, options.Pagination.PageSize)
}
