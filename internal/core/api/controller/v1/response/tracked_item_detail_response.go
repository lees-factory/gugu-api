package response

import (
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type TrackedItemDetail struct {
	TrackedItemID     string       `json:"tracked_item_id"`
	SKUID             string       `json:"sku_id,omitempty"`
	Market            string       `json:"market"`
	ExternalProductID string       `json:"external_product_id"`
	OriginalURL       string       `json:"original_url"`
	Title             string       `json:"title"`
	MainImageURL      string       `json:"main_image_url"`
	CurrentPrice      string       `json:"current_price"`
	Currency          string       `json:"currency"`
	ProductURL        string       `json:"product_url"`
	SKUs              []ProductSKU `json:"skus"`
}

func NewTrackedItemDetail(detail *domaintrackeditem.TrackedItemDetail) TrackedItemDetail {
	return TrackedItemDetail{
		TrackedItemID:     detail.TrackedItem.ID,
		SKUID:             detail.TrackedItem.SKUID,
		Market:            string(detail.Product.Market),
		ExternalProductID: detail.Product.ExternalProductID,
		OriginalURL:       detail.TrackedItem.OriginalURL,
		Title:             detail.Product.Title,
		MainImageURL:      detail.Product.MainImageURL,
		CurrentPrice:      detail.Product.CurrentPrice,
		Currency:          detail.Product.Currency,
		ProductURL:        detail.Product.ProductURL,
		SKUs:              NewProductSKUs(detail.SKUs),
	}
}
