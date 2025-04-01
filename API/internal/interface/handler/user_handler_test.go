package handler

import (
	"bytes"
	"context"
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

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockCreate     func(*domain.User) error
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful creation",
			requestBody: map[string]string{
				"name":     "Test User",
				"email":    "test@example.com",
				"role":     "admin",
				"password": "password123",
			},
			mockCreate: func(u *domain.User) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			expectedError:  "",
		},
		{
			name: "user already exists",
			requestBody: map[string]string{
				"name":     "Existing User",
				"email":    "existing@example.com",
				"role":     "user",
				"password": "password123",
			},
			mockCreate: func(u *domain.User) error {
				return &domain.UserAlreadyExistsError{Email: "existing@example.com"}
			},
			expectedStatus: http.StatusConflict,
			expectedError:  "user with email existing@example.com already exists",
		},
		{
			name: "invalid request body",
			requestBody: map[string]string{
				"name": "Invalid User",
				// Missing required fields
			},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name, email, role and password are required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				CreateFunc: tt.mockCreate,
			}
			useCase := usecase.NewUserUseCase(mockRepo)
			handler := NewUserHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(body))
			w := httptest.NewRecorder()

			handler.CreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockGetByID    func(int64) (*domain.User, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful retrieval",
			userID: "1",
			mockGetByID: func(id int64) (*domain.User, error) {
				if id == 1 {
					return &domain.User{
						ID:        1,
						Name:      "Test User",
						Email:     "test@example.com",
						Role:      "admin",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					}, nil
				}
				return nil, &domain.UserNotFoundError{UserID: id}
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:   "user not found",
			userID: "999",
			mockGetByID: func(id int64) (*domain.User, error) {
				return nil, &domain.UserNotFoundError{UserID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user with ID 999 not found",
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			mockGetByID:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := usecase.NewUserUseCase(mockRepo)
			handler := NewUserHandler(useCase)

			req := httptest.NewRequest("GET", "/api/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.GetUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

		})
	}
}

func TestGetAllUsers(t *testing.T) {
	tests := []struct {
		name           string
		mockGetAll     func() ([]*domain.User, error)
		expectedStatus int
		expectedError  string
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
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        2,
						Name:      "User 2",
						Email:     "user2@example.com",
						Role:      "user",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "database error",
			mockGetAll: func() ([]*domain.User, error) {
				return nil, domain.ErrInternalServer
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				GetAllFunc: tt.mockGetAll,
			}
			useCase := usecase.NewUserUseCase(mockRepo)
			handler := NewUserHandler(useCase)

			req := httptest.NewRequest("GET", "/api/users", nil)
			w := httptest.NewRecorder()

			handler.GetAllUsers(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedError != "" {
				if w.Body.String() != tt.expectedError+"\n" {
					t.Errorf("expected error %s, got %s", tt.expectedError, w.Body.String())
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    map[string]string
		mockUpdate     func(*domain.User) error
		mockGetByID    func(id int64) (*domain.User, error)
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful update",
			userID: "1",
			requestBody: map[string]string{
				"name":  "Updated User",
				"email": "updated@example.com",
				"role":  "admin",
			},
			mockUpdate: func(u *domain.User) error {
				return nil
			},
			mockGetByID: func(id int64) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Role:      "admin",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:   "user not found",
			userID: "999",
			requestBody: map[string]string{
				"name":  "Non-existent User",
				"email": "nonexistent@example.com",
				"role":  "user",
			},
			mockUpdate: func(u *domain.User) error {
				return &domain.UserNotFoundError{UserID: 999}
			},
			mockGetByID: func(id int64) (*domain.User, error) {
				return &domain.User{
					ID:        1,
					Name:      "Test User",
					Email:     "test@example.com",
					Role:      "admin",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user with ID 999 not found",
		},
		{
			name:   "invalid request body",
			userID: "1",
			requestBody: map[string]string{
				"name": "Invalid User",
				// Missing required fields
			},
			mockUpdate:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name, email and role are required fields",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				UpdateFunc:  tt.mockUpdate,
				GetByIDFunc: tt.mockGetByID,
			}
			useCase := usecase.NewUserUseCase(mockRepo)
			handler := NewUserHandler(useCase)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("PUT", "/api/users/"+tt.userID, bytes.NewBuffer(body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.UpdateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockDelete     func(int64) error
		expectedStatus int
		expectedError  string
	}{
		{
			name:   "successful deletion",
			userID: "1",
			mockDelete: func(id int64) error {
				return nil
			},
			expectedStatus: http.StatusNoContent,
			expectedError:  "",
		},
		{
			name:   "user not found",
			userID: "999",
			mockDelete: func(id int64) error {
				return &domain.UserNotFoundError{UserID: 999}
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user with ID 999 not found",
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			mockDelete:     nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repository.MockUserRepository{
				DeleteFunc: tt.mockDelete,
			}
			useCase := usecase.NewUserUseCase(mockRepo)
			handler := NewUserHandler(useCase)

			req := httptest.NewRequest("DELETE", "/api/users/"+tt.userID, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			handler.DeleteUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

		})
	}
}
