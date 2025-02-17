package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetRepository는 테스트용 자산 저장소입니다.
type MockAssetRepository struct {
	mock.Mock
}

func (m *MockAssetRepository) Save(ctx context.Context, asset *domain.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockAssetRepository) FindByID(ctx context.Context, id string) (*domain.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Asset), args.Error(1)
}

func (m *MockAssetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockAssetRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockAssetRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Asset, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Asset), args.Error(1)
}

func (m *MockAssetRepository) FindByType(ctx context.Context, assetType domain.AssetType) ([]*domain.Asset, error) {
	args := m.Called(ctx, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Asset), args.Error(1)
}

// setupTestHandler는 테스트용 핸들러를 생성합니다.
func setupTestHandler() (*Handler, *MockAssetRepository) {
	repo := new(MockAssetRepository)
	handler := NewHandler(repo)
	return handler, repo
}

// setupTestRouter는 테스트용 라우터를 생성합니다.
func setupTestRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	h.RegisterRoutes(r)
	return r
}

// createValidAsset는 테스트용 유효한 자산을 생성합니다.
func createValidAsset() *domain.Asset {
	amount, _ := valueobjects.NewMoney(1000.0, "USD")
	return domain.NewAsset("user-123", domain.Stock, "Test Asset", amount)
}

