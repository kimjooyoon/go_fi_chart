package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/portfolio/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler API 핸들러입니다.
type Handler struct {
	portfolioRepo domain.PortfolioRepository
}

// NewHandler 새로운 API 핸들러를 생성합니다.
func NewHandler() *Handler {
	return &Handler{
		portfolioRepo: domain.NewMemoryPortfolioRepository(),
	}
}

// RegisterRoutes 라우터에 API 핸들러를 등록합니다.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/portfolios", func(r chi.Router) {
		r.Get("/", h.ListPortfolios)
		r.Post("/", h.CreatePortfolio)
		r.Get("/{id}", h.GetPortfolio)
		r.Put("/{id}", h.UpdatePortfolio)
		r.Delete("/{id}", h.DeletePortfolio)
		r.Post("/{id}/assets", h.AddAsset)
		r.Put("/{id}/assets/{assetId}", h.UpdateAssetWeight)
		r.Delete("/{id}/assets/{assetId}", h.RemoveAsset)
	})
}

// PortfolioResponse 포트폴리오 응답 구조체
type PortfolioResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Assets    []Asset   `json:"assets"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Asset 자산 응답 구조체
type Asset struct {
	AssetID string  `json:"assetId"`
	Weight  float64 `json:"weight"`
}

// CreatePortfolioRequest 포트폴리오 생성 요청 구조체
type CreatePortfolioRequest struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

// AddAssetRequest 자산 추가 요청 구조체
type AddAssetRequest struct {
	AssetID string  `json:"assetId"`
	Weight  float64 `json:"weight"`
}

// UpdateAssetWeightRequest 자산 가중치 업데이트 요청 구조체
type UpdateAssetWeightRequest struct {
	Weight float64 `json:"weight"`
}

// ErrorResponse 에러 응답 구조체
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// API 에러 코드 상수
const (
	ErrInvalidRequest = "INVALID_REQUEST"
	ErrNotFound       = "NOT_FOUND"
	ErrInternalServer = "INTERNAL_SERVER_ERROR"
)

func (h *Handler) ListPortfolios(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "사용자 ID가 필요합니다")
		return
	}

	portfolios, err := h.portfolioRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 목록 조회 실패")
		return
	}

	response := make([]PortfolioResponse, len(portfolios))
	for i, portfolio := range portfolios {
		response[i] = toPortfolioResponse(portfolio)
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) CreatePortfolio(w http.ResponseWriter, r *http.Request) {
	var req CreatePortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	portfolio := domain.NewPortfolio(req.UserID, req.Name)
	if err := h.portfolioRepo.Save(r.Context(), portfolio); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 생성 실패")
		return
	}

	respondJSON(w, http.StatusCreated, toPortfolioResponse(portfolio))
}

func (h *Handler) GetPortfolio(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID가 필요합니다")
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "포트폴리오를 찾을 수 없습니다")
		return
	}

	respondJSON(w, http.StatusOK, toPortfolioResponse(portfolio))
}

func (h *Handler) UpdatePortfolio(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID가 필요합니다")
		return
	}

	var req CreatePortfolioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "포트폴리오를 찾을 수 없습니다")
		return
	}

	portfolio.Name = req.Name
	portfolio.UpdatedAt = time.Now()

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 업데이트 실패")
		return
	}

	respondJSON(w, http.StatusOK, toPortfolioResponse(portfolio))
}

func (h *Handler) DeletePortfolio(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID가 필요합니다")
		return
	}

	if err := h.portfolioRepo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 삭제 실패")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

func (h *Handler) AddAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID가 필요합니다")
		return
	}

	var req AddAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "포트폴리오를 찾을 수 없습니다")
		return
	}

	weight, err := valueobjects.NewPercentage(req.Weight)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	if err := portfolio.AddAsset(req.AssetID, weight); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 업데이트 실패")
		return
	}

	respondJSON(w, http.StatusOK, toPortfolioResponse(portfolio))
}

func (h *Handler) UpdateAssetWeight(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	assetID := chi.URLParam(r, "assetId")
	if id == "" || assetID == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID와 자산 ID가 필요합니다")
		return
	}

	var req UpdateAssetWeightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "포트폴리오를 찾을 수 없습니다")
		return
	}

	weight, err := valueobjects.NewPercentage(req.Weight)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	if err := portfolio.UpdateAssetWeight(assetID, weight); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 업데이트 실패")
		return
	}

	respondJSON(w, http.StatusOK, toPortfolioResponse(portfolio))
}

func (h *Handler) RemoveAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	assetID := chi.URLParam(r, "assetId")
	if id == "" || assetID == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "포트폴리오 ID와 자산 ID가 필요합니다")
		return
	}

	portfolio, err := h.portfolioRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "포트폴리오를 찾을 수 없습니다")
		return
	}

	if err := portfolio.RemoveAsset(assetID); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	if err := h.portfolioRepo.Update(r.Context(), portfolio); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "포트폴리오 업데이트 실패")
		return
	}

	respondJSON(w, http.StatusOK, toPortfolioResponse(portfolio))
}

// 응답 헬퍼 함수
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			respondError(w, http.StatusInternalServerError, ErrInternalServer, "응답 생성 중 오류가 발생했습니다")
		}
	}
}

func respondError(w http.ResponseWriter, status int, code string, message string) {
	respondJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
	})
}

func toPortfolioResponse(p *domain.Portfolio) PortfolioResponse {
	assets := make([]Asset, len(p.Assets))
	for i, asset := range p.Assets {
		assets[i] = Asset{
			AssetID: asset.AssetID,
			Weight:  asset.Weight.Value,
		}
	}

	return PortfolioResponse{
		ID:        p.ID,
		UserID:    p.UserID,
		Name:      p.Name,
		Assets:    assets,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}
