package source

import (
	"fmt"
	"time"
)

// SourceError는 데이터 소스에서 발생하는 에러의 기본 인터페이스입니다.
type SourceError interface {
	error
	SourceName() string // 데이터 소스 이름
	ErrorCode() string  // 에러 코드
}

// BaseSourceError는 데이터 소스 에러의 기본 구현입니다.
type BaseSourceError struct {
	Source string // 데이터 소스 이름
	Code   string // 에러 코드
	Msg    string // 에러 메시지
}

// Error는 에러 메시지를 반환합니다.
func (e *BaseSourceError) Error() string {
	return fmt.Sprintf("[%s] %s: %s", e.Source, e.Code, e.Msg)
}

// SourceName은 데이터 소스 이름을 반환합니다.
func (e *BaseSourceError) SourceName() string {
	return e.Source
}

// ErrorCode는 에러 코드를 반환합니다.
func (e *BaseSourceError) ErrorCode() string {
	return e.Code
}

// NewSourceError는 새로운 SourceError를 생성합니다.
func NewSourceError(source, code, message string) SourceError {
	return &BaseSourceError{
		Source: source,
		Code:   code,
		Msg:    message,
	}
}

// RateLimitError는 API 속도 제한 관련 에러입니다.
type RateLimitError struct {
	BaseSourceError
	RetryAfter time.Duration // 재시도 가능 시간
}

// NewRateLimitError는 새로운 RateLimitError를 생성합니다.
func NewRateLimitError(source string, retryAfter time.Duration) *RateLimitError {
	return &RateLimitError{
		BaseSourceError: BaseSourceError{
			Source: source,
			Code:   "RATE_LIMIT_EXCEEDED",
			Msg:    fmt.Sprintf("Rate limit exceeded, retry after %v", retryAfter),
		},
		RetryAfter: retryAfter,
	}
}

// GetRetryAfter는 재시도 가능 시간을 반환합니다.
func (e *RateLimitError) GetRetryAfter() time.Duration {
	return e.RetryAfter
}

// NetworkError는 네트워크 관련 에러입니다.
type NetworkError struct {
	BaseSourceError
	Retryable bool // 재시도 가능 여부
}

// NewNetworkError는 새로운 NetworkError를 생성합니다.
func NewNetworkError(source, message string, retryable bool) *NetworkError {
	return &NetworkError{
		BaseSourceError: BaseSourceError{
			Source: source,
			Code:   "NETWORK_ERROR",
			Msg:    message,
		},
		Retryable: retryable,
	}
}

// IsRetryable은 에러가 재시도 가능한지 여부를 반환합니다.
func (e *NetworkError) IsRetryable() bool {
	return e.Retryable
}

// ParseError는 데이터 파싱 관련 에러입니다.
type ParseError struct {
	BaseSourceError
	RawData string // 원시 데이터
}

// NewParseError는 새로운 ParseError를 생성합니다.
func NewParseError(source, message, rawData string) *ParseError {
	return &ParseError{
		BaseSourceError: BaseSourceError{
			Source: source,
			Code:   "PARSE_ERROR",
			Msg:    message,
		},
		RawData: rawData,
	}
}

// GetRawData는 파싱에 실패한 원시 데이터를 반환합니다.
func (e *ParseError) GetRawData() string {
	return e.RawData
}
