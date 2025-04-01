package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type MySQLUserRepository struct {
	*MySQLBaseRepository
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		MySQLBaseRepository: NewMySQLBaseRepository(db),
	}
}

func (r *MySQLUserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password, role, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		now,
		now,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.UserAlreadyExistsError{Email: user.Email}
		}
		return err
	}

	id, err := r.GetLastInsertID(result)
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *MySQLUserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	var user domain.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MySQLUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	var user domain.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *MySQLUserRepository) GetAll() ([]domain.User, error) {
	query := `
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *MySQLUserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, password = ?, role = ?, updated_at = ?
		WHERE id = ?
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		now,
		user.ID,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.UserAlreadyExistsError{Email: user.Email}
		}
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.UserNotFoundError{UserID: user.ID}
	}

	user.UpdatedAt = now
	return nil
}

func (r *MySQLUserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.UserNotFoundError{UserID: id}
	}

	return nil
}
