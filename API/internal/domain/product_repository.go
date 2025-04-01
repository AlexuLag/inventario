package domain

// ProductRepository defines the interface for product persistence operations
type ProductRepository interface {
	Create(product *Product) error
	GetByID(id int64) (*Product, error)
	GetAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id int64) error
}
