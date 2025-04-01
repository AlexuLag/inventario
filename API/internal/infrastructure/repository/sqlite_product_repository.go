package repository

import (
	"database/sql"
	"inventario/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(db *sql.DB) *SQLiteProductRepository {
	return &SQLiteProductRepository{db: db}
}

func (r *SQLiteProductRepository) Create(product *domain.Product) error {
	result, err := r.db.Exec(`
		INSERT INTO products (name, code, image_url)
		VALUES (?, ?, ?)
	`, product.Name, product.Code, product.ImageURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	product.ID = id
	return nil
}

func (r *SQLiteProductRepository) GetByID(id int64) (*domain.Product, error) {
	var product domain.Product
	err := r.db.QueryRow(`
		SELECT id, name, code, image_url, created_at, updated_at
		FROM products
		WHERE id = ?
	`, id).Scan(&product.ID, &product.Name, &product.Code, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *SQLiteProductRepository) GetAll() ([]domain.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, code, image_url, created_at, updated_at
		FROM products
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Code, &product.ImageURL, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func (r *SQLiteProductRepository) Update(product *domain.Product) error {
	result, err := r.db.Exec(`
		UPDATE products
		SET name = ?, code = ?, image_url = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`, product.Name, product.Code, product.ImageURL, product.ID)
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

func (r *SQLiteProductRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
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

func (r *SQLiteProductRepository) Close() error {
	return r.db.Close()
}
