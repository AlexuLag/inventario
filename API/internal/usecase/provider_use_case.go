package usecase

import (
	"inventario/internal/domain"
)

type ProviderUseCase struct {
	providerRepo domain.IProviderRepository
}

func NewProviderUseCase(repo domain.IProviderRepository) *ProviderUseCase {
	return &ProviderUseCase{
		providerRepo: repo,
	}
}

func (u *ProviderUseCase) CreateProvider(name, email, phone, address string) (*domain.Provider, error) {
	provider := &domain.Provider{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Address: address,
	}

	if err := u.providerRepo.Create(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (u *ProviderUseCase) GetProvider(id int64) (*domain.Provider, error) {
	provider, err := u.providerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if provider == nil {
		return nil, &domain.ProviderNotFoundError{ProviderID: id}
	}

	return provider, nil
}

func (u *ProviderUseCase) GetAllProviders() ([]domain.Provider, error) {
	return u.providerRepo.GetAll()
}

func (u *ProviderUseCase) UpdateProvider(provider *domain.Provider) error {
	existingProvider, err := u.providerRepo.GetByID(provider.ID)
	if err != nil {
		return err
	}

	if existingProvider == nil {
		return &domain.ProviderNotFoundError{ProviderID: provider.ID}
	}

	return u.providerRepo.Update(provider)
}

func (u *ProviderUseCase) DeleteProvider(id int64) error {
	provider, err := u.providerRepo.GetByID(id)
	if err != nil {
		return err
	}

	if provider == nil {
		return &domain.ProviderNotFoundError{ProviderID: id}
	}

	return u.providerRepo.Delete(id)
}
