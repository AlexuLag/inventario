package usecase

import (
	"errors"
	"inventario/internal/domain"
	"inventario/internal/infrastructure/repository"
	"testing"
	"time"
)

func TestCreateProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerName  string
		providerEmail string
		providerPhone string
		providerAddr  string
		mockCreate    func(*domain.Provider) error
		expectedError error
	}{
		{
			name:          "successful creation",
			providerName:  "Test Provider",
			providerEmail: "test@example.com",
			providerPhone: "1234567890",
			providerAddr:  "Test Address",
			mockCreate: func(p *domain.Provider) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:          "provider already exists",
			providerName:  "Existing Provider",
			providerEmail: "existing@example.com",
			providerPhone: "0987654321",
			providerAddr:  "Existing Address",
			mockCreate: func(p *domain.Provider) error {
				return &domain.ProviderAlreadyExistsError{Email: "existing@example.com"}
			},
			expectedError: &domain.ProviderAlreadyExistsError{Email: "existing@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := NewProviderUseCase(mockRepo)

			provider, err := useCase.CreateProvider(tt.providerName, tt.providerEmail, tt.providerPhone, tt.providerAddr)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if provider.Name != tt.providerName {
				t.Errorf("expected name %s, got %s", tt.providerName, provider.Name)
			}
			if provider.Email != tt.providerEmail {
				t.Errorf("expected email %s, got %s", tt.providerEmail, provider.Email)
			}
			if provider.Phone != tt.providerPhone {
				t.Errorf("expected phone %s, got %s", tt.providerPhone, provider.Phone)
			}
			if provider.Address != tt.providerAddr {
				t.Errorf("expected address %s, got %s", tt.providerAddr, provider.Address)
			}
		})
	}
}

func TestGetProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerID    int64
		mockGetByID   func(int64) (*domain.Provider, error)
		expectedError error
	}{
		{
			name:       "successful retrieval",
			providerID: 1,
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
			expectedError: nil,
		},
		{
			name:       "provider not found",
			providerID: 999,
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			expectedError: &domain.ProviderNotFoundError{ProviderID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := NewProviderUseCase(mockRepo)

			provider, err := useCase.GetProvider(tt.providerID)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if provider.ID != tt.providerID {
				t.Errorf("expected ID %d, got %d", tt.providerID, provider.ID)
			}
		})
	}
}

func TestGetAllProviders(t *testing.T) {
	tests := []struct {
		name          string
		mockGetAll    func() ([]domain.Provider, error)
		expectedError error
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
			expectedError: nil,
		},
		{
			name: "database error",
			mockGetAll: func() ([]domain.Provider, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := NewProviderUseCase(mockRepo)

			providers, err := useCase.GetAllProviders()
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if len(providers) != 2 {
				t.Errorf("expected 2 providers, got %d", len(providers))
			}
		})
	}
}

func TestUpdateProvider(t *testing.T) {
	tests := []struct {
		name          string
		provider      *domain.Provider
		mockGetByID   func(int64) (*domain.Provider, error)
		mockUpdate    func(*domain.Provider) error
		expectedError error
	}{
		{
			name: "successful update",
			provider: &domain.Provider{
				ID:        1,
				Name:      "Updated Provider",
				Email:     "updated@example.com",
				Phone:     "1234567890",
				Address:   "Updated Address",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
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
			expectedError: nil,
		},
		{
			name: "provider not found",
			provider: &domain.Provider{
				ID:        999,
				Name:      "Non-existent Provider",
				Email:     "nonexistent@example.com",
				Phone:     "1234567890",
				Address:   "Non-existent Address",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			mockUpdate: func(p *domain.Provider) error {
				return nil
			},
			expectedError: &domain.ProviderNotFoundError{ProviderID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
				UpdateFunc:  tt.mockUpdate,
			}
			useCase := NewProviderUseCase(mockRepo)

			err := useCase.UpdateProvider(tt.provider)
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

func TestDeleteProvider(t *testing.T) {
	tests := []struct {
		name          string
		providerID    int64
		mockGetByID   func(int64) (*domain.Provider, error)
		mockDelete    func(int64) error
		expectedError error
	}{
		{
			name:       "successful deletion",
			providerID: 1,
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
			expectedError: nil,
		},
		{
			name:       "provider not found",
			providerID: 999,
			mockGetByID: func(id int64) (*domain.Provider, error) {
				return nil, &domain.ProviderNotFoundError{ProviderID: 999}
			},
			mockDelete: func(id int64) error {
				return nil
			},
			expectedError: &domain.ProviderNotFoundError{ProviderID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockProviderRepository{
				GetByIDFunc: tt.mockGetByID,
				DeleteFunc:  tt.mockDelete,
			}
			useCase := NewProviderUseCase(mockRepo)

			err := useCase.DeleteProvider(tt.providerID)
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
