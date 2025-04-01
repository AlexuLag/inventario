package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type MySQLProductRepository struct {
	*MySQLBaseRepository
}

func NewMySQLProductRepository(db *sql.DB) *MySQLProductRepository {
	return &MySQLProductRepository{
		MySQLBaseRepository: NewMySQLBaseRepository(db),
	}
}

func (r *MySQLProductRepository) Create(product *domain.Product) error {
	query := `
		INSERT INTO products (name, code, image_url, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		product.Name,
		product.Code,
		product.ImageURL,
		now,
		now,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.ProductAlreadyExistsError{Code: product.Code}
		}
		return err
	}

	id, err := r.GetLastInsertID(result)
	if err != nil {
		return err
	}

	product.ID = id
	product.CreatedAt = now
	product.UpdatedAt = now
	return nil
}

func (r *MySQLProductRepository) GetByID(id int64) (*domain.Product, error) {
	query := `
		SELECT id, name, code, image_url, created_at, updated_at
		FROM products
		WHERE id = ?
	`

	var product domain.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Code,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *MySQLProductRepository) GetAll() ([]domain.Product, error) {
	query := `
		SELECT id, name, code, image_url, created_at, updated_at
		FROM products
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Code,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *MySQLProductRepository) Update(product *domain.Product) error {
	query := `
		UPDATE products
		SET name = ?, code = ?, image_url = ?, updated_at = ?
		WHERE id = ?
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		product.Name,
		product.Code,
		product.ImageURL,
		now,
		product.ID,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.ProductAlreadyExistsError{Code: product.Code}
		}
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.ProductNotFoundError{ProductID: product.ID}
	}

	product.UpdatedAt = now
	return nil
}

func (r *MySQLProductRepository) Delete(id int64) error {
	query := "DELETE FROM products WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.ProductNotFoundError{ProductID: id}
	}

	return nil
}
