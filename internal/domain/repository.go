package domain

import (
	"context"
	"time"
)

// Entity 모든 엔티티가 구현해야 하는 기본 인터페이스
type Entity interface {
	GetID() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

// SearchCriteria 검색 조건 인터페이스
type SearchCriteria interface {
	ToQuery() (string, []interface{})
}

// Repository 기본 레포지토리 인터페이스
type Repository[T Entity, ID comparable] interface {
	// 기본 CRUD 작업
	Save(ctx context.Context, entity T) error
	FindByID(ctx context.Context, id ID) (T, error)
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id ID) error

	// 검색 작업
	FindAll(ctx context.Context, criteria SearchCriteria) ([]T, error)
	FindOne(ctx context.Context, criteria SearchCriteria) (T, error)

	// 트랜잭션 관리
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

// RepositoryError 레포지토리 관련 에러 타입
type RepositoryError struct {
	Op  string // 작업 종류 (예: Save, FindByID 등)
	Err error  // 원본 에러
}

func (e *RepositoryError) Error() string {
	if e.Err == nil {
		return e.Op
	}
	return e.Op + ": " + e.Err.Error()
}

// NewRepositoryError 새로운 레포지토리 에러 생성
func NewRepositoryError(op string, err error) error {
	return &RepositoryError{
		Op:  op,
		Err: err,
	}
}
