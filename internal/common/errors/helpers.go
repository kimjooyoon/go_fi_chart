package errors

import (
	"errors"
)

// As는 표준 errors.As 함수의 래퍼입니다.
// 에러 체인에서 target 타입의 에러를 찾습니다.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// Is는 표준 errors.Is 함수의 래퍼입니다.
// 에러 체인에서 target 에러가 있는지 확인합니다.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// IsErrorType은 주어진 에러가 특정 도메인 에러 타입인지 확인합니다.
func IsErrorType(err error, errorTypeCheck func(error) bool) bool {
	return err != nil && errorTypeCheck(err)
}

// GetDomainError는 에러에서 DomainError 인터페이스를 구현하는 에러를 추출합니다.
// 찾지 못한 경우 nil을 반환합니다.
func GetDomainError(err error) DomainError {
	if err == nil {
		return nil
	}

	var domainErr DomainError
	if errors.As(err, &domainErr) {
		return domainErr
	}

	return nil
}
