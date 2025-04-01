package usecase

import (
	"inventario/internal/domain"
	"time"
)

type StockUseCase struct {
	stockRepo domain.IStockRepository
}

func NewStockUseCase(stockRepo domain.IStockRepository) *StockUseCase {
	return &StockUseCase{
		stockRepo: stockRepo,
	}
}

func (uc *StockUseCase) CreateStock(productID int64, serial string, batch string, purchaseDate time.Time, providerID int64, createdByUserID int64) (*domain.Stock, error) {
	stock := &domain.Stock{
		Product: &domain.Product{
			ID: productID,
		},
		Serial:       serial,
		Batch:        batch,
		PurchaseDate: purchaseDate,
		Provider: &domain.Provider{
			ID: providerID,
		},
		CreatedByUser: &domain.User{
			ID: createdByUserID,
		},
		UpdatedByUser: &domain.User{
			ID: createdByUserID,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.stockRepo.Create(stock); err != nil {
		return nil, err
	}

	return stock, nil
}

func (uc *StockUseCase) GetStock(id int64) (*domain.Stock, error) {
	stock, err := uc.stockRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, &domain.StockNotFoundError{StockID: id}
	}
	return stock, nil
}

func (uc *StockUseCase) GetAllStocks() ([]*domain.Stock, error) {
	stocks, err := uc.stockRepo.GetAll()
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Stock, len(stocks))
	for i, stock := range stocks {
		result[i] = &stock
	}
	return result, nil
}

func (uc *StockUseCase) UpdateStock(stock *domain.Stock) error {
	existingStock, err := uc.stockRepo.GetByID(stock.ID)
	if err != nil {
		return err
	}
	if existingStock == nil {
		return &domain.StockNotFoundError{StockID: stock.ID}
	}

	stock.UpdatedAt = time.Now()
	return uc.stockRepo.Update(stock)
}

func (uc *StockUseCase) DeleteStock(id int64) error {
	stock, err := uc.stockRepo.GetByID(id)
	if err != nil {
		return err
	}
	if stock == nil {
		return &domain.StockNotFoundError{StockID: id}
	}
	return uc.stockRepo.Delete(id)
}

func (uc *StockUseCase) GetStocksByProductID(productID int64) ([]*domain.Stock, error) {
	stocks, err := uc.stockRepo.GetByProductID(productID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Stock, len(stocks))
	for i, stock := range stocks {
		result[i] = &stock
	}
	return result, nil
}

func (uc *StockUseCase) GetStockBySerial(serial string) (*domain.Stock, error) {
	stock, err := uc.stockRepo.GetBySerial(serial)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, &domain.StockNotFoundError{Serial: serial}
	}
	return stock, nil
}
