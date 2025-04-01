package repository

import (
	"inventario/internal/domain"
)

type MockProductRepository struct {
	CreateFunc  func(*domain.Product) error
	GetByIDFunc func(int64) (*domain.Product, error)
	GetAllFunc  func() ([]*domain.Product, error)
	UpdateFunc  func(*domain.Product) error
	DeleteFunc  func(int64) error
	CloseFunc   func() error
}

func (m *MockProductRepository) Create(product *domain.Product) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(product)
	}
	return nil
}

func (m *MockProductRepository) GetByID(id int64) (*domain.Product, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockProductRepository) GetAll() ([]*domain.Product, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockProductRepository) Update(product *domain.Product) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(product)
	}
	return nil
}

func (m *MockProductRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockProductRepository) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}
