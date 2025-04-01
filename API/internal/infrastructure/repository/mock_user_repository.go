package repository

import (
	"inventario/internal/domain"
)

type MockUserRepository struct {
	CreateFunc     func(*domain.User) error
	GetByIDFunc    func(int64) (*domain.User, error)
	GetByEmailFunc func(string) (*domain.User, error)
	GetAllFunc     func() ([]*domain.User, error)
	UpdateFunc     func(*domain.User) error
	DeleteFunc     func(int64) error
}

func (m *MockUserRepository) Create(user *domain.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(id int64) (*domain.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserRepository) GetAll() ([]domain.User, error) {
	if m.GetAllFunc != nil {
		users, err := m.GetAllFunc()
		if err != nil {
			return nil, err
		}
		result := make([]domain.User, len(users))
		for i, u := range users {
			result[i] = *u
		}
		return result, nil
	}
	return nil, nil
}

func (m *MockUserRepository) Update(user *domain.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(user)
	}
	return nil
}

func (m *MockUserRepository) Delete(id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}
