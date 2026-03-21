package product

import "github.com/ljj/gugu-api/internal/core/enum"

type NewProduct struct {
	Market            enum.Market
	ExternalProductID string
	OriginalURL       string
	Title             string
	MainImageURL      string
	CurrentPrice      string
	Currency          string
	ProductURL        string
	CollectionSource  string
	SKUs              []NewSKU
}
