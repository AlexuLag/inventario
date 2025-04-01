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

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockCreate     func(*domain.Product) error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful creation",
			requestBody: map[string]string{
				"name":      "Test Product",
				"code":      "TEST123",
				"image_url": "http://example.com/image.jpg",
			},
			mockCreate: func(p *domain.Product) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "product already exists",
			requestBody: map[string]string{
				"name":      "Existing Product",
				"code":      "EXIST123",
				"image_url": "http://example.com/existing.jpg",
			},
			mockCreate: func(p *domain.Product) error {
				return &domain.ProductAlreadyExistsError{Code: "EXIST123"}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "product with code EXIST123 already exists",
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"name": "Invalid Product",
				// Missing required fields
			},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name and code are required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := usecase.NewProductUseCase(mockRepo)
			handler := NewProductHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/products", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockGetByID    func(int64) (*domain.Product, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful retrieval",
			productID: "1",
			mockGetByID: func(id int64) (*domain.Product, error) {
				return &domain.Product{
					ID:        1,
					Name:      "Test Product",
					Code:      "TEST123",
					ImageURL:  "http://example.com/image.jpg",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:      "product not found",
			productID: "999",
			mockGetByID: func(id int64) (*domain.Product, error) {
				return nil, &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "product with ID 999 not found",
		},
		{
			name:           "invalid product ID",
			productID:      "invalid",
			mockGetByID:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid product ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := usecase.NewProductUseCase(mockRepo)
			handler := NewProductHandler(useCase)

			req := httptest.NewRequest("GET", "/api/products/"+tt.productID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.productID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.GetProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAll     func() ([]*domain.Product, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful retrieval",
			mockGetAll: func() ([]*domain.Product, error) {
				return []*domain.Product{
					{
						ID:        1,
						Name:      "Product 1",
						Code:      "CODE1",
						ImageURL:  "http://example.com/1.jpg",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Name:      "Product 2",
						Code:      "CODE2",
						ImageURL:  "http://example.com/2.jpg",
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
			mockGetAll: func() ([]*domain.Product, error) {
				return nil, domain.ErrInternalServer
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := usecase.NewProductUseCase(mockRepo)
			handler := NewProductHandler(useCase)

			req := httptest.NewRequest("GET", "/api/products", nil)
			w := httptest.NewRecorder()

			handler.GetAllProducts(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		requestBody    map[string]string
		mockGetByID    func(int64) (*domain.Product, error)
		mockUpdate     func(*domain.Product) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful update",
			productID: "1",
			requestBody: map[string]string{
				"name":      "Updated Product",
				"code":      "UPD123",
				"image_url": "http://example.com/updated.jpg",
			},
			mockGetByID: func(id int64) (*domain.Product, error) {
				return &domain.Product{
					ID:        1,
					Name:      "Original Product",
					Code:      "ORIG123",
					ImageURL:  "http://example.com/original.jpg",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate: func(p *domain.Product) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:      "product not found",
			productID: "999",
			requestBody: map[string]string{
				"name":      "Non-existent Product",
				"code":      "NON123",
				"image_url": "http://example.com/none.jpg",
			},
			mockGetByID: func(id int64) (*domain.Product, error) {
				return nil, &domain.ProductNotFoundError{ProductID: 999}
			},
			mockUpdate: func(p *domain.Product) error {
				return &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "product with ID 999 not found",
		},
		{
			name:      "invalid request body",
			productID: "1",
			requestBody: map[string]string{
				"name": "Invalid Product",
				// Missing required fields
			},
			mockGetByID:    nil,
			mockUpdate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name and code are required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetByIDFunc: tt.mockGetByID,
				UpdateFunc:  tt.mockUpdate,
			}
			useCase := usecase.NewProductUseCase(mockRepo)
			handler := NewProductHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/products/"+tt.productID, bytes.NewBuffer(body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.productID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.UpdateProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name           string
		productID      string
		mockGetByID    func(int64) (*domain.Product, error)
		mockDelete     func(int64) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:      "successful deletion",
			productID: "1",
			mockGetByID: func(id int64) (*domain.Product, error) {
				return &domain.Product{
					ID:        1,
					Name:      "Test Product",
					Code:      "TEST123",
					ImageURL:  "http://example.com/image.jpg",
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
			name:      "product not found",
			productID: "999",
			mockGetByID: func(id int64) (*domain.Product, error) {
				return nil, &domain.ProductNotFoundError{ProductID: 999}
			},
			mockDelete: func(id int64) error {
				return &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "product with ID 999 not found",
		},
		{
			name:           "invalid product ID",
			productID:      "invalid",
			mockGetByID:    nil,
			mockDelete:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid product ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetByIDFunc: tt.mockGetByID,
				DeleteFunc:  tt.mockDelete,
			}
			useCase := usecase.NewProductUseCase(mockRepo)
			handler := NewProductHandler(useCase)

			req := httptest.NewRequest("DELETE", "/api/products/"+tt.productID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.productID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.DeleteProduct(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}
