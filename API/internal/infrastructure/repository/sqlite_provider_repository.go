package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type SQLiteProviderRepository struct {
	db *sql.DB
}

func NewSQLiteProviderRepository(db *sql.DB) *SQLiteProviderRepository {
	return &SQLiteProviderRepository{db: db}
}

func (r *SQLiteProviderRepository) Create(provider *domain.Provider) error {
	result, err := r.db.Exec(`
		INSERT INTO providers (name, address, phone)
		VALUES (?, ?, ?)
	`, provider.Name, provider.Address, provider.Phone)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	provider.ID = id
	return nil
}

func (r *SQLiteProviderRepository) GetByID(id int64) (*domain.Provider, error) {
	var provider domain.Provider
	err := r.db.QueryRow(`
		SELECT id, name, address, phone
		FROM providers
		WHERE id = ?
	`, id).Scan(&provider.ID, &provider.Name, &provider.Address, &provider.Phone)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

func (r *SQLiteProviderRepository) GetAll() ([]domain.Provider, error) {
	rows, err := r.db.Query(`
		SELECT id, name, address, phone
		FROM providers
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []domain.Provider
	for rows.Next() {
		var provider domain.Provider
		err := rows.Scan(&provider.ID, &provider.Name, &provider.Address, &provider.Phone)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (r *SQLiteProviderRepository) Update(provider *domain.Provider) error {
	result, err := r.db.Exec(`
		UPDATE providers
		SET name = ?, address = ?, phone = ?
		WHERE id = ?
	`, provider.Name, provider.Address, provider.Phone, provider.ID)
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

func (r *SQLiteProviderRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM providers WHERE id = ?", id)
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

func (r *SQLiteProviderRepository) Close() error {
	return r.db.Close()
}
