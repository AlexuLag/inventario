package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type SQLiteStockRepository struct {
	db *sql.DB
}

func NewSQLiteStockRepository(db *sql.DB) *SQLiteStockRepository {
	return &SQLiteStockRepository{db: db}
}

func (r *SQLiteStockRepository) Create(stock *domain.Stock) error {
	result, err := r.db.Exec(`
		INSERT INTO stocks (
			product_id, serial, created_at, updated_at, 
			created_by_user_id, updated_by_user_id, 
			batch, purchase_date, provider_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, stock.Product.ID, stock.Serial, stock.CreatedAt, stock.UpdatedAt,
		stock.CreatedByUser.ID, stock.UpdatedByUser.ID,
		stock.Batch, stock.PurchaseDate, stock.Provider.ID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	stock.ID = id
	return nil
}

func (r *SQLiteStockRepository) GetByID(id int64) (*domain.Stock, error) {
	var stock domain.Stock
	var productID, createdByUserID, updatedByUserID, providerID int64

	err := r.db.QueryRow(`
		SELECT id, product_id, serial, created_at, updated_at,
			created_by_user_id, updated_by_user_id,
			batch, purchase_date, provider_id
		FROM stocks
		WHERE id = ?
	`, id).Scan(
		&stock.ID, &productID, &stock.Serial, &stock.CreatedAt, &stock.UpdatedAt,
		&createdByUserID, &updatedByUserID,
		&stock.Batch, &stock.PurchaseDate, &providerID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch related entities
	stock.Product = &domain.Product{ID: productID}
	stock.CreatedByUser = &domain.User{ID: createdByUserID}
	stock.UpdatedByUser = &domain.User{ID: updatedByUserID}
	stock.Provider = &domain.Provider{ID: providerID}

	return &stock, nil
}

func (r *SQLiteStockRepository) GetAll() ([]domain.Stock, error) {
	rows, err := r.db.Query(`
		SELECT id, product_id, serial, created_at, updated_at,
			created_by_user_id, updated_by_user_id,
			batch, purchase_date, provider_id
		FROM stocks
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		var productID, createdByUserID, updatedByUserID, providerID int64

		err := rows.Scan(
			&stock.ID, &productID, &stock.Serial, &stock.CreatedAt, &stock.UpdatedAt,
			&createdByUserID, &updatedByUserID,
			&stock.Batch, &stock.PurchaseDate, &providerID,
		)
		if err != nil {
			return nil, err
		}

		// Set related entities
		stock.Product = &domain.Product{ID: productID}
		stock.CreatedByUser = &domain.User{ID: createdByUserID}
		stock.UpdatedByUser = &domain.User{ID: updatedByUserID}
		stock.Provider = &domain.Provider{ID: providerID}

		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *SQLiteStockRepository) Update(stock *domain.Stock) error {
	result, err := r.db.Exec(`
		UPDATE stocks
		SET product_id = ?, serial = ?, updated_at = CURRENT_TIMESTAMP,
			updated_by_user_id = ?, batch = ?, purchase_date = ?, provider_id = ?
		WHERE id = ?
	`, stock.Product.ID, stock.Serial, stock.UpdatedByUser.ID,
		stock.Batch, stock.PurchaseDate, stock.Provider.ID, stock.ID)
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

func (r *SQLiteStockRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM stocks WHERE id = ?", id)
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

func (r *SQLiteStockRepository) GetByProductID(productID int64) ([]domain.Stock, error) {
	rows, err := r.db.Query(`
		SELECT id, product_id, serial, created_at, updated_at,
			created_by_user_id, updated_by_user_id,
			batch, purchase_date, provider_id
		FROM stocks
		WHERE product_id = ?
	`, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		var productID, createdByUserID, updatedByUserID, providerID int64

		err := rows.Scan(
			&stock.ID, &productID, &stock.Serial, &stock.CreatedAt, &stock.UpdatedAt,
			&createdByUserID, &updatedByUserID,
			&stock.Batch, &stock.PurchaseDate, &providerID,
		)
		if err != nil {
			return nil, err
		}

		// Set related entities
		stock.Product = &domain.Product{ID: productID}
		stock.CreatedByUser = &domain.User{ID: createdByUserID}
		stock.UpdatedByUser = &domain.User{ID: updatedByUserID}
		stock.Provider = &domain.Provider{ID: providerID}

		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *SQLiteStockRepository) GetBySerial(serial string) (*domain.Stock, error) {
	var stock domain.Stock
	var productID, createdByUserID, updatedByUserID, providerID int64

	err := r.db.QueryRow(`
		SELECT id, product_id, serial, created_at, updated_at,
			created_by_user_id, updated_by_user_id,
			batch, purchase_date, provider_id
		FROM stocks
		WHERE serial = ?
	`, serial).Scan(
		&stock.ID, &productID, &stock.Serial, &stock.CreatedAt, &stock.UpdatedAt,
		&createdByUserID, &updatedByUserID,
		&stock.Batch, &stock.PurchaseDate, &providerID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Set related entities
	stock.Product = &domain.Product{ID: productID}
	stock.CreatedByUser = &domain.User{ID: createdByUserID}
	stock.UpdatedByUser = &domain.User{ID: updatedByUserID}
	stock.Provider = &domain.Provider{ID: providerID}

	return &stock, nil
}

func (r *SQLiteStockRepository) Close() error {
	return r.db.Close()
}
