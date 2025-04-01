package repository

import (
	"inventario/internal/domain"
	"sync"
)

type MemoryProductRepository struct {
	products map[int64]*domain.Product
	mutex    sync.RWMutex
}

func NewMemoryProductRepository() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[int64]*domain.Product),
	}
}

func (r *MemoryProductRepository) Create(product *domain.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[product.ID]; exists {
		return &domain.ProductAlreadyExistsError{
			Code: product.Code,
		}
	}

	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) GetByID(id int64) (*domain.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, &domain.ProductNotFoundError{
			ProductID: id,
		}
	}

	return product, nil
}

func (r *MemoryProductRepository) GetAll() ([]*domain.Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	products := make([]*domain.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, product)
	}

	return products, nil
}

func (r *MemoryProductRepository) Update(product *domain.Product) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return &domain.ProductNotFoundError{
			ProductID: product.ID,
		}
	}

	r.products[product.ID] = product
	return nil
}

func (r *MemoryProductRepository) Delete(id int64) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[id]; !exists {
		return &domain.ProductNotFoundError{
			ProductID: id,
		}
	}

	delete(r.products, id)
	return nil
}
