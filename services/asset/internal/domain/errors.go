package domain

import "errors"

var (
	// ErrAssetNotFound 자산을 찾을 수 없을 때 발생하는 에러입니다.
	ErrAssetNotFound = errors.New("asset not found")

	// ErrInvalidAssetType 잘못된 자산 유형일 때 발생하는 에러입니다.
	ErrInvalidAssetType = errors.New("invalid asset type")

	// ErrInvalidAmount 잘못된 금액일 때 발생하는 에러입니다.
	ErrInvalidAmount = errors.New("invalid amount")

	// ErrInvalidCurrency 잘못된 통화일 때 발생하는 에러입니다.
	ErrInvalidCurrency = errors.New("invalid currency")

	// ErrAssetDeleted 삭제된 자산에 대한 작업을 시도할 때 발생하는 에러입니다.
	ErrAssetDeleted = errors.New("삭제된 자산입니다")
)
