package handler

import (
	"encoding/json"
	"inventario/internal/domain"
	"inventario/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type StockHandler struct {
	stockUseCase *usecase.StockUseCase
}

func NewStockHandler(stockUseCase *usecase.StockUseCase) *StockHandler {
	return &StockHandler{
		stockUseCase: stockUseCase,
	}
}

type CreateStockRequest struct {
	ProductID       int64  `json:"product_id"`
	Serial          string `json:"serial"`
	Batch           string `json:"batch"`
	PurchaseDate    string `json:"purchase_date"`
	ProviderID      int64  `json:"provider_id"`
	CreatedByUserID int64  `json:"created_by_user_id"`
	UpdatedByUserID int64  `json:"updated_by_user_id"`
}

func (h *StockHandler) CreateStock(w http.ResponseWriter, r *http.Request) {
	var req CreateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ProductID == 0 || req.Serial == "" || req.ProviderID == 0 || req.CreatedByUserID == 0 || req.UpdatedByUserID == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse purchase date
	var purchaseDate time.Time
	var err error
	if req.PurchaseDate != "" {
		purchaseDate, err = time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			http.Error(w, "Invalid purchase date format", http.StatusBadRequest)
			return
		}
	}

	createdStock, err := h.stockUseCase.CreateStock(
		req.ProductID,
		req.Serial,
		req.Batch,
		purchaseDate,
		req.ProviderID,
		req.CreatedByUserID,
	)
	if err != nil {
		switch e := err.(type) {
		case *domain.StockAlreadyExistsError:
			http.Error(w, "stock with serial "+e.Serial+" already exists", http.StatusConflict)
		default:
			http.Error(w, "Error creating stock", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdStock)
}

func (h *StockHandler) GetStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	stock, err := h.stockUseCase.GetStock(id)
	if err != nil {
		switch e := err.(type) {
		case *domain.StockNotFoundError:
			http.Error(w, "stock with ID "+strconv.FormatInt(e.StockID, 10)+" not found", http.StatusNotFound)
		default:
			http.Error(w, "Error fetching stock", http.StatusInternalServerError)
		}
		return
	}

	if stock == nil {
		http.Error(w, "stock with ID "+strconv.FormatInt(id, 10)+" not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) GetAllStocks(w http.ResponseWriter, r *http.Request) {
	stocks, err := h.stockUseCase.GetAllStocks()
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

type UpdateStockRequest struct {
	ProductID       int64  `json:"product_id"`
	Serial          string `json:"serial"`
	Batch           string `json:"batch"`
	PurchaseDate    string `json:"purchase_date"`
	ProviderID      int64  `json:"provider_id"`
	UpdatedByUserID int64  `json:"updated_by_user_id"`
}

func (h *StockHandler) UpdateStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	var req UpdateStockRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.ProductID == 0 || req.Serial == "" || req.ProviderID == 0 || req.UpdatedByUserID == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse purchase date
	var purchaseDate time.Time
	if req.PurchaseDate != "" {
		purchaseDate, err = time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			http.Error(w, "Invalid purchase date format", http.StatusBadRequest)
			return
		}
	}

	stock := &domain.Stock{
		ID:            id,
		Serial:        req.Serial,
		Batch:         req.Batch,
		PurchaseDate:  purchaseDate,
		Product:       &domain.Product{ID: req.ProductID},
		Provider:      &domain.Provider{ID: req.ProviderID},
		UpdatedByUser: &domain.User{ID: req.UpdatedByUserID},
	}

	if err := h.stockUseCase.UpdateStock(stock); err != nil {
		switch e := err.(type) {
		case *domain.StockNotFoundError:
			http.Error(w, "stock with ID "+strconv.FormatInt(e.StockID, 10)+" not found", http.StatusNotFound)
		default:
			http.Error(w, "Error updating stock", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func (h *StockHandler) DeleteStock(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	if err := h.stockUseCase.DeleteStock(id); err != nil {
		switch e := err.(type) {
		case *domain.StockNotFoundError:
			http.Error(w, "stock with ID "+strconv.FormatInt(e.StockID, 10)+" not found", http.StatusNotFound)
		default:
			http.Error(w, "Error deleting stock", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StockHandler) GetStocksByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productId")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	stocks, err := h.stockUseCase.GetStocksByProductID(productID)
	if err != nil {
		http.Error(w, "Error fetching stocks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

func (h *StockHandler) GetStockBySerial(w http.ResponseWriter, r *http.Request) {
	serial := chi.URLParam(r, "serial")
	if serial == "" {
		http.Error(w, "Serial number is required", http.StatusBadRequest)
		return
	}

	stock, err := h.stockUseCase.GetStockBySerial(serial)
	if err != nil {
		switch e := err.(type) {
		case *domain.StockNotFoundError:
			http.Error(w, "stock with serial "+e.Serial+" not found", http.StatusNotFound)
		default:
			http.Error(w, "Error fetching stock", http.StatusInternalServerError)
		}
		return
	}

	if stock == nil {
		http.Error(w, "stock with serial "+serial+" not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}
