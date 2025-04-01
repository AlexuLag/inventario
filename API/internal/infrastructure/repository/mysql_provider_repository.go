package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type MySQLProviderRepository struct {
	*MySQLBaseRepository
}

func NewMySQLProviderRepository(db *sql.DB) *MySQLProviderRepository {
	return &MySQLProviderRepository{
		MySQLBaseRepository: NewMySQLBaseRepository(db),
	}
}

func (r *MySQLProviderRepository) Create(provider *domain.Provider) error {
	query := `
		INSERT INTO providers (name, email, phone, address, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Address,
		now,
		now,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.ProviderAlreadyExistsError{Email: provider.Email}
		}
		return err
	}

	id, err := r.GetLastInsertID(result)
	if err != nil {
		return err
	}

	provider.ID = id
	provider.CreatedAt = now
	provider.UpdatedAt = now
	return nil
}

func (r *MySQLProviderRepository) GetByID(id int64) (*domain.Provider, error) {
	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM providers
		WHERE id = ?
	`

	var provider domain.Provider
	err := r.db.QueryRow(query, id).Scan(
		&provider.ID,
		&provider.Name,
		&provider.Email,
		&provider.Phone,
		&provider.Address,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &provider, nil
}

func (r *MySQLProviderRepository) GetAll() ([]domain.Provider, error) {
	query := `
		SELECT id, name, email, phone, address, created_at, updated_at
		FROM providers
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []domain.Provider
	for rows.Next() {
		var provider domain.Provider
		err := rows.Scan(
			&provider.ID,
			&provider.Name,
			&provider.Email,
			&provider.Phone,
			&provider.Address,
			&provider.CreatedAt,
			&provider.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, provider)
	}
	return providers, nil
}

func (r *MySQLProviderRepository) Update(provider *domain.Provider) error {
	query := `
		UPDATE providers
		SET name = ?, email = ?, phone = ?, address = ?, updated_at = ?
		WHERE id = ?
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		provider.Name,
		provider.Email,
		provider.Phone,
		provider.Address,
		now,
		provider.ID,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.ProviderAlreadyExistsError{Email: provider.Email}
		}
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.ProviderNotFoundError{ProviderID: provider.ID}
	}

	provider.UpdatedAt = now
	return nil
}

func (r *MySQLProviderRepository) Delete(id int64) error {
	query := "DELETE FROM providers WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.ProviderNotFoundError{ProviderID: id}
	}

	return nil
}
