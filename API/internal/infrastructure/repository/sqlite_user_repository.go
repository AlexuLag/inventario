package repository

import (
	"database/sql"
	"inventario/internal/domain"
	"time"
)

type SQLiteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, role, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Role,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
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
	query := `
		SELECT id, name, email, role, password, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, &domain.UserNotFoundError{UserID: id}
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, role, password, created_at, updated_at
		FROM users
		WHERE email = ?
	`
	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Role,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, &domain.UserNotFoundError{UserID: 0}
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *SQLiteUserRepository) GetAll() ([]*domain.User, error) {
	query := `
		SELECT id, name, email, role, password, created_at, updated_at
		FROM users
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Role,
			&user.Password,
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

func (r *SQLiteUserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, role = ?, password = ?, updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.Exec(query,
		user.Name,
		user.Email,
		user.Role,
		user.Password,
		time.Now(),
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &domain.UserNotFoundError{UserID: user.ID}
	}

	return nil
}

func (r *SQLiteUserRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &domain.UserNotFoundError{UserID: id}
	}

	return nil
}

func (r *SQLiteUserRepository) Close() error {
	return r.db.Close()
}
