package repository

import (
	"database/sql"
	"inventario/internal/domain"
)

type MySQLStockRepository struct {
	*MySQLBaseRepository
}

func NewMySQLStockRepository(db *sql.DB) *MySQLStockRepository {
	return &MySQLStockRepository{
		MySQLBaseRepository: NewMySQLBaseRepository(db),
	}
}

func (r *MySQLStockRepository) Create(stock *domain.Stock) error {
	query := `
		INSERT INTO stocks (
			product_id, serial, created_at, updated_at,
			created_by_user_id, updated_by_user_id,
			batch, purchase_date, provider_id
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		stock.Product.ID,
		stock.Serial,
		now,
		now,
		stock.CreatedByUser.ID,
		stock.UpdatedByUser.ID,
		stock.Batch,
		stock.PurchaseDate,
		stock.Provider.ID,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.StockAlreadyExistsError{Serial: stock.Serial}
		}
		return err
	}

	id, err := r.GetLastInsertID(result)
	if err != nil {
		return err
	}

	stock.ID = id
	stock.CreatedAt = now
	stock.UpdatedAt = now
	return nil
}

func (r *MySQLStockRepository) GetByID(id int64) (*domain.Stock, error) {
	query := `
		SELECT 
			s.id, s.serial, s.created_at, s.updated_at,
			s.batch, s.purchase_date,
			p.id, p.name, p.code, p.image_url,
			u1.id, u1.name, u1.email, u1.role,
			u2.id, u2.name, u2.email, u2.role,
			pr.id, pr.name, pr.email, pr.phone, pr.address
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN users u1 ON s.created_by_user_id = u1.id
		JOIN users u2 ON s.updated_by_user_id = u2.id
		JOIN providers pr ON s.provider_id = pr.id
		WHERE s.id = ?
	`

	var stock domain.Stock
	err := r.db.QueryRow(query, id).Scan(
		&stock.ID,
		&stock.Serial,
		&stock.CreatedAt,
		&stock.UpdatedAt,
		&stock.Batch,
		&stock.PurchaseDate,
		&stock.Product.ID,
		&stock.Product.Name,
		&stock.Product.Code,
		&stock.Product.ImageURL,
		&stock.CreatedByUser.ID,
		&stock.CreatedByUser.Name,
		&stock.CreatedByUser.Email,
		&stock.CreatedByUser.Role,
		&stock.UpdatedByUser.ID,
		&stock.UpdatedByUser.Name,
		&stock.UpdatedByUser.Email,
		&stock.UpdatedByUser.Role,
		&stock.Provider.ID,
		&stock.Provider.Name,
		&stock.Provider.Email,
		&stock.Provider.Phone,
		&stock.Provider.Address,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (r *MySQLStockRepository) GetAll() ([]domain.Stock, error) {
	query := `
		SELECT 
			s.id, s.serial, s.created_at, s.updated_at,
			s.batch, s.purchase_date,
			p.id, p.name, p.code, p.image_url,
			u1.id, u1.name, u1.email, u1.role,
			u2.id, u2.name, u2.email, u2.role,
			pr.id, pr.name, pr.email, pr.phone, pr.address
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN users u1 ON s.created_by_user_id = u1.id
		JOIN users u2 ON s.updated_by_user_id = u2.id
		JOIN providers pr ON s.provider_id = pr.id
		ORDER BY s.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.Serial,
			&stock.CreatedAt,
			&stock.UpdatedAt,
			&stock.Batch,
			&stock.PurchaseDate,
			&stock.Product.ID,
			&stock.Product.Name,
			&stock.Product.Code,
			&stock.Product.ImageURL,
			&stock.CreatedByUser.ID,
			&stock.CreatedByUser.Name,
			&stock.CreatedByUser.Email,
			&stock.CreatedByUser.Role,
			&stock.UpdatedByUser.ID,
			&stock.UpdatedByUser.Name,
			&stock.UpdatedByUser.Email,
			&stock.UpdatedByUser.Role,
			&stock.Provider.ID,
			&stock.Provider.Name,
			&stock.Provider.Email,
			&stock.Provider.Phone,
			&stock.Provider.Address,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *MySQLStockRepository) GetByProductID(productID int64) ([]domain.Stock, error) {
	query := `
		SELECT 
			s.id, s.serial, s.created_at, s.updated_at,
			s.batch, s.purchase_date,
			p.id, p.name, p.code, p.image_url,
			u1.id, u1.name, u1.email, u1.role,
			u2.id, u2.name, u2.email, u2.role,
			pr.id, pr.name, pr.email, pr.phone, pr.address
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN users u1 ON s.created_by_user_id = u1.id
		JOIN users u2 ON s.updated_by_user_id = u2.id
		JOIN providers pr ON s.provider_id = pr.id
		WHERE s.product_id = ?
		ORDER BY s.id
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.Serial,
			&stock.CreatedAt,
			&stock.UpdatedAt,
			&stock.Batch,
			&stock.PurchaseDate,
			&stock.Product.ID,
			&stock.Product.Name,
			&stock.Product.Code,
			&stock.Product.ImageURL,
			&stock.CreatedByUser.ID,
			&stock.CreatedByUser.Name,
			&stock.CreatedByUser.Email,
			&stock.CreatedByUser.Role,
			&stock.UpdatedByUser.ID,
			&stock.UpdatedByUser.Name,
			&stock.UpdatedByUser.Email,
			&stock.UpdatedByUser.Role,
			&stock.Provider.ID,
			&stock.Provider.Name,
			&stock.Provider.Email,
			&stock.Provider.Phone,
			&stock.Provider.Address,
		)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, nil
}

func (r *MySQLStockRepository) GetBySerial(serial string) (*domain.Stock, error) {
	query := `
		SELECT 
			s.id, s.serial, s.created_at, s.updated_at,
			s.batch, s.purchase_date,
			p.id, p.name, p.code, p.image_url,
			u1.id, u1.name, u1.email, u1.role,
			u2.id, u2.name, u2.email, u2.role,
			pr.id, pr.name, pr.email, pr.phone, pr.address
		FROM stocks s
		JOIN products p ON s.product_id = p.id
		JOIN users u1 ON s.created_by_user_id = u1.id
		JOIN users u2 ON s.updated_by_user_id = u2.id
		JOIN providers pr ON s.provider_id = pr.id
		WHERE s.serial = ?
	`

	var stock domain.Stock
	err := r.db.QueryRow(query, serial).Scan(
		&stock.ID,
		&stock.Serial,
		&stock.CreatedAt,
		&stock.UpdatedAt,
		&stock.Batch,
		&stock.PurchaseDate,
		&stock.Product.ID,
		&stock.Product.Name,
		&stock.Product.Code,
		&stock.Product.ImageURL,
		&stock.CreatedByUser.ID,
		&stock.CreatedByUser.Name,
		&stock.CreatedByUser.Email,
		&stock.CreatedByUser.Role,
		&stock.UpdatedByUser.ID,
		&stock.UpdatedByUser.Name,
		&stock.UpdatedByUser.Email,
		&stock.UpdatedByUser.Role,
		&stock.Provider.ID,
		&stock.Provider.Name,
		&stock.Provider.Email,
		&stock.Provider.Phone,
		&stock.Provider.Address,
	)
	if err != nil {
		if r.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &stock, nil
}

func (r *MySQLStockRepository) Update(stock *domain.Stock) error {
	query := `
		UPDATE stocks
		SET 
			product_id = ?, serial = ?, updated_at = ?,
			updated_by_user_id = ?, batch = ?, purchase_date = ?,
			provider_id = ?
		WHERE id = ?
	`

	now := r.GetCurrentTimestamp()
	result, err := r.db.Exec(query,
		stock.Product.ID,
		stock.Serial,
		now,
		stock.UpdatedByUser.ID,
		stock.Batch,
		stock.PurchaseDate,
		stock.Provider.ID,
		stock.ID,
	)
	if err != nil {
		if r.IsDuplicateEntry(err) {
			return &domain.StockAlreadyExistsError{Serial: stock.Serial}
		}
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.StockNotFoundError{StockID: stock.ID}
	}

	stock.UpdatedAt = now
	return nil
}

func (r *MySQLStockRepository) Delete(id int64) error {
	query := "DELETE FROM stocks WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := r.GetRowsAffected(result)
	if err != nil {
		return err
	}

	if rows == 0 {
		return &domain.StockNotFoundError{StockID: id}
	}

	return nil
}
