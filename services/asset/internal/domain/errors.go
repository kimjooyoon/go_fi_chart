package domain

import (
	"fmt"
	"net/http"

	commonerrors "github.com/aske/go_fi_chart/internal/common/errors"
)

// AssetNotFoundError 자산을 찾을 수 없을 때 반환되는 에러입니다.
type AssetNotFoundError struct {
	AssetID string
}

func (e AssetNotFoundError) Error() string {
	return fmt.Sprintf("자산을 찾을 수 없습니다: %s", e.AssetID)
}

func (e AssetNotFoundError) Code() string {
	return commonerrors.ErrorCodeAssetNotFound
}

func (e AssetNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// NewAssetNotFoundError 자산을 찾을 수 없는 에러를 생성합니다.
func NewAssetNotFoundError(assetID string) error {
	return AssetNotFoundError{AssetID: assetID}
}

// AssetAlreadyExistsError 자산이 이미 존재할 때 반환되는 에러입니다.
type AssetAlreadyExistsError struct {
	AssetID string
}

func (e AssetAlreadyExistsError) Error() string {
	return fmt.Sprintf("자산이 이미 존재합니다: %s", e.AssetID)
}

func (e AssetAlreadyExistsError) Code() string {
	return commonerrors.ErrorCodeAssetAlreadyExists
}

func (e AssetAlreadyExistsError) StatusCode() int {
	return http.StatusConflict
}

// NewAssetAlreadyExistsError 자산이 이미 존재하는 에러를 생성합니다.
func NewAssetAlreadyExistsError(assetID string) error {
	return AssetAlreadyExistsError{AssetID: assetID}
}

// AssetInvalidDataError 자산 데이터가 유효하지 않을 때 반환되는 에러입니다.
type AssetInvalidDataError struct {
	AssetID string
	Message string
}

func (e AssetInvalidDataError) Error() string {
	return fmt.Sprintf("자산 데이터가 유효하지 않습니다. %s: %s", e.AssetID, e.Message)
}

func (e AssetInvalidDataError) Code() string {
	return commonerrors.ErrorCodeAssetInvalidData
}

func (e AssetInvalidDataError) StatusCode() int {
	return http.StatusBadRequest
}

// NewAssetInvalidDataError 자산 데이터가 유효하지 않은 에러를 생성합니다.
func NewAssetInvalidDataError(assetID, message string) error {
	return AssetInvalidDataError{AssetID: assetID, Message: message}
}

// AssetDeletedError 이미 삭제된 에셋에 대한 작업을 시도할 때 발생하는 에러입니다.
type AssetDeletedError struct {
	AssetID string
}

// Error 에러 메시지를 반환합니다.
func (e AssetDeletedError) Error() string {
	return fmt.Sprintf("자산이 이미 삭제되었습니다: %s", e.AssetID)
}

// Code 에러 코드를 반환합니다.
func (e AssetDeletedError) Code() string {
	return "ERROR_ASSET_DELETED"
}

// StatusCode 에러의 HTTP 상태 코드를 반환합니다.
func (e AssetDeletedError) StatusCode() int {
	return http.StatusBadRequest
}

// NewAssetDeletedError는 새로운 AssetDeletedError를 생성합니다.
func NewAssetDeletedError(assetID string) error {
	return AssetDeletedError{AssetID: assetID}
}

// 에러 타입 확인 유틸리티 함수들

// IsAssetNotFound는 주어진 에러가 AssetNotFoundError 타입인지 확인합니다.
func IsAssetNotFound(err error) bool {
	var e AssetNotFoundError
	return commonerrors.As(err, &e)
}

// IsAssetAlreadyExists는 주어진 에러가 AssetAlreadyExistsError 타입인지 확인합니다.
func IsAssetAlreadyExists(err error) bool {
	var e AssetAlreadyExistsError
	return commonerrors.As(err, &e)
}

// IsAssetInvalidData는 주어진 에러가 AssetInvalidDataError 타입인지 확인합니다.
func IsAssetInvalidData(err error) bool {
	var e AssetInvalidDataError
	return commonerrors.As(err, &e)
}

// IsAssetDeleted는 주어진 에러가 AssetDeletedError 타입인지 확인합니다.
func IsAssetDeleted(err error) bool {
	var e AssetDeletedError
	return commonerrors.As(err, &e)
}
