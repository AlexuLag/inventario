package usecase

import (
	"inventario/internal/domain"
	"time"
)

type UserUseCase struct {
	userRepo domain.IUserRepository
}

func NewUserUseCase(userRepo domain.IUserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) CreateUser(name, email, role, password string) (*domain.User, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return nil, &domain.UserAlreadyExistsError{Email: email}
	}

	now := time.Now()
	user := &domain.User{
		Name:      name,
		Email:     email,
		Role:      role,
		Password:  password, // Note: In a real application, this should be hashed
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) GetUser(id int64) (*domain.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, &domain.UserNotFoundError{UserID: id}
	}
	return user, nil
}

func (u *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, &domain.UserNotFoundError{UserID: 0} // We don't have the ID in this case
	}
	return user, nil
}

func (u *UserUseCase) GetAllUsers() ([]*domain.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.User, len(users))
	for i := range users {
		result[i] = &users[i]
	}
	return result, nil
}

func (u *UserUseCase) UpdateUser(user *domain.User) error {
	// Check if user exists
	existingUser, err := u.userRepo.GetByID(user.ID)
	if err != nil {
		return &domain.UserNotFoundError{UserID: user.ID}
	}

	// Update only allowed fields
	existingUser.Name = user.Name
	existingUser.Email = user.Email
	existingUser.Role = user.Role
	existingUser.UpdatedAt = time.Now()

	if user.Password != "" {
		existingUser.Password = user.Password // Note: In a real application, this should be hashed
	}

	return u.userRepo.Update(existingUser)
}

func (u *UserUseCase) DeleteUser(id int64) error {
	// Check if user exists
	_, err := u.userRepo.GetByID(id)
	if err != nil {
		return &domain.UserNotFoundError{UserID: id}
	}

	return u.userRepo.Delete(id)
}
