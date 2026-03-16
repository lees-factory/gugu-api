package product

import "time"

type Product struct {
	ID                string
	Market            Market
	ExternalProductID string
	OriginalURL       string
	Title             string
	MainImageURL      string
	CurrentPrice      string
	Currency          string
	ProductURL        string
	CollectionSource  string
	LastCollectedAt   time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
