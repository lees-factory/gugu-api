package response

import domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"

type HotProductItem struct {
	ProductID    string `json:"product_id"`
	Market       string `json:"market"`
	Title        string `json:"title"`
	MainImageURL string `json:"main_image_url"`
	Currency     string `json:"currency"`
	ProductURL   string `json:"product_url"`
}

func NewHotProductList(products []domainproduct.Product, variants map[string]*domainproduct.Variant, fallbackCurrency string) []HotProductItem {
	items := make([]HotProductItem, 0, len(products))
	for _, p := range products {
		title := p.Title
		mainImageURL := p.MainImageURL
		currency := fallbackCurrency
		productURL := p.ProductURL

		if variant := variants[p.ID]; variant != nil {
			if variant.Title != "" {
				title = variant.Title
			}
			if variant.MainImageURL != "" {
				mainImageURL = variant.MainImageURL
			}
			if variant.Currency != "" {
				currency = variant.Currency
			}
			if variant.ProductURL != "" {
				productURL = variant.ProductURL
			}
		}

		items = append(items, HotProductItem{
			ProductID:    p.ID,
			Market:       string(p.Market),
			Title:        title,
			MainImageURL: mainImageURL,
			Currency:     currency,
			ProductURL:   productURL,
		})
	}
	return items
}
