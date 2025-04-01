package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInternalServer = errors.New("internal server error")
)

// ProductAlreadyExistsError represents an error when a product already exists
type ProductAlreadyExistsError struct {
	Code string
}

func (e *ProductAlreadyExistsError) Error() string {
	return fmt.Sprintf("product with code %s already exists", e.Code)
}

// ProductNotFoundError represents an error when a product is not found
type ProductNotFoundError struct {
	ProductID int64
}

func (e *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product with ID %d not found", e.ProductID)
}
