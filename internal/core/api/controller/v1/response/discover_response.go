package response

import (
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type HotProductItem struct {
	ProductID    string `json:"product_id"`
	Market       string `json:"market"`
	Title        string `json:"title"`
	MainImageURL string `json:"main_image_url"`
	CurrentPrice string `json:"current_price"`
	Currency     string `json:"currency"`
	ProductURL   string `json:"product_url"`
}

func NewHotProductList(products []domainproduct.Product) []HotProductItem {
	items := make([]HotProductItem, 0, len(products))
	for _, p := range products {
		items = append(items, HotProductItem{
			ProductID:    p.ID,
			Market:       string(p.Market),
			Title:        p.Title,
			MainImageURL: p.MainImageURL,
			CurrentPrice: p.CurrentPrice,
			Currency:     p.Currency,
			ProductURL:   p.ProductURL,
		})
	}
	return items
}
