package errors

import (
	"net/http"
)

// DomainError는 도메인 에러의 공통 인터페이스입니다.
type DomainError interface {
	error
	Code() string    // 에러 코드 반환
	StatusCode() int // HTTP 상태 코드 반환
}

// 도메인 에러 타입 상수
const (
	// 일반 에러 코드
	ErrorCodeNotFound      = "ERROR_NOT_FOUND"
	ErrorCodeAlreadyExists = "ERROR_ALREADY_EXISTS"
	ErrorCodeInvalidInput  = "ERROR_INVALID_INPUT"
	ErrorCodeInternalError = "ERROR_INTERNAL"
	ErrorCodeUnauthorized  = "ERROR_UNAUTHORIZED"
	ErrorCodeForbidden     = "ERROR_FORBIDDEN"

	// 포트폴리오 관련 에러 코드
	ErrorCodePortfolioNotFound   = "ERROR_PORTFOLIO_NOT_FOUND"
	ErrorCodeDuplicateAsset      = "ERROR_DUPLICATE_ASSET"
	ErrorCodeTotalWeightExceeded = "ERROR_TOTAL_WEIGHT_EXCEEDED"

	// 자산 관련 에러 코드
	ErrorCodeAssetNotFound      = "ERROR_ASSET_NOT_FOUND"
	ErrorCodeAssetAlreadyExists = "ERROR_ASSET_ALREADY_EXISTS"
	ErrorCodeAssetInvalidData   = "ERROR_ASSET_INVALID_DATA"
)

// ErrorType은 에러 타입을 식별하는 상수입니다.
const (
	ErrorTypeNotFound      = "NOT_FOUND"
	ErrorTypeAlreadyExists = "ALREADY_EXISTS"
	ErrorTypeInvalidData   = "INVALID_DATA"
	ErrorTypeAccessDenied  = "ACCESS_DENIED"
	ErrorTypeInternalError = "INTERNAL_ERROR"
)

// DefaultStatusCode는 에러 타입별 기본 HTTP 상태 코드를 반환합니다.
func DefaultStatusCode(errorType string) int {
	switch errorType {
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeAlreadyExists:
		return http.StatusConflict
	case ErrorTypeInvalidData:
		return http.StatusBadRequest
	case ErrorTypeAccessDenied:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
