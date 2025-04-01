package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(user *domain.User) error {
	// Check if user with same email exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return &domain.UserAlreadyExistsError{
			Email: user.Email,
		}
	}

	result, err := r.db.Exec(`
		INSERT INTO users (name, email, password, role)
		VALUES (?, ?, ?, ?)
	`, user.Name, user.Email, user.Password, user.Role)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *SQLiteUserRepository) GetByID(id int64) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.QueryRow(`
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
		WHERE email = ?
	`, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLiteUserRepository) GetAll() ([]domain.User, error) {
	rows, err := r.db.Query(`
		SELECT id, name, email, password, role, created_at, updated_at
		FROM users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *SQLiteUserRepository) Update(user *domain.User) error {
	result, err := r.db.Exec(`
		UPDATE users
		SET name = ?, email = ?, password = ?, role = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, user.Name, user.Email, user.Password, user.Role, user.ID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteUserRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteUserRepository) Close() error {
	return r.db.Close()
}
