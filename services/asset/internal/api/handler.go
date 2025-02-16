package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/asset/internal/domain"
	"github.com/go-chi/chi/v5"
)

// Handler API 핸들러입니다.
type Handler struct {
	assetRepo domain.AssetRepository
}

// NewHandler 새로운 API 핸들러를 생성합니다.
func NewHandler() *Handler {
	return &Handler{
		assetRepo: domain.NewMemoryAssetRepository(),
	}
}

// RegisterRoutes 라우터에 API 핸들러를 등록합니다.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/assets", func(r chi.Router) {
		r.Get("/", h.ListAssets)
		r.Post("/", h.CreateAsset)
		r.Get("/{id}", h.GetAsset)
		r.Put("/{id}", h.UpdateAsset)
		r.Delete("/{id}", h.DeleteAsset)
	})
}

// AssetResponse 자산 응답 구조체
type AssetResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateAssetRequest 자산 생성 요청 구조체
type CreateAssetRequest struct {
	UserID   string  `json:"userId"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
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

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "사용자 ID가 필요합니다")
		return
	}

	assets, err := h.assetRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "자산 목록 조회 실패")
		return
	}

	response := make([]AssetResponse, len(assets))
	for i, asset := range assets {
		response[i] = AssetResponse{
			ID:        asset.ID,
			UserID:    asset.UserID,
			Type:      string(asset.Type),
			Name:      asset.Name,
			Amount:    asset.Amount.Amount,
			Currency:  asset.Amount.Currency,
			CreatedAt: asset.CreatedAt,
			UpdatedAt: asset.UpdatedAt,
		}
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	money, err := valueobjects.NewMoney(req.Amount, req.Currency)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	asset := domain.NewAsset(req.UserID, domain.AssetType(req.Type), req.Name, money)
	if err := h.assetRepo.Save(r.Context(), asset); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "자산 생성 실패")
		return
	}

	response := AssetResponse{
		ID:        asset.ID,
		UserID:    asset.UserID,
		Type:      string(asset.Type),
		Name:      asset.Name,
		Amount:    asset.Amount.Amount,
		Currency:  asset.Amount.Currency,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "자산 ID가 필요합니다")
		return
	}

	asset, err := h.assetRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "자산을 찾을 수 없습니다")
		return
	}

	response := AssetResponse{
		ID:        asset.ID,
		UserID:    asset.UserID,
		Type:      string(asset.Type),
		Name:      asset.Name,
		Amount:    asset.Amount.Amount,
		Currency:  asset.Amount.Currency,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "자산 ID가 필요합니다")
		return
	}

	var req CreateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "잘못된 요청 형식")
		return
	}

	asset, err := h.assetRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, ErrNotFound, "자산을 찾을 수 없습니다")
		return
	}

	money, err := valueobjects.NewMoney(req.Amount, req.Currency)
	if err != nil {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, err.Error())
		return
	}

	asset.Update(req.Name, domain.AssetType(req.Type), money)
	if err := h.assetRepo.Update(r.Context(), asset); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "자산 업데이트 실패")
		return
	}

	response := AssetResponse{
		ID:        asset.ID,
		UserID:    asset.UserID,
		Type:      string(asset.Type),
		Name:      asset.Name,
		Amount:    asset.Amount.Amount,
		Currency:  asset.Amount.Currency,
		CreatedAt: asset.CreatedAt,
		UpdatedAt: asset.UpdatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, ErrInvalidRequest, "자산 ID가 필요합니다")
		return
	}

	if err := h.assetRepo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, ErrInternalServer, "자산 삭제 실패")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
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
