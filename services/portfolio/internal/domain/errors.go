package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrPortfolioNotFound 포트폴리오를 찾을 수 없을 때 발생하는 에러입니다.
	ErrPortfolioNotFound = errors.New("portfolio not found")

	// ErrDuplicateAsset 이미 존재하는 자산을 추가하려 할 때 발생하는 에러입니다.
	ErrDuplicateAsset = errors.New("duplicate asset")

	// ErrTotalWeightExceeded 총 가중치가 100%를 초과할 때 발생하는 에러입니다.
	ErrTotalWeightExceeded = errors.New("total weight exceeded")
)

// AssetNotFoundError 자산을 찾을 수 없을 때 발생하는 에러입니다.
type AssetNotFoundError struct {
	AssetID string
}

// Error 에러 메시지를 반환합니다.
func (e AssetNotFoundError) Error() string {
	return fmt.Sprintf("asset not found: %s", e.AssetID)
}

// NewAssetNotFoundError 새로운 AssetNotFoundError를 생성합니다.
func NewAssetNotFoundError(assetID string) error {
	return AssetNotFoundError{AssetID: assetID}
}
