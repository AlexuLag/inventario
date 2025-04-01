package domain

import "time"

// Product represents the core product entity in the domain
type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ImageURL  string    `json:"image_url"`
}

// NewProduct creates a new Product instance with default values
func NewProduct(name, code, imageURL string) *Product {
	now := time.Now()
	return &Product{
		Name:      name,
		Code:      code,
		ImageURL:  imageURL,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
