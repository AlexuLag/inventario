package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Password  string    `json:"-"` // The "-" tag ensures the password is never sent in JSON responses
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	GetByEmail(email string) (*User, error)
	GetAll() ([]*User, error)
	Update(user *User) error
	Delete(id int64) error
	Close() error
}

// UserNotFoundError represents an error when a user is not found
type UserNotFoundError struct {
	UserID int64
}

func (e *UserNotFoundError) Error() string {
	return "user not found"
}

// UserAlreadyExistsError represents an error when a user already exists
type UserAlreadyExistsError struct {
	Email string
}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}
