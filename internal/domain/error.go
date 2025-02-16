package domain

import "fmt"

// Error 도메인 에러 인터페이스입니다.
type Error interface {
	error
	// Code 에러 코드를 반환합니다.
	Code() string
	// Domain 에러가 발생한 도메인을 반환합니다.
	Domain() string
}

// BaseError 기본 도메인 에러 구현체입니다.
type BaseError struct {
	domain string
	code   string
	msg    string
}

// NewError 새로운 도메인 에러를 생성합니다.
func NewError(domain string, code string, msg string) Error {
	return &BaseError{
		domain: domain,
		code:   code,
		msg:    msg,
	}
}

func (e *BaseError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.domain, e.code, e.msg)
}

func (e *BaseError) Code() string {
	return e.code
}

func (e *BaseError) Domain() string {
	return e.domain
}

// 자주 사용되는 에러 코드 상수
const (
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeAlreadyExists    = "ALREADY_EXISTS"
	ErrCodeInvalidArgument  = "INVALID_ARGUMENT"
	ErrCodeInvalidOperation = "INVALID_OPERATION"
	ErrCodeNotImplemented   = "NOT_IMPLEMENTED"
	ErrCodeInternal         = "INTERNAL"
)
