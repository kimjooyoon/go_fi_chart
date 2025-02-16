package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aske/go_fi_chart/pkg/domain/valueobjects"
	"github.com/aske/go_fi_chart/services/transaction/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// Handler는 Transaction 서비스의 HTTP 핸들러입니다
type Handler struct {
	repository domain.TransactionRepository
}

// NewHandler는 새로운 Handler를 생성합니다
func NewHandler() *Handler {
	return &Handler{
		repository: domain.NewMemoryTransactionRepository(),
	}
}

// RegisterRoutes는 라우터에 API 엔드포인트를 등록합니다
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Route("/api/v1/transactions", func(r chi.Router) {
		r.Get("/", h.ListTransactions)
		r.Post("/", h.CreateTransaction)
		r.Get("/{id}", h.GetTransaction)
		r.Put("/{id}", h.UpdateTransaction)
		r.Delete("/{id}", h.DeleteTransaction)
		r.Get("/user/{userID}", h.ListUserTransactions)
		r.Get("/portfolio/{portfolioID}", h.ListPortfolioTransactions)
		r.Get("/asset/{assetID}", h.ListAssetTransactions)
	})
}

type createTransactionRequest struct {
	UserID        string  `json:"userID"`
	PortfolioID   string  `json:"portfolioID"`
	AssetID       string  `json:"assetID"`
	Type          string  `json:"type"`
	Amount        float64 `json:"amount"`
	Quantity      float64 `json:"quantity"`
	ExecutedPrice float64 `json:"executedPrice"`
	ExecutedAt    string  `json:"executedAt"`
}

type TransactionResponse struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	PortfolioID   string    `json:"portfolio_id"`
	AssetID       string    `json:"asset_id"`
	Type          string    `json:"type"`
	Amount        float64   `json:"amount"`
	Quantity      float64   `json:"quantity"`
	ExecutedPrice float64   `json:"executed_price"`
	ExecutedAt    time.Time `json:"executed_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// CreateTransaction은 새로운 거래를 생성합니다
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	portfolioID, err := uuid.Parse(req.PortfolioID)
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	assetID, err := uuid.Parse(req.AssetID)
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	executedAt, err := time.Parse(time.RFC3339, req.ExecutedAt)
	if err != nil {
		http.Error(w, "Invalid executed at time", http.StatusBadRequest)
		return
	}

	amount, err := valueobjects.NewMoney(req.Amount, "USD")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	executedPrice, err := valueobjects.NewMoney(req.ExecutedPrice, "USD")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := domain.NewTransaction(
		userID,
		portfolioID,
		assetID,
		domain.TransactionType(req.Type),
		amount,
		req.Quantity,
		executedPrice,
		executedAt,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repository.Save(r.Context(), transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(toTransactionResponse(transaction)); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// GetTransaction은 특정 거래를 조회합니다
func (h *Handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	transaction, err := h.repository.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(toTransactionResponse(transaction)); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListTransactions은 모든 거래를 조회합니다
func (h *Handler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// 실제 구현에서는 페이지네이션을 추가해야 합니다
	transactions, err := h.repository.FindByUserID(r.Context(), uuid.Nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = toTransactionResponse(t)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListUserTransactions은 사용자의 모든 거래를 조회합니다
func (h *Handler) ListUserTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.repository.FindByUserID(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = toTransactionResponse(t)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListPortfolioTransactions은 포트폴리오의 모든 거래를 조회합니다
func (h *Handler) ListPortfolioTransactions(w http.ResponseWriter, r *http.Request) {
	portfolioID, err := uuid.Parse(chi.URLParam(r, "portfolioID"))
	if err != nil {
		http.Error(w, "Invalid portfolio ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.repository.FindByPortfolioID(r.Context(), portfolioID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = toTransactionResponse(t)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// ListAssetTransactions은 자산의 모든 거래를 조회합니다
func (h *Handler) ListAssetTransactions(w http.ResponseWriter, r *http.Request) {
	assetID, err := uuid.Parse(chi.URLParam(r, "assetID"))
	if err != nil {
		http.Error(w, "Invalid asset ID", http.StatusBadRequest)
		return
	}

	transactions, err := h.repository.FindByAssetID(r.Context(), assetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := make([]TransactionResponse, len(transactions))
	for i, t := range transactions {
		response[i] = toTransactionResponse(t)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// UpdateTransaction은 거래를 업데이트합니다
func (h *Handler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	var req createTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transaction, err := h.repository.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 업데이트 로직 구현
	// 실제 구현에서는 더 세밀한 업데이트 로직이 필요할 수 있습니다

	if err := h.repository.Update(r.Context(), transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(toTransactionResponse(transaction)); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// DeleteTransaction은 거래를 삭제합니다
func (h *Handler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	if err := h.repository.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func toTransactionResponse(t *domain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:            t.ID.String(),
		UserID:        t.UserID.String(),
		PortfolioID:   t.PortfolioID.String(),
		AssetID:       t.AssetID.String(),
		Type:          string(t.Type),
		Amount:        t.Amount.Amount,
		Quantity:      t.Quantity,
		ExecutedPrice: t.ExecutedPrice.Amount,
		ExecutedAt:    t.ExecutedAt,
		CreatedAt:     t.CreatedAt,
	}
}
