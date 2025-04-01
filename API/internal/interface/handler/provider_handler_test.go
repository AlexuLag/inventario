package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"inventario/internal/domain"
	"inventario/internal/infrastructure/repository"
	"inventario/internal/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockCreate     func(*domain.Provider) error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful creation",
			requestBody: map[string]string{
				"name":    "Test Provider",
				"email":   "test@example.com",
				"phone":   "1234567890",
				"address": "Test Address",
			},
			mockCreate: func(p *domain.Provider) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "provider already exists",
			requestBody: map[string]string{
				"name":    "Existing Provider",
				"email":   "existing@example.com",
				"phone":   "0987654321",
				"address": "Existing Address",
			},
			mockCreate: func(p *domain.Provider) error {
				return &domain.ProviderAlreadyExistsError{Email: "existing@example.com"}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "provider already exists",
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"name": "Invalid Provider",
				// Missing required fields
			},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := usecase.NewProviderUseCase(mockRepo)
			handler := NewProviderHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/providers", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateProvider(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name           string
		providerID     string
		mockGetByID    func(int64) (*domain.Provider, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful retrieval",
			providerID: "1",
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return &domain.Provider{
					ID:        1,
					Name:      "Test Provider",
					Email:     "test@example.com",
					Phone:     "1234567890",
					Address:   "Test Address",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:       "provider not found",
			providerID: "999",
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "provider not found",
		},
		{
			name:           "invalid provider ID",
			providerID:     "invalid",
			mockGetByID:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid provider ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := usecase.NewProviderUseCase(mockRepo)
			handler := NewProviderHandler(useCase)

			req := httptest.NewRequest("GET", "/api/providers/"+tt.providerID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.providerID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.GetProvider(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetAllProviders(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAll     func() ([]domain.Provider, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful retrieval",
			mockGetAll: func() ([]domain.Provider, error) {
				return []domain.Provider{
					{
						ID:        1,
						Name:      "Provider 1",
						Email:     "provider1@example.com",
						Phone:     "1234567890",
						Address:   "Address 1",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Name:      "Provider 2",
						Email:     "provider2@example.com",
						Phone:     "0987654321",
						Address:   "Address 2",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "database error",
			mockGetAll: func() ([]domain.Provider, error) {
				return nil, domain.ErrInternalServer
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := usecase.NewProviderUseCase(mockRepo)
			handler := NewProviderHandler(useCase)

			req := httptest.NewRequest("GET", "/api/providers", nil)
			w := httptest.NewRecorder()

			handler.GetAllProviders(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestUpdateProvider(t *testing.T) {
	tests := []struct {
		name           string
		providerID     string
		requestBody    map[string]string
		mockGetByID    func(int64) (*domain.Provider, error)
		mockUpdate     func(*domain.Provider) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful update",
			providerID: "1",
			requestBody: map[string]string{
				"name":    "Updated Provider",
				"email":   "updated@example.com",
				"phone":   "1234567890",
				"address": "Updated Address",
			},
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return &domain.Provider{
					ID:        1,
					Name:      "Original Provider",
					Email:     "original@example.com",
					Phone:     "0987654321",
					Address:   "Original Address",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate: func(p *domain.Provider) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:       "provider not found",
			providerID: "999",
			requestBody: map[string]string{
				"name":    "Non-existent Provider",
				"email":   "nonexistent@example.com",
				"phone":   "1234567890",
				"address": "Non-existent Address",
			},
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			mockUpdate: func(p *domain.Provider) error {
				return nil
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "provider not found",
		},
		{
			name:       "invalid request body",
			providerID: "1",
			requestBody: map[string]string{
				"name": "Invalid Provider",
				// Missing required fields
			},
			mockGetByID:    nil,
			mockUpdate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
				UpdateFunc:  tt.mockUpdate,
			}
			useCase := usecase.NewProviderUseCase(mockRepo)
			handler := NewProviderHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/providers/"+tt.providerID, bytes.NewBuffer(body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.providerID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.UpdateProvider(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteProvider(t *testing.T) {
	tests := []struct {
		name           string
		providerID     string
		mockGetByID    func(int64) (*domain.Provider, error)
		mockDelete     func(int64) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:       "successful deletion",
			providerID: "1",
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return &domain.Provider{
					ID:        1,
					Name:      "Test Provider",
					Email:     "test@example.com",
					Phone:     "1234567890",
					Address:   "Test Address",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockDelete: func(id int64) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  "",
		},
		{
			name:       "provider not found",
			providerID: "999",
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			mockDelete: func(id int64) error {
				return nil
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "provider not found",
		},
		{
			name:           "invalid provider ID",
			providerID:     "invalid",
			mockGetByID:    nil,
			mockDelete:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid provider ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
				DeleteFunc:  tt.mockDelete,
			}
			useCase := usecase.NewProviderUseCase(mockRepo)
			handler := NewProviderHandler(useCase)

			req := httptest.NewRequest("DELETE", "/api/providers/"+tt.providerID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.providerID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.DeleteProvider(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
