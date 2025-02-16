package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/internal/domain"
	"github.com/aske/go_fi_chart/internal/domain/asset"
	"github.com/aske/go_fi_chart/internal/domain/gamification"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 모의 리포지토리 구현
type mockAssetRepository struct {
	mock.Mock
}

func (m *mockAssetRepository) Save(ctx context.Context, asset *asset.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *mockAssetRepository) FindByID(ctx context.Context, id string) (*asset.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*asset.Asset), args.Error(1)
}

func (m *mockAssetRepository) FindByUserID(ctx context.Context, userID string) ([]*asset.Asset, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*asset.Asset), args.Error(1)
}

func (m *mockAssetRepository) Update(ctx context.Context, asset *asset.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *mockAssetRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockAssetRepository) FindByType(ctx context.Context, assetType asset.Type) ([]*asset.Asset, error) {
	args := m.Called(ctx, assetType)
	return args.Get(0).([]*asset.Asset), args.Error(1)
}

func (m *mockAssetRepository) UpdateAmount(ctx context.Context, id string, amount asset.Money) error {
	args := m.Called(ctx, id, amount)
	return args.Error(0)
}

func (m *mockAssetRepository) FindAll(ctx context.Context, criteria domain.SearchCriteria) ([]*asset.Asset, error) {
	args := m.Called(ctx, criteria)
	return args.Get(0).([]*asset.Asset), args.Error(1)
}

func (m *mockAssetRepository) FindOne(ctx context.Context, criteria domain.SearchCriteria) (*asset.Asset, error) {
	args := m.Called(ctx, criteria)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*asset.Asset), args.Error(1)
}

func (m *mockAssetRepository) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// 테스트 헬퍼 함수
func setupTestHandler() (*Handler, *mockAssetRepository) {
	mockRepo := new(mockAssetRepository)
	mockTransactionRepo := new(asset.MockTransactionRepository)
	mockPortfolioRepo := new(asset.MockPortfolioRepository)
	mockGamificationRepo := new(gamification.MockRepository)

	handler := NewHandler(mockRepo, mockTransactionRepo, mockPortfolioRepo, mockGamificationRepo)
	return handler, mockRepo
}

func TestListAssets(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 테스트 데이터 준비
	testUserID := "test-user"
	testAssets := []*asset.Asset{
		{
			ID:     "asset-1",
			UserID: testUserID,
			Type:   asset.Stock,
			Name:   "삼성전자",
			Amount: asset.Money{
				Amount:   1000000,
				Currency: "KRW",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// 모의 동작 설정
	mockRepo.On("FindByUserID", mock.Anything, testUserID).Return(testAssets, nil)

	// 테스트 요청 생성
	req := httptest.NewRequest("GET", "/assets?userId="+testUserID, nil)
	w := httptest.NewRecorder()

	// 핸들러 실행
	handler.ListAssets(w, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, w.Code)

	var response []AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, testAssets[0].ID, response[0].ID)
	assert.Equal(t, testAssets[0].UserID, response[0].UserID)
	assert.Equal(t, string(testAssets[0].Type), response[0].Type)
	assert.Equal(t, testAssets[0].Name, response[0].Name)
	assert.Equal(t, testAssets[0].Amount.Amount, response[0].Amount)
	assert.Equal(t, testAssets[0].Amount.Currency, response[0].Currency)

	// 모의 객체 호출 검증
	mockRepo.AssertExpectations(t)
}

func TestCreateAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 테스트 데이터 준비
	createReq := CreateAssetRequest{
		UserID:   "test-user",
		Type:     string(asset.Stock),
		Name:     "삼성전자",
		Amount:   1000000,
		Currency: "KRW",
	}

	// 모의 동작 설정
	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)

	// 테스트 요청 생성
	reqBody, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/assets", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	// 핸들러 실행
	handler.CreateAsset(w, req)

	// 응답 검증
	assert.Equal(t, http.StatusCreated, w.Code)

	var response AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, createReq.UserID, response.UserID)
	assert.Equal(t, createReq.Type, response.Type)
	assert.Equal(t, createReq.Name, response.Name)
	assert.Equal(t, createReq.Amount, response.Amount)
	assert.Equal(t, createReq.Currency, response.Currency)

	// 모의 객체 호출 검증
	mockRepo.AssertExpectations(t)
}

func TestGetAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 테스트 데이터 준비
	testAsset := &asset.Asset{
		ID:     "test-asset",
		UserID: "test-user",
		Type:   asset.Stock,
		Name:   "삼성전자",
		Amount: asset.Money{
			Amount:   1000000,
			Currency: "KRW",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 모의 동작 설정
	mockRepo.On("FindByID", mock.Anything, testAsset.ID).Return(testAsset, nil)

	// 테스트 요청 생성
	req := httptest.NewRequest("GET", "/assets/"+testAsset.ID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", testAsset.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// 핸들러 실행
	handler.GetAsset(w, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, w.Code)

	var response AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, testAsset.ID, response.ID)
	assert.Equal(t, testAsset.UserID, response.UserID)
	assert.Equal(t, string(testAsset.Type), response.Type)
	assert.Equal(t, testAsset.Name, response.Name)
	assert.Equal(t, testAsset.Amount.Amount, response.Amount)
	assert.Equal(t, testAsset.Amount.Currency, response.Currency)

	// 모의 객체 호출 검증
	mockRepo.AssertExpectations(t)
}

func TestUpdateAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 테스트 데이터 준비
	testAsset := &asset.Asset{
		ID:     "test-asset",
		UserID: "test-user",
		Type:   asset.Stock,
		Name:   "삼성전자",
		Amount: asset.Money{
			Amount:   1000000,
			Currency: "KRW",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateReq := UpdateAssetRequest{
		Name:     "삼성전자 우선주",
		Amount:   2000000,
		Currency: "KRW",
	}

	// 모의 동작 설정
	mockRepo.On("FindByID", mock.Anything, testAsset.ID).Return(testAsset, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)

	// 테스트 요청 생성
	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/assets/"+testAsset.ID, bytes.NewBuffer(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", testAsset.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// 핸들러 실행
	handler.UpdateAsset(w, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, w.Code)

	var response AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, testAsset.ID, response.ID)
	assert.Equal(t, updateReq.Name, response.Name)
	assert.Equal(t, updateReq.Amount, response.Amount)
	assert.Equal(t, updateReq.Currency, response.Currency)

	// 모의 객체 호출 검증
	mockRepo.AssertExpectations(t)
}

func TestDeleteAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 테스트 데이터 준비
	testAsset := &asset.Asset{
		ID:     "test-asset",
		UserID: "test-user",
		Type:   asset.Stock,
		Name:   "삼성전자",
		Amount: asset.Money{
			Amount:   1000000,
			Currency: "KRW",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 모의 동작 설정
	mockRepo.On("FindByID", mock.Anything, testAsset.ID).Return(testAsset, nil)
	mockRepo.On("Delete", mock.Anything, testAsset.ID).Return(nil)

	// 테스트 요청 생성
	req := httptest.NewRequest("DELETE", "/assets/"+testAsset.ID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", testAsset.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	// 핸들러 실행
	handler.DeleteAsset(w, req)

	// 응답 검증
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 모의 객체 호출 검증
	mockRepo.AssertExpectations(t)
}

// 에러 케이스 테스트 추가
func TestListAssets_InvalidRequest(t *testing.T) {
	handler, _ := setupTestHandler()

	// userId 없이 요청
	req := httptest.NewRequest("GET", "/assets", nil)
	w := httptest.NewRecorder()

	handler.ListAssets(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRequest, response.Code)
}

func TestCreateAsset_InvalidRequest(t *testing.T) {
	handler, _ := setupTestHandler()

	// 잘못된 JSON 요청
	req := httptest.NewRequest("POST", "/assets", bytes.NewBufferString("invalid json"))
	w := httptest.NewRecorder()

	handler.CreateAsset(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, ErrInvalidRequest, response.Code)
}

func TestGetAsset_NotFound(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 존재하지 않는 자산 ID
	assetID := "non-existent"
	mockRepo.On("FindByID", mock.Anything, assetID).Return(nil, fmt.Errorf("자산을 찾을 수 없습니다"))

	req := httptest.NewRequest("GET", "/assets/"+assetID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", assetID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetAsset(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, ErrNotFound, response.Code)

	mockRepo.AssertExpectations(t)
}

func TestUpdateAsset_NotFound(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 존재하지 않는 자산 ID
	assetID := "non-existent"
	mockRepo.On("FindByID", mock.Anything, assetID).Return(nil, fmt.Errorf("자산을 찾을 수 없습니다"))

	updateReq := UpdateAssetRequest{
		Name:     "삼성전자 우선주",
		Amount:   2000000,
		Currency: "KRW",
	}

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PUT", "/assets/"+assetID, bytes.NewBuffer(reqBody))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", assetID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateAsset(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, ErrNotFound, response.Code)

	mockRepo.AssertExpectations(t)
}

func TestDeleteAsset_NotFound(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	// 존재하지 않는 자산 ID
	assetID := "non-existent"
	mockRepo.On("FindByID", mock.Anything, assetID).Return(nil, fmt.Errorf("자산을 찾을 수 없습니다"))

	req := httptest.NewRequest("DELETE", "/assets/"+assetID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", assetID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteAsset(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, ErrNotFound, response.Code)

	mockRepo.AssertExpectations(t)
}
