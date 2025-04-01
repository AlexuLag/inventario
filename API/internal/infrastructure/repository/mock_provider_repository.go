package repository

import (
	"inventario/internal/domain"
)

type MockProviderRepository struct {
	CreateFunc  func(*domain.Provider) error
	GetByIDFunc func(int64) (*domain.Provider, error)
	GetAllFunc  func() ([]domain.Provider, error)
	UpdateFunc  func(*domain.Provider) error
	DeleteFunc  func(int64) error
}

func (m *MockProviderRepository) Create(provider *domain.Provider) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(provider)
	}
	return nil
}

func (m *MockProviderRepository) GetByID(id int64) (*domain.Provider, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockProviderRepository) GetAll() ([]domain.Provider, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockProviderRepository) Update(provider *domain.Provider) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(provider)
	}
	return nil
}

func (m *MockProviderRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
