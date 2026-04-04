package product

import "github.com/ljj/gugu-api/internal/core/enum"

type NewProduct struct {
	Market            enum.Market
	ExternalProductID string
	OriginalURL       string
	Language          string
	Title             string
	MainImageURL      string
	CurrentPrice      string
	Currency          string
	ProductURL        string
	PromotionLink     string
	CollectionSource  string
	SKUs              []NewSKU
}
