package usecase

import (
	"errors"
	"inventario/internal/domain"
	"inventario/internal/infrastructure/repository"
	"testing"
)

func TestCreateProduct(t *testing.T) {
	tests := []struct {
		name          string
		productName   string
		productCode   string
		imageURL      string
		mockCreate    func(*domain.Product) error
		expectedError error
	}{
		{
			name:        "successful creation",
			productName: "Test Product",
			productCode: "TEST123",
			imageURL:    "http://example.com/image.jpg",
			mockCreate: func(p *domain.Product) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:        "product already exists",
			productName: "Existing Product",
			productCode: "EXIST123",
			imageURL:    "http://example.com/existing.jpg",
			mockCreate: func(p *domain.Product) error {
				return &domain.ProductAlreadyExistsError{Code: "EXIST123"}
			},
			expectedError: &domain.ProductAlreadyExistsError{Code: "EXIST123"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := NewProductUseCase(mockRepo)

			product, err := useCase.CreateProduct(tt.productName, tt.productCode, tt.imageURL)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if product.Name != tt.productName {
				t.Errorf("expected name %s, got %s", tt.productName, product.Name)
			}
			if product.Code != tt.productCode {
				t.Errorf("expected code %s, got %s", tt.productCode, product.Code)
			}
			if product.ImageURL != tt.imageURL {
				t.Errorf("expected image URL %s, got %s", tt.imageURL, product.ImageURL)
			}
		})
	}
}

func TestGetProduct(t *testing.T) {
	tests := []struct {
		name          string
		productID     int64
		mockGetByID   func(int64) (*domain.Product, error)
		expectedError error
	}{
		{
			name:      "successful retrieval",
			productID: 1,
			mockGetByID: func(id int64) (*domain.Product, error) {
				return &domain.Product{
					ID:       1,
					Name:     "Test Product",
					Code:     "TEST123",
					ImageURL: "http://example.com/image.jpg",
				}, nil
			},
			expectedError: nil,
		},
		{
			name:      "product not found",
			productID: 999,
			mockGetByID: func(id int64) (*domain.Product, error) {
				return nil, &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedError: &domain.ProductNotFoundError{ProductID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := NewProductUseCase(mockRepo)

			product, err := useCase.GetProduct(tt.productID)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if product.ID != tt.productID {
				t.Errorf("expected ID %d, got %d", tt.productID, product.ID)
			}
		})
	}
}

func TestGetAllProducts(t *testing.T) {
	tests := []struct {
		name          string
		mockGetAll    func() ([]*domain.Product, error)
		expectedError error
	}{
		{
			name: "successful retrieval",
			mockGetAll: func() ([]*domain.Product, error) {
				return []*domain.Product{
					{
						ID:       1,
						Name:     "Product 1",
						Code:     "CODE1",
						ImageURL: "http://example.com/1.jpg",
					},
					{
						ID:       2,
						Name:     "Product 2",
						Code:     "CODE2",
						ImageURL: "http://example.com/2.jpg",
					},
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "database error",
			mockGetAll: func() ([]*domain.Product, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := NewProductUseCase(mockRepo)

			products, err := useCase.GetAllProducts()
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if len(products) != 2 {
				t.Errorf("expected 2 products, got %d", len(products))
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {
	tests := []struct {
		name          string
		product       *domain.Product
		mockUpdate    func(*domain.Product) error
		expectedError error
	}{
		{
			name: "successful update",
			product: &domain.Product{
				ID:       1,
				Name:     "Updated Product",
				Code:     "UPD123",
				ImageURL: "http://example.com/updated.jpg",
			},
			mockUpdate: func(p *domain.Product) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "product not found",
			product: &domain.Product{
				ID:       999,
				Name:     "Non-existent Product",
				Code:     "NON123",
				ImageURL: "http://example.com/none.jpg",
			},
			mockUpdate: func(p *domain.Product) error {
				return &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedError: &domain.ProductNotFoundError{ProductID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				UpdateFunc: tt.mockUpdate,
			}
			useCase := NewProductUseCase(mockRepo)

			err := useCase.UpdateProduct(tt.product)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}
		})
	}
}

func TestDeleteProduct(t *testing.T) {
	tests := []struct {
		name          string
		productID     int64
		mockDelete    func(int64) error
		expectedError error
	}{
		{
			name:      "successful deletion",
			productID: 1,
			mockDelete: func(id int64) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:      "product not found",
			productID: 999,
			mockDelete: func(id int64) error {
				return &domain.ProductNotFoundError{ProductID: 999}
			},
			expectedError: &domain.ProductNotFoundError{ProductID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProductRepository{
				DeleteFunc: tt.mockDelete,
			}
			useCase := NewProductUseCase(mockRepo)

			err := useCase.DeleteProduct(tt.productID)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}
		})
	}
}
