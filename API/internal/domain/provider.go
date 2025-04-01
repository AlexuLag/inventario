package domain

import "time"

type Provider struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProviderNotFoundError struct {
	ProviderID int64
}

func (e *ProviderNotFoundError) Error() string {
	return "provider not found"
}

type ProviderAlreadyExistsError struct {
	Email string
}

func (e *ProviderAlreadyExistsError) Error() string {
	return "provider already exists"
}

type IProviderRepository interface {
	Create(provider *Provider) error
	GetByID(id int64) (*Provider, error)
	GetAll() ([]Provider, error)
	Update(provider *Provider) error
	Delete(id int64) error
}
