package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/portfolio/internal/domain"
	"github.com/gorilla/mux"
)

type Handler struct {
	portfolioRepo domain.PortfolioRepository
	logger        *slog.Logger
}

func NewHandler(portfolioRepo domain.PortfolioRepository, logger *slog.Logger) *Handler {
	return &Handler{
		portfolioRepo: portfolioRepo,
		logger:        logger,
	}
}

func (h *Handler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/portfolios", h.CreatePortfolio).Methods("POST")
	r.HandleFunc("/portfolios/{id}", h.GetPortfolio).Methods("GET")
	r.HandleFunc("/portfolios/{id}", h.UpdatePortfolio).Methods("PUT")
	r.HandleFunc("/portfolios/{id}", h.DeletePortfolio).Methods("DELETE")
	r.HandleFunc("/portfolios/{id}/assets", h.AddAsset).Methods("POST")
	r.HandleFunc("/portfolios/{id}/assets/{assetId}", h.UpdateAssetWeight).Methods("PUT")
	r.HandleFunc("/portfolios/{id}/assets/{assetId}", h.RemoveAsset).Methods("DELETE")
	r.HandleFunc("/users/{userId}/portfolios", h.ListUserPortfolios).Methods("GET")
	r.HandleFunc("/portfolios", h.ListPortfolios).Methods("GET")
}

type assetResponse struct {
	AssetID string  `json:"assetId"`
	Weight  float64 `json:"weight"`
}

type portfolioResponse struct {
	ID        string          `json:"id"`
	UserID    string          `json:"userId"`
	Name      string          `json:"name"`
	Assets    []assetResponse `json:"assets"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
}

type createPortfolioRequest struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

type addAssetRequest struct {
	AssetID string  `json:"assetId"`
	Weight  float64 `json:"weight"`
}

type updateAssetWeightRequest struct {
	Weight float64 `json:"weight"`
}

func toPortfolioResponse(p *domain.Portfolio) portfolioResponse {
	assets := make([]assetResponse, len(p.Assets))
	for i, asset := range p.Assets {
		assets[i] = assetResponse{
			AssetID: asset.AssetID,
			Weight:  asset.Weight.Value,
		}
	}
	return portfolioResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Name:      p.Name,
		Assets:    assets,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (h *Handler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	var req createPortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	portfolio := domain.NewPortfolio(req.UserID, req.Name)
	if err := h.portfolioRepo.Save(r.Context(), portfolio); err != nil {
		h.logger.Error("failed to save portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toPortfolioResponse(portfolio)); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toPortfolioResponse(portfolio)); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req createPortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	portfolio.Name = req.Name
	portfolio.UpdatedAt = time.Now()

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		h.logger.Error("failed to update portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toPortfolioResponse(portfolio)); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	if err := h.portfolioRepo.Delete(r.Context(), id); err != nil {
		h.logger.Error("failed to delete portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) AddAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req addAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	weight, err := valueobjects.NewPercentage(req.Weight)
	if err != nil {
		h.logger.Error("failed to create percentage", "error", err)
		http.Error(w, "invalid weight", http.StatusBadRequest)
		return
	}

	if err := portfolio.AddAsset(req.AssetID, weight); err != nil {
		h.logger.Error("failed to add asset", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		h.logger.Error("failed to update portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toPortfolioResponse(portfolio)); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateAssetWeight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	assetID := vars["assetId"]

	var req updateAssetWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	weight, err := valueobjects.NewPercentage(req.Weight)
	if err != nil {
		h.logger.Error("failed to create percentage", "error", err)
		http.Error(w, "invalid weight", http.StatusBadRequest)
		return
	}

	if err := portfolio.UpdateAssetWeight(assetID, weight); err != nil {
		h.logger.Error("failed to update asset weight", "error", err)
		var assetNotFoundError domain.AssetNotFoundError
		switch {
		case errors.As(err, &assetNotFoundError):
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		h.logger.Error("failed to update portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(toPortfolioResponse(portfolio)); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) RemoveAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	assetID := vars["assetId"]

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		h.logger.Error("failed to find portfolio", "error", err)
		http.Error(w, "portfolio not found", http.StatusNotFound)
		return
	}

	if err := portfolio.RemoveAsset(assetID); err != nil {
		h.logger.Error("failed to remove asset", "error", err)
		var assetNotFoundError domain.AssetNotFoundError
		switch {
		case errors.As(err, &assetNotFoundError):
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		h.logger.Error("failed to update portfolio", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListUserPortfolios(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	portfolios, err := h.portfolioRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to find portfolios", "error", err)
		if errors.Is(err, domain.ErrPortfolioNotFound) {
			http.Error(w, "portfolios not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]portfolioResponse, len(portfolios))
	for i, p := range portfolios {
		response[i] = toPortfolioResponse(p)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ListPortfolios(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	portfolios, err := h.portfolioRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		h.logger.Error("failed to find portfolios", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]portfolioResponse, len(portfolios))
	for i, p := range portfolios {
		response[i] = toPortfolioResponse(p)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
