package repository

import (
	"errors"
	"fmt"
)

// 표준 에러 변수
var (
	// ErrEntityNotFound 엔티티를 찾을 수 없을 때 발생하는 에러
	ErrEntityNotFound = errors.New("entity not found")

	// ErrDuplicateEntity 중복된 엔티티가 존재할 때 발생하는 에러
	ErrDuplicateEntity = errors.New("duplicate entity")

	// ErrInvalidEntity 유효하지 않은 엔티티일 때 발생하는 에러
	ErrInvalidEntity = errors.New("invalid entity")

	// ErrRepositoryError 레포지토리 내부 에러
	ErrRepositoryError = errors.New("repository error")

	// ErrTransactionFailed 트랜잭션 실패 에러
	ErrTransactionFailed = errors.New("transaction failed")
)

// Error 레포지토리 관련 상세 에러 타입
type Error struct {
	// Operation 수행하던 작업
	Operation string

	// Entity 관련 엔티티 타입 또는 ID
	Entity string

	// Message 추가 에러 메시지
	Message string

	// Err 원본 에러
	Err error
}

// Error RepositoryError의 문자열 표현 반환
func (e *Error) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("%s: %s %s", e.Operation, e.Entity, e.Message)
	}
	return fmt.Sprintf("%s: %s %s: %v", e.Operation, e.Entity, e.Message, e.Err)
}

// Unwrap 원본 에러 반환
func (e *Error) Unwrap() error {
	return e.Err
}

// NewRepositoryError 새로운 레포지토리 에러 생성
func NewRepositoryError(operation, entity, message string, err error) *Error {
	return &Error{
		Operation: operation,
		Entity:    entity,
		Message:   message,
		Err:       err,
	}
}

// Is 에러 비교 지원
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}

	// 직접 원본 에러와 비교
	if errors.Is(e.Err, target) {
		return true
	}

	// 기본 에러 타입과 비교
	switch target {
	case ErrEntityNotFound, ErrDuplicateEntity, ErrInvalidEntity, ErrRepositoryError, ErrTransactionFailed:
		return errors.Is(e.Err, target)
	default:
		return false
	}
}
