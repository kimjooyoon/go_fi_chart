package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/portfolio/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPortfolioRepository는 테스트용 포트폴리오 저장소입니다.
type MockPortfolioRepository struct {
	mock.Mock
}

func (m *MockPortfolioRepository) Save(ctx context.Context, portfolio *domain.Portfolio) error {
	args := m.Called(ctx, portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) FindByID(ctx context.Context, id string) (*domain.Portfolio, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Portfolio), args.Error(1)
}

func (m *MockPortfolioRepository) Update(ctx context.Context, portfolio *domain.Portfolio) error {
	args := m.Called(ctx, portfolio)
	return args.Error(0)
}

func (m *MockPortfolioRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockPortfolioRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Portfolio, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Portfolio), args.Error(1)
}

// UpdatePortfolioRequest는 포트폴리오 업데이트 요청 구조체입니다.
type UpdatePortfolioRequest struct {
	Name string `json:"name"`
}

// CreatePortfolioRequest는 포트폴리오 생성 요청 구조체입니다.
type CreatePortfolioRequest struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

// PortfolioResponse는 포트폴리오 응답 구조체입니다.
type PortfolioResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Assets    []Asset   `json:"assets"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Asset은 포트폴리오 자산 응답 구조체입니다.
type Asset struct {
	ID     string  `json:"id"`
	Weight float64 `json:"weight"`
}

// AddAssetRequest는 자산 추가 요청 구조체입니다.
type AddAssetRequest struct {
	AssetID string  `json:"assetId"`
	Weight  float64 `json:"weight"`
}

// UpdateAssetWeightRequest는 자산 가중치 업데이트 요청 구조체입니다.
type UpdateAssetWeightRequest struct {
	Weight float64 `json:"weight"`
}

// setupTestHandler는 테스트용 핸들러를 생성합니다.
func setupTestHandler() (*Handler, *MockPortfolioRepository) {
	repo := new(MockPortfolioRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	handler := NewHandler(repo, logger)
	return handler, repo
}

// setupTestRouter는 테스트용 라우터를 생성합니다.
func setupTestRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r
}

// createValidPortfolio는 테스트용 유효한 포트폴리오를 생성합니다.
func createValidPortfolio() *domain.Portfolio {
	portfolio := domain.NewPortfolio("user-123", "Test Portfolio")
	portfolio.ID = "portfolio-123"
	weight, _ := valueobjects.NewPercentage(0.5)
	portfolio.AddAsset("asset-123", weight)
	return portfolio
}

func TestCreatePortfolio(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("유효한 포트폴리오 생성", func(t *testing.T) {
		req := CreatePortfolioRequest{
			UserID: uuid.New().String(),
			Name:   "테스트 포트폴리오",
		}

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Portfolio")).Return(nil)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/portfolios", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusCreated, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("잘못된 요청 데이터", func(t *testing.T) {
		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/portfolios", bytes.NewReader([]byte("invalid json")))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetPortfolio(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("존재하는 포트폴리오 조회", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, portfolio.ID).Return(portfolio, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/portfolios/"+portfolio.ID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오 조회", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		repo.On("FindByID", mock.Anything, nonExistentID).Return(nil, domain.ErrAssetNotFound)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/portfolios/"+nonExistentID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestUpdatePortfolio(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("포트폴리오 업데이트 성공", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, "portfolio-123").Return(portfolio, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Portfolio")).Return(nil)

		updateReq := UpdatePortfolioRequest{
			Name: "업데이트된 포트폴리오",
		}

		body, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/portfolios/portfolio-123", bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오 업데이트", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, "not-found").Return(nil, domain.ErrPortfolioNotFound)

		updateReq := UpdatePortfolioRequest{
			Name: "업데이트된 포트폴리오",
		}

		body, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/portfolios/not-found", bytes.NewReader(body))
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestDeletePortfolio(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("포트폴리오 삭제 성공", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, portfolio.ID).Return(portfolio, nil)
		repo.On("Delete", mock.Anything, portfolio.ID).Return(nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("DELETE", "/portfolios/"+portfolio.ID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오 삭제", func(t *testing.T) {
		nonExistentID := uuid.New().String()
		repo.On("FindByID", mock.Anything, nonExistentID).Return(nil, domain.ErrPortfolioNotFound)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("DELETE", "/portfolios/"+nonExistentID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestListUserPortfolios(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("사용자의 포트폴리오 목록 조회", func(t *testing.T) {
		userID := uuid.New().String()
		portfolios := []*domain.Portfolio{createValidPortfolio(), createValidPortfolio()}

		repo.On("FindByUserID", mock.Anything, userID).Return(portfolios, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/users/"+userID+"/portfolios", nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []portfolioResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		repo.AssertExpectations(t)
	})
}

func TestGetPortfolioByUserID(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("사용자의 포트폴리오 조회", func(t *testing.T) {
		portfolios := []*domain.Portfolio{portfolio}
		repo.On("FindByUserID", mock.Anything, portfolio.UserID).Return(portfolios, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/users/"+portfolio.UserID+"/portfolios", nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("포트폴리오가 없는 사용자 조회", func(t *testing.T) {
		nonExistentUserID := uuid.New().String()
		repo.On("FindByUserID", mock.Anything, nonExistentUserID).Return([]*domain.Portfolio{}, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/users/"+nonExistentUserID+"/portfolios", nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestListPortfolios(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)

	t.Run("사용자의 포트폴리오 목록 조회", func(t *testing.T) {
		userID := uuid.New().String()
		portfolios := []*domain.Portfolio{createValidPortfolio(), createValidPortfolio()}

		repo.On("FindByUserID", mock.Anything, userID).Return(portfolios, nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("GET", "/portfolios?userId="+userID, nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)

		var response []portfolioResponse
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		repo.AssertExpectations(t)
	})
}

func TestAddAsset(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("자산 추가 성공", func(t *testing.T) {
		req := AddAssetRequest{
			AssetID: uuid.New().String(),
			Weight:  0.5,
		}

		repo.On("FindByID", mock.Anything, portfolio.ID).Return(portfolio, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Portfolio")).Return(nil)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/portfolios/"+portfolio.ID+"/assets", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오에 자산 추가", func(t *testing.T) {
		req := AddAssetRequest{
			AssetID: uuid.New().String(),
			Weight:  0.5,
		}

		repo.On("FindByID", mock.Anything, "not-found").Return(nil, domain.ErrPortfolioNotFound)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("POST", "/portfolios/not-found/assets", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestUpdateAssetWeight(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("자산 가중치 업데이트 성공", func(t *testing.T) {
		req := UpdateAssetWeightRequest{
			Weight: 0.7,
		}

		repo.On("FindByID", mock.Anything, portfolio.ID).Return(portfolio, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Portfolio")).Return(nil)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/portfolios/"+portfolio.ID+"/assets/asset-123", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusOK, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오의 자산 가중치 업데이트", func(t *testing.T) {
		req := UpdateAssetWeightRequest{
			Weight: 0.7,
		}

		repo.On("FindByID", mock.Anything, "not-found").Return(nil, domain.ErrPortfolioNotFound)

		body, err := json.Marshal(req)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("PUT", "/portfolios/not-found/assets/asset-123", bytes.NewReader(body))
		httpReq.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}

func TestRemoveAsset(t *testing.T) {
	handler, repo := setupTestHandler()
	router := setupTestRouter(handler)
	portfolio := createValidPortfolio()

	t.Run("자산 제거 성공", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, portfolio.ID).Return(portfolio, nil)
		repo.On("Update", mock.Anything, mock.AnythingOfType("*domain.Portfolio")).Return(nil)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("DELETE", "/portfolios/"+portfolio.ID+"/assets/asset-123", nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNoContent, w.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 포트폴리오의 자산 제거", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, "not-found").Return(nil, domain.ErrPortfolioNotFound)

		w := httptest.NewRecorder()
		httpReq := httptest.NewRequest("DELETE", "/portfolios/not-found/assets/asset-123", nil)
		router.ServeHTTP(w, httpReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		repo.AssertExpectations(t)
	})
}
