package usecase

import (
	"inventario/internal/domain"
)

// ProductUseCase handles the business logic for product operations
type ProductUseCase struct {
	productRepo domain.IProductRepository
}

// NewProductUseCase creates a new ProductUseCase instance
func NewProductUseCase(repo domain.IProductRepository) *ProductUseCase {
	return &ProductUseCase{
		productRepo: repo,
	}
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(name, code, imageURL string) (*domain.Product, error) {
	product := domain.NewProduct(name, code, imageURL)
	err := uc.productRepo.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

// GetProduct retrieves a product by ID
func (uc *ProductUseCase) GetProduct(id int64) (*domain.Product, error) {
	return uc.productRepo.GetByID(id)
}

// GetAllProducts retrieves all products
func (uc *ProductUseCase) GetAllProducts() ([]*domain.Product, error) {
	products, err := uc.productRepo.GetAll()
	if err != nil {
		return nil, err
	}

	result := make([]*domain.Product, len(products))
	for i := range products {
		result[i] = &products[i]
	}
	return result, nil
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(product *domain.Product) error {
	return uc.productRepo.Update(product)
}

// DeleteProduct deletes a product by ID
func (uc *ProductUseCase) DeleteProduct(id int64) error {
	return uc.productRepo.Delete(id)
}
