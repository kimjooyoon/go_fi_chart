package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/events"
	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/transaction/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTransactionRepository는 테스트를 위한 mock repository입니다.
type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Save(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByPortfolioID(ctx context.Context, portfolioID uuid.UUID) ([]*domain.Transaction, error) {
	args := m.Called(ctx, portfolioID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByAssetID(ctx context.Context, assetID uuid.UUID) ([]*domain.Transaction, error) {
	args := m.Called(ctx, assetID)
	return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *domain.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// setupTestHandler는 테스트용 핸들러를 생성합니다.
func setupTestHandler() *Handler {
	eventBus := events.NewSimplePublisher()
	repository := domain.NewMemoryTransactionRepository(eventBus)
	return NewHandler(repository)
}

// setupTestRouter는 테스트용 라우터를 생성합니다.
func setupTestRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/transactions", h.CreateTransaction)
		r.Get("/transactions/{id}", h.GetTransaction)
		r.Get("/users/{userID}/transactions", h.ListUserTransactions)
		r.Get("/portfolios/{portfolioID}/transactions", h.ListPortfolioTransactions)
		r.Get("/assets/{assetID}/transactions", h.ListAssetTransactions)
		r.Put("/transactions/{id}", h.UpdateTransaction)
		r.Delete("/transactions/{id}", h.DeleteTransaction)
	})
	return r
}

// createValidTransaction은 테스트용 유효한 거래를 생성합니다.
func createValidTransaction(t *testing.T) *domain.Transaction {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, err := valueobjects.NewMoney(100.0, "USD")
	assert.NoError(t, err)
	executedPrice, err := valueobjects.NewMoney(50.0, "USD")
	assert.NoError(t, err)
	executedAt := time.Now()

	transaction, err := domain.NewTransaction(
		userID,
		portfolioID,
		assetID,
		domain.Buy,
		amount,
		2.0,
		executedPrice,
		executedAt,
	)
	assert.NoError(t, err)
	return transaction
}

func TestCreateTransaction(t *testing.T) {
	t.Run("유효한 거래 생성", func(t *testing.T) {
		// Given
		repo := new(MockTransactionRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		userID := uuid.New()
		portfolioID := uuid.New()
		assetID := uuid.New()
		executedAt := time.Now()

		reqBody := createTransactionRequest{
			UserID:        userID.String(),
			PortfolioID:   portfolioID.String(),
			AssetID:       assetID.String(),
			Type:          string(domain.Buy),
			Amount:        5000.0,
			Quantity:      100,
			ExecutedPrice: 50.0,
			ExecutedAt:    executedAt.Format(time.RFC3339),
		}

		repo.On("Save", mock.Anything, mock.AnythingOfType("*domain.Transaction")).Return(nil)

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusCreated, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("잘못된 요청 데이터", func(t *testing.T) {
		// Given
		repo := new(MockTransactionRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		testCases := []struct {
			name    string
			reqBody interface{}
		}{
			{
				name: "필수 필드 누락",
				reqBody: struct {
					Type string `json:"type"`
				}{
					Type: "BUY",
				},
			},
			{
				name: "잘못된 타입",
				reqBody: struct {
					Type          string    `json:"type"`
					Amount        float64   `json:"amount"`
					Quantity      float64   `json:"quantity"`
					ExecutedPrice float64   `json:"executed_price"`
					Currency      string    `json:"currency"`
					ExecutedAt    time.Time `json:"executed_at"`
				}{
					Type:          "INVALID_TYPE",
					Amount:        -100.0,
					Quantity:      -1,
					ExecutedPrice: -50.0,
					Currency:      "",
					ExecutedAt:    time.Now(),
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, err := json.Marshal(tc.reqBody)
				assert.NoError(t, err)

				req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				rr := httptest.NewRecorder()

				router.ServeHTTP(rr, req)

				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})
		}
	})
}

func TestGetTransaction(t *testing.T) {
	t.Run("존재하는 거래 조회", func(t *testing.T) {
		// Given
		repo := new(MockTransactionRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		transactionID := uuid.MustParse("93910f31-9327-4531-bd3b-9d8c1ae9d0a4")
		userID := uuid.New()
		portfolioID := uuid.New()
		assetID := uuid.New()
		amount, _ := valueobjects.NewMoney(5000.0, "USD")
		executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
		executedAt := time.Now()

		transaction := &domain.Transaction{
			ID:            transactionID,
			UserID:        userID,
			PortfolioID:   portfolioID,
			AssetID:       assetID,
			Type:          domain.Buy,
			Amount:        amount,
			Quantity:      100,
			ExecutedPrice: executedPrice,
			ExecutedAt:    executedAt,
		}

		repo.On("FindByID", mock.Anything, transactionID).Return(transaction, nil)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+transactionID.String(), nil)
		rr := httptest.NewRecorder()

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusOK, rr.Code)
		repo.AssertExpectations(t)
	})

	t.Run("존재하지 않는 거래 조회", func(t *testing.T) {
		// Given
		repo := new(MockTransactionRepository)
		handler := NewHandler(repo)
		router := setupTestRouter(handler)

		transactionID := uuid.New()
		repo.On("FindByID", mock.Anything, transactionID).Return(nil, domain.ErrTransactionNotFound)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+transactionID.String(), nil)
		rr := httptest.NewRecorder()

		// When
		router.ServeHTTP(rr, req)

		// Then
		assert.Equal(t, http.StatusNotFound, rr.Code)
		repo.AssertExpectations(t)
	})
}

func TestListUserTransactions(t *testing.T) {
	// Given
	repo := new(MockTransactionRepository)
	handler := NewHandler(repo)
	router := setupTestRouter(handler)

	userID := uuid.New()
	transactions := []*domain.Transaction{createValidTransaction(t), createValidTransaction(t)}

	repo.On("FindByUserID", mock.Anything, userID).Return(transactions, nil)

	// When
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/"+userID.String()+"/transactions", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Then
	assert.Equal(t, http.StatusOK, rr.Code)
	repo.AssertExpectations(t)
}

func TestListPortfolioTransactions(t *testing.T) {
	// Given
	repo := new(MockTransactionRepository)
	handler := NewHandler(repo)
	router := setupTestRouter(handler)

	portfolioID := uuid.New()
	amount, _ := valueobjects.NewMoney(5000.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")

	transactions := []*domain.Transaction{
		{
			ID:            uuid.New(),
			PortfolioID:   portfolioID,
			Type:          domain.Buy,
			Amount:        amount,
			Quantity:      100,
			ExecutedPrice: executedPrice,
			ExecutedAt:    time.Now(),
		},
		{
			ID:            uuid.New(),
			PortfolioID:   portfolioID,
			Type:          domain.Sell,
			Amount:        amount,
			Quantity:      50,
			ExecutedPrice: executedPrice,
			ExecutedAt:    time.Now(),
		},
	}

	repo.On("FindByPortfolioID", mock.Anything, portfolioID).Return(transactions, nil)

	// When
	req := httptest.NewRequest(http.MethodGet, "/api/v1/portfolios/"+portfolioID.String()+"/transactions", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Then
	assert.Equal(t, http.StatusOK, rr.Code)

	var response []TransactionResponse
	err := json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}
