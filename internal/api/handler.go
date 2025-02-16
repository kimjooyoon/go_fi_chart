package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/aske/go_fi_chart/internal/domain/asset"
	"github.com/aske/go_fi_chart/internal/domain/gamification"
	"github.com/go-chi/chi/v5"
)

// Handler API 핸들러입니다.
type Handler struct {
	assetRepo        asset.Repository
	transactionRepo  asset.TransactionRepository
	portfolioRepo    asset.PortfolioRepository
	gamificationRepo gamification.Repository
}

// NewHandler 새로운 API 핸들러를 생성합니다.
func NewHandler(
	assetRepo asset.Repository,
	transactionRepo asset.TransactionRepository,
	portfolioRepo asset.PortfolioRepository,
	gamificationRepo gamification.Repository,
) *Handler {
	return &Handler{
		assetRepo:        assetRepo,
		transactionRepo:  transactionRepo,
		portfolioRepo:    portfolioRepo,
		gamificationRepo: gamificationRepo,
	}
}

// RegisterRoutes 라우터에 API 핸들러를 등록합니다.
func (h *Handler) RegisterRoutes(r chi.Router) {
	// 자산 관리 API
	r.Route("/assets", func(r chi.Router) {
		r.Get("/", h.ListAssets)
		r.Post("/", h.CreateAsset)
		r.Get("/{id}", h.GetAsset)
		r.Put("/{id}", h.UpdateAsset)
		r.Delete("/{id}", h.DeleteAsset)
	})

	// 거래 내역 API
	r.Route("/transactions", func(r chi.Router) {
		r.Get("/", h.ListTransactions)
		r.Post("/", h.CreateTransaction)
		r.Get("/{id}", h.GetTransaction)
	})

	// 포트폴리오 API
	r.Route("/portfolios", func(r chi.Router) {
		r.Get("/", h.GetPortfolio)
		r.Put("/", h.UpdatePortfolio)
	})

	// 게임화 API
	r.Route("/gamification", func(r chi.Router) {
		r.Get("/profile", h.GetProfile)
		r.Get("/badges", h.ListBadges)
		r.Get("/streaks", h.ListStreaks)
		r.Get("/stats", h.GetStats)
	})
}

// 응답 헬퍼 함수
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// 요청/응답 구조체
type CreateAssetRequest struct {
	UserID   string  `json:"userId"`
	Type     string  `json:"type"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type UpdateAssetRequest struct {
	Name     string  `json:"name,omitempty"`
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

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

func (h *Handler) ListAssets(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		respondError(w, http.StatusBadRequest, "사용자 ID가 필요합니다")
		return
	}

	assets, err := h.assetRepo.FindByUserID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "자산 목록 조회 실패")
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
		respondError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	asset := asset.NewAsset(req.UserID, asset.Type(req.Type), req.Name, req.Amount, req.Currency)
	if err := h.assetRepo.Save(r.Context(), asset); err != nil {
		respondError(w, http.StatusInternalServerError, "자산 생성 실패")
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
		respondError(w, http.StatusBadRequest, "자산 ID가 필요합니다")
		return
	}

	asset, err := h.assetRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "자산을 찾을 수 없습니다")
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
		respondError(w, http.StatusBadRequest, "자산 ID가 필요합니다")
		return
	}

	var req UpdateAssetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	asset, err := h.assetRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "자산을 찾을 수 없습니다")
		return
	}

	if req.Name != "" {
		asset.Name = req.Name
	}
	if req.Amount != 0 {
		asset.Amount.Amount = req.Amount
	}
	if req.Currency != "" {
		asset.Amount.Currency = req.Currency
	}
	asset.UpdatedAt = time.Now()

	if err := h.assetRepo.Update(r.Context(), asset); err != nil {
		respondError(w, http.StatusInternalServerError, "자산 업데이트 실패")
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
		respondError(w, http.StatusBadRequest, "자산 ID가 필요합니다")
		return
	}

	if err := h.assetRepo.Delete(r.Context(), id); err != nil {
		respondError(w, http.StatusInternalServerError, "자산 삭제 실패")
		return
	}

	respondJSON(w, http.StatusNoContent, nil)
}

// 거래 내역 요청/응답 구조체
type CreateTransactionRequest struct {
	AssetID     string  `json:"assetId"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

type TransactionResponse struct {
	ID          string    `json:"id"`
	AssetID     string    `json:"assetId"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("assetId")
	if assetID == "" {
		respondError(w, http.StatusBadRequest, "자산 ID가 필요합니다")
		return
	}

	transactions, err := h.transactionRepo.FindByAssetID(r.Context(), assetID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "거래 내역 조회 실패")
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, tx := range transactions {
		response[i] = TransactionResponse{
			ID:          tx.ID,
			AssetID:     tx.AssetID,
			Type:        string(tx.Type),
			Amount:      tx.Amount,
			Category:    tx.Category,
			Description: tx.Description,
			Date:        tx.Date,
			CreatedAt:   tx.CreatedAt,
		}
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "잘못된 요청 형식")
		return
	}

	// 자산 존재 여부 확인
	targetAsset, err := h.assetRepo.FindByID(r.Context(), req.AssetID)
	if err != nil {
		respondError(w, http.StatusNotFound, "자산을 찾을 수 없습니다")
		return
	}

	// 거래 생성
	tx := asset.NewTransaction(
		req.AssetID,
		asset.TransactionType(req.Type),
		req.Amount,
		req.Category,
		req.Description,
	)

	// 거래 유효성 검증
	if err := targetAsset.ValidateTransaction(tx); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 거래 처리 및 자산 업데이트
	if err := targetAsset.ProcessTransaction(tx); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// 트랜잭션 저장
	if err := h.transactionRepo.Save(r.Context(), tx); err != nil {
		respondError(w, http.StatusInternalServerError, "거래 생성 실패")
		return
	}

	// 자산 상태 저장
	if err := h.assetRepo.Update(r.Context(), targetAsset); err != nil {
		// 롤백 처리가 필요할 수 있음
		respondError(w, http.StatusInternalServerError, "자산 상태 업데이트 실패")
		return
	}

	response := TransactionResponse{
		ID:          tx.ID,
		AssetID:     tx.AssetID,
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
		CreatedAt:   tx.CreatedAt,
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "거래 ID가 필요합니다")
		return
	}

	tx, err := h.transactionRepo.FindByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "거래를 찾을 수 없습니다")
		return
	}

	response := TransactionResponse{
		ID:          tx.ID,
		AssetID:     tx.AssetID,
		Type:        string(tx.Type),
		Amount:      tx.Amount,
		Category:    tx.Category,
		Description: tx.Description,
		Date:        tx.Date,
		CreatedAt:   tx.CreatedAt,
	}

	respondJSON(w, http.StatusOK, response)
}

func (h *Handler) GetPortfolio(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}

func (h *Handler) UpdatePortfolio(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}

func (h *Handler) GetProfile(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}

func (h *Handler) ListBadges(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}

func (h *Handler) ListStreaks(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}

func (h *Handler) GetStats(w http.ResponseWriter, _ *http.Request) {
	// TODO: 구현
	respondError(w, http.StatusNotImplemented, "아직 구현되지 않았습니다")
}
