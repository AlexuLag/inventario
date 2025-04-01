package handler

import (
	"bytes"
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

func TestCreateStock(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockCreate     func(*domain.Stock) error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"product_id":         1,
				"serial":             "SERIAL123",
				"batch":              "BATCH001",
				"purchase_date":      time.Now().Format("2006-01-02"),
				"provider_id":        1,
				"created_by_user_id": 1,
				"updated_by_user_id": 1,
			},
			mockCreate: func(s *domain.Stock) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "stock already exists",
			requestBody: map[string]interface{}{
				"product_id":         1,
				"serial":             "EXISTING123",
				"batch":              "BATCH001",
				"purchase_date":      time.Now().Format("2006-01-02"),
				"provider_id":        1,
				"created_by_user_id": 1,
				"updated_by_user_id": 1,
			},
			mockCreate: func(s *domain.Stock) error {
				return &domain.StockAlreadyExistsError{Serial: "EXISTING123"}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "stock with serial EXISTING123 already exists",
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"product_id": 1,
				// Missing required fields
			},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/stocks", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateStock(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetStock(t *testing.T) {
	tests := []struct {
		name           string
		stockID        string
		mockGetByID    func(int64) (*domain.Stock, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful retrieval",
			stockID: "1",
			mockGetByID: func(id int64) (*domain.Stock, error) {
				return &domain.Stock{
					ID:           1,
					Serial:       "SERIAL123",
					Batch:        "BATCH001",
					PurchaseDate: time.Now(),
					Product: &domain.Product{
						ID:       1,
						Name:     "Test Product",
						Code:     "TEST123",
						ImageURL: "http://example.com/image.jpg",
					},
					Provider: &domain.Provider{
						ID:      1,
						Name:    "Test Provider",
						Email:   "provider@example.com",
						Phone:   "1234567890",
						Address: "Test Address",
					},
					CreatedByUser: &domain.User{
						ID:    1,
						Name:  "Test User",
						Email: "user@example.com",
						Role:  "admin",
					},
					UpdatedByUser: &domain.User{
						ID:    1,
						Name:  "Test User",
						Email: "user@example.com",
						Role:  "admin",
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:    "stock not found",
			stockID: "999",
			mockGetByID: func(id int64) (*domain.Stock, error) {
				return nil, &domain.StockNotFoundError{StockID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "stock with ID 999 not found",
		},
		{
			name:           "invalid stock ID",
			stockID:        "invalid",
			mockGetByID:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid stock ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			// Create a new chi router and add the URL parameter
			r := chi.NewRouter()
			r.Get("/{id}", handler.GetStock)

			req := httptest.NewRequest("GET", "/"+tt.stockID, nil)
			w := httptest.NewRecorder()

			// Use the router to handle the request
			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetAllStocks(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAll     func() ([]domain.Stock, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful retrieval",
			mockGetAll: func() ([]domain.Stock, error) {
				return []domain.Stock{
					{
						ID:           1,
						Serial:       "SERIAL123",
						Batch:        "BATCH001",
						PurchaseDate: time.Now(),
						Product: &domain.Product{
							ID:       1,
							Name:     "Test Product",
							Code:     "TEST123",
							ImageURL: "http://example.com/image.jpg",
						},
						Provider: &domain.Provider{
							ID:      1,
							Name:    "Test Provider",
							Email:   "provider@example.com",
							Phone:   "1234567890",
							Address: "Test Address",
						},
						CreatedByUser: &domain.User{
							ID:    1,
							Name:  "Test User",
							Email: "user@example.com",
							Role:  "admin",
						},
						UpdatedByUser: &domain.User{
							ID:    1,
							Name:  "Test User",
							Email: "user@example.com",
							Role:  "admin",
						},
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "database error",
			mockGetAll: func() ([]domain.Stock, error) {
				return nil, &domain.StockNotFoundError{StockID: 1}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Error fetching stocks",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			req := httptest.NewRequest("GET", "/api/stocks", nil)
			w := httptest.NewRecorder()

			handler.GetAllStocks(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetStocksByProductID(t *testing.T) {
	tests := []struct {
		name             string
		productID        string
		mockGetByProduct func(int64) ([]domain.Stock, error)
		expectedStatus   int
		expectedError    string
	}{
		{
			name:      "successful retrieval",
			productID: "1",
			mockGetByProduct: func(id int64) ([]domain.Stock, error) {
				return []domain.Stock{
					{
						ID:           1,
						Serial:       "SERIAL123",
						Batch:        "BATCH001",
						PurchaseDate: time.Now(),
						Product: &domain.Product{
							ID:       1,
							Name:     "Test Product",
							Code:     "TEST123",
							ImageURL: "http://example.com/image.jpg",
						},
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:             "invalid product ID",
			productID:        "invalid",
			mockGetByProduct: nil,
			expectedStatus:   http.StatusBadRequest,
			expectedError:    "Invalid product ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				GetByProductIDFunc: tt.mockGetByProduct,
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			r := chi.NewRouter()
			r.Get("/{productId}", handler.GetStocksByProductID)

			req := httptest.NewRequest("GET", "/"+tt.productID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestGetStockBySerial(t *testing.T) {
	tests := []struct {
		name            string
		serial          string
		mockGetBySerial func(string) (*domain.Stock, error)
		expectedStatus  int
		expectedError   string
	}{
		{
			name:   "successful retrieval",
			serial: "SERIAL123",
			mockGetBySerial: func(serial string) (*domain.Stock, error) {
				return &domain.Stock{
					ID:           1,
					Serial:       "SERIAL123",
					Batch:        "BATCH001",
					PurchaseDate: time.Now(),
					Product: &domain.Product{
						ID:       1,
						Name:     "Test Product",
						Code:     "TEST123",
						ImageURL: "http://example.com/image.jpg",
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:   "stock not found",
			serial: "NONEXISTENT",
			mockGetBySerial: func(serial string) (*domain.Stock, error) {
				return nil, &domain.StockNotFoundError{Serial: "NONEXISTENT"}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "stock with serial NONEXISTENT not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				GetBySerialFunc: tt.mockGetBySerial,
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			r := chi.NewRouter()
			r.Get("/{serial}", handler.GetStockBySerial)

			req := httptest.NewRequest("GET", "/"+tt.serial, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestUpdateStock(t *testing.T) {
	tests := []struct {
		name           string
		stockID        string
		requestBody    map[string]interface{}
		mockUpdate     func(*domain.Stock) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful update",
			stockID: "1",
			requestBody: map[string]interface{}{
				"product_id":         1,
				"serial":             "SERIAL123",
				"batch":              "BATCH001",
				"purchase_date":      time.Now().Format("2006-01-02"),
				"provider_id":        1,
				"updated_by_user_id": 1,
			},
			mockUpdate: func(s *domain.Stock) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:    "stock not found",
			stockID: "999",
			requestBody: map[string]interface{}{
				"product_id":         1,
				"serial":             "SERIAL123",
				"batch":              "BATCH001",
				"purchase_date":      time.Now().Format("2006-01-02"),
				"provider_id":        1,
				"updated_by_user_id": 1,
			},
			mockUpdate: func(s *domain.Stock) error {
				return &domain.StockNotFoundError{StockID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "stock with ID 999 not found",
		},
		{
			name:    "invalid request body",
			stockID: "1",
			requestBody: map[string]interface{}{
				"product_id": 1,
				// Missing required fields
			},
			mockUpdate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				UpdateFunc: tt.mockUpdate,
				GetByIDFunc: func(id int64) (*domain.Stock, error) {
					if tt.name == "successful update" {
						return &domain.Stock{
							ID:            1,
							Serial:        "SERIAL123",
							Batch:         "BATCH001",
							Product:       &domain.Product{ID: 1},
							Provider:      &domain.Provider{ID: 1},
							UpdatedByUser: &domain.User{ID: 1},
						}, nil
					}
					return nil, &domain.StockNotFoundError{StockID: id}
				},
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			r := chi.NewRouter()
			r.Put("/{id}", handler.UpdateStock)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/"+tt.stockID, bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestDeleteStock(t *testing.T) {
	tests := []struct {
		name           string
		stockID        string
		mockDelete     func(int64) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:    "successful deletion",
			stockID: "1",
			mockDelete: func(id int64) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  "",
		},
		{
			name:    "stock not found",
			stockID: "999",
			mockDelete: func(id int64) error {
				return &domain.StockNotFoundError{StockID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "stock with ID 999 not found",
		},
		{
			name:           "invalid stock ID",
			stockID:        "invalid",
			mockDelete:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid stock ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockStockRepository{
				DeleteFunc: tt.mockDelete,
				GetByIDFunc: func(id int64) (*domain.Stock, error) {
					if tt.name == "successful deletion" {
						return &domain.Stock{
							ID:            1,
							Serial:        "SERIAL123",
							Batch:         "BATCH001",
							Product:       &domain.Product{ID: 1},
							Provider:      &domain.Provider{ID: 1},
							UpdatedByUser: &domain.User{ID: 1},
						}, nil
					}
					return nil, &domain.StockNotFoundError{StockID: id}
				},
			}
			useCase := usecase.NewStockUseCase(mockRepo)
			handler := NewStockHandler(useCase)

			r := chi.NewRouter()
			r.Delete("/{id}", handler.DeleteStock)

			req := httptest.NewRequest("DELETE", "/"+tt.stockID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %q, got %q", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}
