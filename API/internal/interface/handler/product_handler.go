package handler

import (
	"encoding/json"
	"inventario/internal/domain"
	"inventario/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
}

func NewProductHandler(useCase *usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{
		productUseCase: useCase,
	}
}

type createProductRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	ImageURL string `json:"image_url"`
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Code == "" {
		http.Error(w, "Name and code are required fields", http.StatusBadRequest)
		return
	}

	product, err := h.productUseCase.CreateProduct(req.Name, req.Code, req.ImageURL)
	if err != nil {
		switch e := err.(type) {
		case *domain.ProductAlreadyExistsError:
			http.Error(w, e.Error(), http.StatusConflict)
		default:
			http.Error(w, "Error creating product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		http.Error(w, "Invalid request context", http.StatusBadRequest)
		return
	}
	idStr := rctx.URLParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.productUseCase.GetProduct(id)
	if err != nil {
		switch e := err.(type) {
		case *domain.ProductNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error fetching product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productUseCase.GetAllProducts()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		http.Error(w, "Invalid request context", http.StatusBadRequest)
		return
	}
	idStr := rctx.URLParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if product.Name == "" || product.Code == "" {
		http.Error(w, "Name and code are required fields", http.StatusBadRequest)
		return
	}

	product.ID = id
	if err := h.productUseCase.UpdateProduct(&product); err != nil {
		switch e := err.(type) {
		case *domain.ProductNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error updating product", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		http.Error(w, "Invalid request context", http.StatusBadRequest)
		return
	}
	idStr := rctx.URLParam("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.productUseCase.DeleteProduct(id); err != nil {
		switch e := err.(type) {
		case *domain.ProductNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error deleting product", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
