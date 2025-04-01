package repository

import (
	"inventario/internal/domain"
)

type MockStockRepository struct {
	CreateFunc         func(*domain.Stock) error
	GetByIDFunc        func(int64) (*domain.Stock, error)
	GetAllFunc         func() ([]domain.Stock, error)
	GetByProductIDFunc func(int64) ([]domain.Stock, error)
	GetBySerialFunc    func(string) (*domain.Stock, error)
	UpdateFunc         func(*domain.Stock) error
	DeleteFunc         func(int64) error
}

func (m *MockStockRepository) Create(stock *domain.Stock) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(stock)
	}
	return nil
}

func (m *MockStockRepository) GetByID(id int64) (*domain.Stock, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockStockRepository) GetAll() ([]domain.Stock, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc()
	}
	return nil, nil
}

func (m *MockStockRepository) GetByProductID(productID int64) ([]domain.Stock, error) {
	if m.GetByProductIDFunc != nil {
		return m.GetByProductIDFunc(productID)
	}
	return nil, nil
}

func (m *MockStockRepository) GetBySerial(serial string) (*domain.Stock, error) {
	if m.GetBySerialFunc != nil {
		return m.GetBySerialFunc(serial)
	}
	return nil, nil
}

func (m *MockStockRepository) Update(stock *domain.Stock) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(stock)
	}
	return nil
}

func (m *MockStockRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
