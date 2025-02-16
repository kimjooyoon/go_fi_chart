package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/transaction/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func setupTestHandler() *Handler {
	return NewHandler()
}

func setupTestRouter(h *Handler) *chi.Mux {
	r := chi.NewRouter()
	h.RegisterRoutes(r)
	return r
}

func createValidTransaction(t *testing.T) *domain.Transaction {
	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	amount, _ := valueobjects.NewMoney(100.0, "USD")
	executedPrice, _ := valueobjects.NewMoney(50.0, "USD")
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
	h := setupTestHandler()
	router := setupTestRouter(h)

	userID := uuid.New()
	portfolioID := uuid.New()
	assetID := uuid.New()
	executedAt := time.Now()

	reqBody := createTransactionRequest{
		UserID:        userID.String(),
		PortfolioID:   portfolioID.String(),
		AssetID:       assetID.String(),
		Type:          string(domain.Buy),
		Amount:        100.0,
		Quantity:      2.0,
		ExecutedPrice: 50.0,
		ExecutedAt:    executedAt.Format(time.RFC3339),
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)

	var response TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, userID.String(), response.UserID)
	assert.Equal(t, portfolioID.String(), response.PortfolioID)
	assert.Equal(t, assetID.String(), response.AssetID)
}

func TestGetTransaction(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+transaction.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, transaction.ID.String(), response.ID)
}

func TestListUserTransactions(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/user/"+transaction.UserID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, transaction.ID.String(), response[0].ID)
}

func TestListPortfolioTransactions(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/portfolio/"+transaction.PortfolioID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, transaction.ID.String(), response[0].ID)
}

func TestListAssetTransactions(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/transactions/asset/"+transaction.AssetID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response []TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, transaction.ID.String(), response[0].ID)
}

func TestUpdateTransaction(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	newAmount := 200.0
	newQuantity := 4.0
	newExecutedPrice := 100.0
	newExecutedAt := time.Now()

	reqBody := createTransactionRequest{
		UserID:        transaction.UserID.String(),
		PortfolioID:   transaction.PortfolioID.String(),
		AssetID:       transaction.AssetID.String(),
		Type:          string(domain.Sell),
		Amount:        newAmount,
		Quantity:      newQuantity,
		ExecutedPrice: newExecutedPrice,
		ExecutedAt:    newExecutedAt.Format(time.RFC3339),
	}

	body, err := json.Marshal(reqBody)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/transactions/"+transaction.ID.String(), bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var response TransactionResponse
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, string(domain.Sell), response.Type)
	assert.Equal(t, newAmount, response.Amount)
	assert.Equal(t, newQuantity, response.Quantity)
	assert.Equal(t, newExecutedPrice, response.ExecutedPrice)
}

func TestDeleteTransaction(t *testing.T) {
	h := setupTestHandler()
	router := setupTestRouter(h)

	transaction := createValidTransaction(t)
	err := h.repository.Save(context.TODO(), transaction)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/transactions/"+transaction.ID.String(), nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	// 삭제된 거래 조회 시도
	req = httptest.NewRequest(http.MethodGet, "/api/v1/transactions/"+transaction.ID.String(), nil)
	rec = httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
}
