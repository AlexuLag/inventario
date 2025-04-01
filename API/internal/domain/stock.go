package domain

import (
	"strconv"
	"time"
)

type Stock struct {
	ID            int64     `json:"id"`
	Product       *Product  `json:"product"`
	Serial        string    `json:"serial"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedByUser *User     `json:"created_by_user"`
	UpdatedByUser *User     `json:"updated_by_user"`
	Batch         string    `json:"batch"`
	PurchaseDate  time.Time `json:"purchase_date"`
	Provider      *Provider `json:"provider"`
}

type StockNotFoundError struct {
	StockID int64
	Serial  string
}

func (e *StockNotFoundError) Error() string {
	if e.Serial != "" {
		return "stock not found with serial: " + e.Serial
	}
	return "stock not found with id: " + strconv.FormatInt(e.StockID, 10)
}

type StockAlreadyExistsError struct {
	Serial string
}

func (e *StockAlreadyExistsError) Error() string {
	return "stock already exists"
}

type IStockRepository interface {
	Create(stock *Stock) error
	GetByID(id int64) (*Stock, error)
	GetAll() ([]Stock, error)
	Update(stock *Stock) error
	Delete(id int64) error
	GetByProductID(productID int64) ([]Stock, error)
	GetBySerial(serial string) (*Stock, error)
}
