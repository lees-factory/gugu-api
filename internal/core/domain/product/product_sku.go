package product

import "time"

type ProductSKU struct {
	ID            string
	ProductID     string
	ExternalSKUID string
	SKUName       string
	Color         string
	Size          string
	Price         string
	OriginalPrice string
	Currency      string
	ImageURL      string
	SKUProperties string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
