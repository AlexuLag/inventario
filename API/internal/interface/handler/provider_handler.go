package handler

import (
	"encoding/json"
	"inventario/internal/domain"
	"inventario/internal/usecase"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ProviderHandler struct {
	providerUseCase *usecase.ProviderUseCase
}

func NewProviderHandler(useCase *usecase.ProviderUseCase) *ProviderHandler {
	return &ProviderHandler{
		providerUseCase: useCase,
	}
}

type createProviderRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

func (h *ProviderHandler) CreateProvider(w http.ResponseWriter, r *http.Request) {
	var req createProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	provider, err := h.providerUseCase.CreateProvider(req.Name, req.Email, req.Phone, req.Address)
	if err != nil {
		switch e := err.(type) {
		case *domain.ProviderAlreadyExistsError:
			http.Error(w, e.Error(), http.StatusConflict)
		default:
			http.Error(w, "Error creating provider", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(provider)
}

func (h *ProviderHandler) GetProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	provider, err := h.providerUseCase.GetProvider(id)
	if err != nil {
		switch e := err.(type) {
		case *domain.ProviderNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error fetching provider", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

func (h *ProviderHandler) GetAllProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.providerUseCase.GetAllProviders()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

func (h *ProviderHandler) UpdateProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	var provider domain.Provider
	if err := json.NewDecoder(r.Body).Decode(&provider); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if provider.Name == "" || provider.Email == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	provider.ID = id
	if err := h.providerUseCase.UpdateProvider(&provider); err != nil {
		switch e := err.(type) {
		case *domain.ProviderNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error updating provider", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(provider)
}

func (h *ProviderHandler) DeleteProvider(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid provider ID", http.StatusBadRequest)
		return
	}

	if err := h.providerUseCase.DeleteProvider(id); err != nil {
		switch e := err.(type) {
		case *domain.ProviderNotFoundError:
			http.Error(w, e.Error(), http.StatusNotFound)
		default:
			http.Error(w, "Error deleting provider", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
