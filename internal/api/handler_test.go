package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aske/go_fi_chart/internal/domain/asset"
	"github.com/aske/go_fi_chart/internal/domain/gamification"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTestHandler() (*Handler, *asset.MockRepository) {
	mockAssetRepo := new(asset.MockRepository)
	mockTransactionRepo := new(asset.MockTransactionRepository)
	mockPortfolioRepo := new(asset.MockPortfolioRepository)
	mockGamificationRepo := new(gamification.MockRepository)

	handler := NewHandler(mockAssetRepo, mockTransactionRepo, mockPortfolioRepo, mockGamificationRepo)
	return handler, mockAssetRepo
}

func TestListAssets(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	assets := []*asset.Asset{
		{
			ID:     "1",
			UserID: "user1",
			Type:   asset.Cash,
			Name:   "Test Asset",
			Amount: asset.Money{
				Amount:   100.0,
				Currency: "USD",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockRepo.On("FindByUserID", mock.Anything, "user1").Return(assets, nil)

	req := httptest.NewRequest(http.MethodGet, "/assets?userId=user1", nil)
	w := httptest.NewRecorder()

	handler.ListAssets(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 1)
	assert.Equal(t, assets[0].ID, response[0].ID)
}

func TestCreateAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	createReq := CreateAssetRequest{
		UserID:   "user1",
		Type:     "CASH",
		Name:     "Test Asset",
		Amount:   100.0,
		Currency: "USD",
	}

	mockRepo.On("Save", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateAsset(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, createReq.UserID, response.UserID)
	assert.Equal(t, createReq.Name, response.Name)
}

func TestGetAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	testAsset := &asset.Asset{
		ID:     "1",
		UserID: "user1",
		Type:   asset.Cash,
		Name:   "Test Asset",
		Amount: asset.Money{
			Amount:   100.0,
			Currency: "USD",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("FindByID", mock.Anything, "1").Return(testAsset, nil)

	req := httptest.NewRequest(http.MethodGet, "/assets/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.GetAsset(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response AssetResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, testAsset.ID, response.ID)
}

func TestUpdateAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	testAsset := &asset.Asset{
		ID:     "1",
		UserID: "user1",
		Type:   asset.Cash,
		Name:   "Test Asset",
		Amount: asset.Money{
			Amount:   100.0,
			Currency: "USD",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updateReq := UpdateAssetRequest{
		Name:     "Updated Asset",
		Amount:   200.0,
		Currency: "USD",
	}

	mockRepo.On("FindByID", mock.Anything, "1").Return(testAsset, nil)
	mockRepo.On("Update", mock.Anything, mock.AnythingOfType("*asset.Asset")).Return(nil)

	body, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/assets/1", bytes.NewReader(body))
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.UpdateAsset(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteAsset(t *testing.T) {
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", mock.Anything, "1").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/assets/1", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	handler.DeleteAsset(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
