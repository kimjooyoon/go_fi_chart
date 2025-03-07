package domain

import (
	"fmt"
	"net/http"

	commonerrors "github.com/aske/go_fi_chart/internal/common/errors"
)

// PortfolioNotFoundError 포트폴리오를 찾을 수 없을 때 발생하는 에러입니다.
type PortfolioNotFoundError struct {
	PortfolioID string
}

// Error 에러 메시지를 반환합니다.
func (e PortfolioNotFoundError) Error() string {
	return fmt.Sprintf("포트폴리오를 찾을 수 없습니다: %s", e.PortfolioID)
}

// Code 에러 코드를 반환합니다.
func (e PortfolioNotFoundError) Code() string {
	return commonerrors.ErrorCodePortfolioNotFound
}

// StatusCode HTTP 상태 코드를 반환합니다.
func (e PortfolioNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// NewPortfolioNotFoundError 새로운 PortfolioNotFoundError를 생성합니다.
func NewPortfolioNotFoundError(portfolioID string) error {
	return PortfolioNotFoundError{PortfolioID: portfolioID}
}

// DuplicateAssetError 이미 존재하는 자산을 추가하려 할 때 발생하는 에러입니다.
type DuplicateAssetError struct {
	PortfolioID string
	AssetID     string
}

// Error 에러 메시지를 반환합니다.
func (e DuplicateAssetError) Error() string {
	return fmt.Sprintf("이미 존재하는 자산입니다: %s (포트폴리오: %s)", e.AssetID, e.PortfolioID)
}

// Code 에러 코드를 반환합니다.
func (e DuplicateAssetError) Code() string {
	return commonerrors.ErrorCodeDuplicateAsset
}

// StatusCode HTTP 상태 코드를 반환합니다.
func (e DuplicateAssetError) StatusCode() int {
	return http.StatusConflict
}

// NewDuplicateAssetError 새로운 DuplicateAssetError를 생성합니다.
func NewDuplicateAssetError(portfolioID, assetID string) error {
	return DuplicateAssetError{
		PortfolioID: portfolioID,
		AssetID:     assetID,
	}
}

// TotalWeightExceededError 총 가중치가 100%를 초과할 때 발생하는 에러입니다.
type TotalWeightExceededError struct {
	PortfolioID string
	TotalWeight float64
}

// Error 에러 메시지를 반환합니다.
func (e TotalWeightExceededError) Error() string {
	return fmt.Sprintf("총 가중치가 100%%를 초과합니다: %.2f%% (포트폴리오: %s)", e.TotalWeight, e.PortfolioID)
}

// Code 에러 코드를 반환합니다.
func (e TotalWeightExceededError) Code() string {
	return commonerrors.ErrorCodeTotalWeightExceeded
}

// StatusCode HTTP 상태 코드를 반환합니다.
func (e TotalWeightExceededError) StatusCode() int {
	return http.StatusBadRequest
}

// NewTotalWeightExceededError 새로운 TotalWeightExceededError를 생성합니다.
func NewTotalWeightExceededError(portfolioID string, totalWeight float64) error {
	return TotalWeightExceededError{
		PortfolioID: portfolioID,
		TotalWeight: totalWeight,
	}
}

// AssetNotFoundError 자산을 찾을 수 없을 때 발생하는 에러입니다.
type AssetNotFoundError struct {
	AssetID string
}

// Error 에러 메시지를 반환합니다.
func (e AssetNotFoundError) Error() string {
	return fmt.Sprintf("자산을 찾을 수 없습니다: %s", e.AssetID)
}

// Code 에러 코드를 반환합니다.
func (e AssetNotFoundError) Code() string {
	return commonerrors.ErrorCodeAssetNotFound
}

// StatusCode HTTP 상태 코드를 반환합니다.
func (e AssetNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// NewAssetNotFoundError 새로운 AssetNotFoundError를 생성합니다.
func NewAssetNotFoundError(assetID string) error {
	return AssetNotFoundError{AssetID: assetID}
}

// 에러 타입 확인 유틸리티 함수들

// IsPortfolioNotFound는 주어진 에러가 PortfolioNotFoundError 타입인지 확인합니다.
func IsPortfolioNotFound(err error) bool {
	var e PortfolioNotFoundError
	return commonerrors.As(err, &e)
}

// IsDuplicateAsset는 주어진 에러가 DuplicateAssetError 타입인지 확인합니다.
func IsDuplicateAsset(err error) bool {
	var e DuplicateAssetError
	return commonerrors.As(err, &e)
}

// IsTotalWeightExceeded는 주어진 에러가 TotalWeightExceededError 타입인지 확인합니다.
func IsTotalWeightExceeded(err error) bool {
	var e TotalWeightExceededError
	return commonerrors.As(err, &e)
}

// IsAssetNotFound는 주어진 에러가 AssetNotFoundError 타입인지 확인합니다.
func IsAssetNotFound(err error) bool {
	var e AssetNotFoundError
	return commonerrors.As(err, &e)
}
