package repository

import (
	"database/sql"
	"inventario/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteProductRepository struct {
	db *sql.DB
}

func NewSQLiteProductRepository(dbPath string) (*SQLiteProductRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create products table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			code TEXT NOT NULL UNIQUE,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			image_url TEXT
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteProductRepository{db: db}, nil
}

func (r *SQLiteProductRepository) Create(product *domain.Product) error {
	// Check if product with same code exists
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE code = ?)", product.Code).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return &domain.ProductAlreadyExistsError{
			Code: product.Code,
		}
	}

	result, err := r.db.Exec(`
		INSERT INTO products (name, code, created_at, updated_at, image_url)
		VALUES (?, ?, ?, ?, ?)
	`,
		product.Name,
		product.Code,
		product.CreatedAt,
		product.UpdatedAt,
		product.ImageURL,
	)
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
		SELECT id, name, code, created_at, updated_at, image_url
		FROM products
		WHERE id = ?
	`, id).Scan(
		&product.ID,
		&product.Name,
		&product.Code,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.ImageURL,
	)

	if err == sql.ErrNoRows {
		return nil, &domain.ProductNotFoundError{
			ProductID: id,
		}
	}
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *SQLiteProductRepository) GetAll() ([]*domain.Product, error) {
	rows, err := r.db.Query(`
		SELECT id, name, code, created_at, updated_at, image_url
		FROM products
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Code,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.ImageURL,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	return products, nil
}

func (r *SQLiteProductRepository) Update(product *domain.Product) error {
	result, err := r.db.Exec(`
		UPDATE products
		SET name = ?, code = ?, updated_at = ?, image_url = ?
		WHERE id = ?
	`,
		product.Name,
		product.Code,
		product.UpdatedAt,
		product.ImageURL,
		product.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &domain.ProductNotFoundError{
			ProductID: product.ID,
		}
	}

	return nil
}

func (r *SQLiteProductRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return &domain.ProductNotFoundError{
			ProductID: id,
		}
	}

	return nil
}

func (r *SQLiteProductRepository) Close() error {
	return r.db.Close()
}
