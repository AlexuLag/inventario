package usecase

import (
	"errors"
	"inventario/internal/domain"
	"inventario/internal/infrastructure/repository"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name          string
		userName      string
		userEmail     string
		userRole      string
		userPassword  string
		mockCreate    func(*domain.User) error
		expectedError error
	}{
		{
			name:         "successful creation",
			userName:     "Test User",
			userEmail:    "test@example.com",
			userRole:     "admin",
			userPassword: "password123",
			mockCreate: func(u *domain.User) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:         "user already exists",
			userName:     "Existing User",
			userEmail:    "existing@example.com",
			userRole:     "user",
			userPassword: "password123",
			mockCreate: func(u *domain.User) error {
				return &domain.UserAlreadyExistsError{Email: "existing@example.com"}
			},
			expectedError: &domain.UserAlreadyExistsError{Email: "existing@example.com"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := NewUserUseCase(mockRepo)

			user, err := useCase.CreateUser(tt.userName, tt.userEmail, tt.userRole, tt.userPassword)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if user.Name != tt.userName {
				t.Errorf("expected name %s, got %s", tt.userName, user.Name)
			}
			if user.Email != tt.userEmail {
				t.Errorf("expected email %s, got %s", tt.userEmail, user.Email)
			}
			if user.Role != tt.userRole {
				t.Errorf("expected role %s, got %s", tt.userRole, user.Role)
			}
			if user.Password != tt.userPassword {
				t.Errorf("expected password %s, got %s", tt.userPassword, user.Password)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockGetByID   func(int64) (*domain.User, error)
		expectedError error
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			mockGetByID: func(id int64) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Role:      "admin",
					Password:  "password123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: 999,
			mockGetByID: func(id int64) (*domain.User, error) {
				return nil, &domain.UserNotFoundError{UserID: 999}
			},
			expectedError: &domain.UserNotFoundError{UserID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := NewUserUseCase(mockRepo)

			user, err := useCase.GetUser(tt.userID)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if user.ID != tt.userID {
				t.Errorf("expected ID %d, got %d", tt.userID, user.ID)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	tests := []struct {
		name           string
		userEmail      string
		mockGetByEmail func(string) (*domain.User, error)
		expectedError  error
	}{
		{
			name:      "successful retrieval",
			userEmail: "test@example.com",
			mockGetByEmail: func(email string) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Role:      "admin",
					Password:  "password123",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name:      "user not found",
			userEmail: "nonexistent@example.com",
			mockGetByEmail: func(email string) (*domain.User, error) {
				return nil, &domain.UserNotFoundError{UserID: 0}
			},
			expectedError: &domain.UserNotFoundError{UserID: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetByEmailFunc: tt.mockGetByEmail,
			}
			useCase := NewUserUseCase(mockRepo)

			user, err := useCase.GetUserByEmail(tt.userEmail)
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if user.Email != tt.userEmail {
				t.Errorf("expected email %s, got %s", tt.userEmail, user.Email)
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name          string
		mockGetAll    func() ([]*domain.User, error)
		expectedError error
	}{
		{
			name: "successful retrieval",
			mockGetAll: func() ([]*domain.User, error) {
				return []*domain.User{
					{
						ID:        1,
						Name:      "User 1",
						Email:     "user1@example.com",
						Role:      "admin",
						Password:  "password123",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Name:      "User 2",
						Email:     "user2@example.com",
						Role:      "user",
						Password:  "password456",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "database error",
			mockGetAll: func() ([]*domain.User, error) {
				return nil, errors.New("database error")
			},
			expectedError: errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := NewUserUseCase(mockRepo)

			users, err := useCase.GetAllUsers()
			if err != nil {
				if tt.expectedError == nil {
					t.Errorf("unexpected error: %v", err)
				}
				if err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if len(users) != 2 {
				t.Errorf("expected 2 users, got %d", len(users))
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		user          *domain.User
		mockGetByID   func(int64) (*domain.User, error)
		mockUpdate    func(*domain.User) error
		expectedError error
	}{
		{
			name: "successful update",
			user: &domain.User{
				ID:        1,
				Name:      "Updated User",
				Email:     "updated@example.com",
				Role:      "admin",
				Password:  "newpassword",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockGetByID: func(id int64) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					Name:      "Original User",
					Email:     "original@example.com",
					Role:      "user",
					Password:  "oldpassword",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			mockUpdate: func(u *domain.User) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "user not found",
			user: &domain.User{
				ID:        999,
				Name:      "Non-existent User",
				Email:     "nonexistent@example.com",
				Role:      "user",
				Password:  "password123",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			mockGetByID: func(id int64) (*domain.User, error) {
				return nil, &domain.UserNotFoundError{UserID: 999}
			},
			mockUpdate: func(u *domain.User) error {
				return nil
			},
			expectedError: &domain.UserNotFoundError{UserID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetByIDFunc: tt.mockGetByID,
				UpdateFunc:  tt.mockUpdate,
			}
			useCase := NewUserUseCase(mockRepo)

			err := useCase.UpdateUser(tt.user)
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

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        int64
		mockDelete    func(int64) error
		expectedError error
	}{
		{
			name:   "successful deletion",
			userID: 1,
			mockDelete: func(id int64) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name:   "user not found",
			userID: 999,
			mockDelete: func(id int64) error {
				return &domain.UserNotFoundError{UserID: 999}
			},
			expectedError: &domain.UserNotFoundError{UserID: 999},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				DeleteFunc: tt.mockDelete,
			}
			useCase := NewUserUseCase(mockRepo)

			err := useCase.DeleteUser(tt.userID)
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