func TestCreateAsset(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("유효한 자산 생성", func(t *testing.T) {
		req := CreateAssetRequest{
			UserID:   uuid.New().String(),
			Type:     "STOCK",
			Name:     "테스트 자산",
			Amount:   1000.0,
			Currency: "USD",
		}

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Asset")).Return(nil)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/api/v1/assets", bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("잘못된 요청 데이터", func(t *testing.T) {
		req := CreateAssetRequest{
			UserID: "invalid-uuid",
		}

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/api/v1/assets", bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAsset(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	asset := createValidAsset()

	t.Run("존재하는 자산 조회", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, asset.ID).Return(asset, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/api/v1/assets/"+asset.ID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 자산 조회", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		repo.On("FindByID", mock.Anything, nonExistentID).Return(nil, domain.ErrAssetNotFound)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/api/v1/assets/"+nonExistentID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestUpdateAsset(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	asset := createValidAsset()

	t.Run("자산 업데이트 성공", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, asset.ID).Return(asset, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Asset")).Return(nil)

		updateReq := UpdateAssetRequest{
			Name:     "업데이트된 자산",
			Amount:   2000.0,
			Currency: "USD",
		}

		body, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/api/v1/assets/"+asset.ID, bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 자산 업데이트", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		repo.On("FindByID", mock.Anything, nonExistentID).Return(nil, domain.ErrAssetNotFound)

		updateReq := UpdateAssetRequest{
			Name: "업데이트된 자산",
		}

		body, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/api/v1/assets/"+nonExistentID, bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestDeleteAsset(t *testing.T) {
	t.Run("자산 삭제 성공", func(t *testing.T) {
		// Given
		repo := new(MockAssetRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		assetID := "e6ac8e5e-61c5-4466-bbab-b9e5b096cd37"
		asset := createValidAsset()
		asset.ID = assetID

		repo.On("FindByID", mock.Anything, assetID).Return(asset, nil)
		repo.On("Delete", mock.Anything, assetID).Return(nil)

		// When
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/"+assetID, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusNoContent, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 자산 삭제", func(t *testing.T) {
		// Given
		repo := new(MockAssetRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		assetID := "non-existent-id"
		repo.On("FindByID", mock.Anything, assetID).Return(nil, domain.ErrAssetNotFound)

		// When
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/"+assetID, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusNotFound, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("자산 조회 실패", func(t *testing.T) {
		// Given
		repo := new(MockAssetRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		assetID := "e6ac8e5e-61c5-4466-bbab-b9e5b096cd37"
		repo.On("FindByID", mock.Anything, assetID).Return(nil, errors.New("database error"))

		// When
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/"+assetID, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("자산 삭제 실패", func(t *testing.T) {
		// Given
		repo := new(MockAssetRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		assetID := "e6ac8e5e-61c5-4466-bbab-b9e5b096cd37"
		asset := createValidAsset()
		asset.ID = assetID

		repo.On("FindByID", mock.Anything, assetID).Return(asset, nil)
		repo.On("Delete", mock.Anything, assetID).Return(errors.New("database error"))

		// When
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/assets/"+assetID, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		repo.AssertExpectations(t)
	})
}

func TestListAssets(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	userID := uuid.New().String()
	assets := []*domain.Asset{
		createValidAsset(),
		createValidAsset(),
	}

	t.Run("사용자의 자산 목록 조회", func(t *testing.T) {
		repo.On("FindByUserID", mock.Anything, userID).Return(assets, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/api/v1/assets?userId="+userID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp []AssetResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, len(assets))
		repo.AssertExpectations(t)
	})
}

func TestListAssetsByType(t *testing.T) {
	t.Run("자산 유형별 목록 조회", func(t *testing.T) {
		// Given
		mockRepo := new(MockAssetRepository)
		handler := NewHandler(mockRepo)
		router := setupTestRouter(handler)

		assets := []*domain.Asset{
			createValidAsset(),
			createValidAsset(),
		}

		mockRepo.On("FindByType", mock.Anything, domain.Stock).Return(assets, nil)

		// When
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/assets/types/STOCK", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)

		var response []AssetResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)
		mockRepo.AssertExpectations(t)
	})

	t.Run("잘못된 자산 유형", func(t *testing.T) {
		// Given
		mockRepo := new(MockAssetRepository)
		handler := NewHandler(mockRepo)
		router := setupTestRouter(handler)

		// When
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/assets/types/INVALID", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}

func TestCreateAsset_Success(t *testing.T) {
	// Given
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Asset")).Return(nil)

	reqBody := CreateAssetRequest{
		UserID:   "user-123",
		Type:     "STOCK",
		Name:     "Test Asset",
		Amount:   1000.0,
		Currency: "USD",
	}
	body, _ := json.Marshal(reqBody)

	// When
	req := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusCreated, res.Code)

	var response AssetResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, reqBody.UserID, response.UserID)
	assert.Equal(t, reqBody.Type, response.Type)
	assert.Equal(t, reqBody.Name, response.Name)
	assert.Equal(t, reqBody.Amount, response.Amount)
	assert.Equal(t, reqBody.Currency, response.Currency)
}

func TestCreateAsset_InvalidRequest(t *testing.T) {
	// Given
	handler, _ := setupTestHandler()
	router := setupTestRouter(handler)

	invalidReqBody := CreateAssetRequest{
		UserID:   "", // 필수 필드 누락
		Type:     "INVALID_TYPE",
		Name:     "Test Asset",
		Amount:   -1000.0, // 잘못된 금액
		Currency: "",      // 필수 필드 누락
	}
	body, _ := json.Marshal(invalidReqBody)

	// When
	req := httptest.NewRequest(http.MethodPost, "/api/v1/assets", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestGetAsset_Success(t *testing.T) {
	// Given
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	asset := createValidAsset()
	repo.On("FindByID", mock.Anything, asset.ID).Return(asset, nil)

	// When
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/assets/%s", asset.ID), nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusOK, res.Code)

	var response AssetResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, asset.ID, response.ID)
}

func TestGetAsset_NotFound(t *testing.T) {
	// Given
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	assetID := "non-existent-id"
	repo.On("FindByID", mock.Anything, assetID).Return(nil, domain.ErrAssetNotFound)

	// When
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets/non-existent-id", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusNotFound, res.Code)
}

func TestUpdateAsset_Success(t *testing.T) {
	// Given
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	asset := createValidAsset()
	repo.On("FindByID", mock.Anything, asset.ID).Return(asset, nil)
	repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Asset")).Return(nil)

	updateReq := UpdateAssetRequest{
		Name:     "Updated Asset",
		Amount:   2000.0,
		Currency: "USD",
	}
	body, _ := json.Marshal(updateReq)

	// When
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/assets/%s", asset.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusOK, res.Code)

	var response AssetResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, updateReq.Name, response.Name)
	assert.Equal(t, updateReq.Amount, response.Amount)
	assert.Equal(t, updateReq.Currency, response.Currency)
}

func TestListAssets_Success(t *testing.T) {
	// Given
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	assets := []*domain.Asset{
		createValidAsset(),
		createValidAsset(),
	}
	repo.On("FindByUserID", mock.Anything, "user-123").Return(assets, nil)

	// When
	req := httptest.NewRequest(http.MethodGet, "/api/v1/assets?userId=user-123", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	// Then
	assert.Equal(t, http.StatusOK, res.Code)

	var response []AssetResponse
	err := json.Unmarshal(res.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, len(assets))
}
