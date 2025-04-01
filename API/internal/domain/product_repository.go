package domain

// IProductRepository defines the interface for product persistence operations
type IProductRepository interface {
	Create(product *Product) error
	GetAll() ([]Product, error)
	GetByID(id int64) (*Product, error)
	Update(product *Product) error
	Delete(id int64) error
}
