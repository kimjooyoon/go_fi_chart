package repository

import (
	"context"
)

// Repository 제네릭 기반 기본 레포지토리 인터페이스
// T: 엔티티 타입, ID: 식별자 타입 (comparable 제약조건 적용)
type Repository[T any, ID comparable] interface {
	// 기본 CRUD 작업
	FindByID(ctx context.Context, id ID) (T, error)
	FindAll(ctx context.Context, opts ...FindOption) ([]T, error)
	Save(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id ID) error

	// 추가 조회 기능
	Count(ctx context.Context, opts ...FindOption) (int64, error)
}

// ReadRepository 읽기 전용 레포지토리 인터페이스
type ReadRepository[T any, ID comparable] interface {
	FindByID(ctx context.Context, id ID) (T, error)
	FindAll(ctx context.Context, opts ...FindOption) ([]T, error)
	Count(ctx context.Context, opts ...FindOption) (int64, error)
}

// WriteRepository 쓰기 전용 레포지토리 인터페이스
type WriteRepository[T any, ID comparable] interface {
	Save(ctx context.Context, entity T) error
	Update(ctx context.Context, entity T) error
	Delete(ctx context.Context, id ID) error
}

// BatchRepository 배치 작업을 지원하는 레포지토리 인터페이스
type BatchRepository[T any] interface {
	SaveAll(ctx context.Context, entities []T) error
	UpdateAll(ctx context.Context, entities []T) error
}

// TransactionalRepository 트랜잭션을 지원하는 레포지토리 인터페이스
type TransactionalRepository interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
