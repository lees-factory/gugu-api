package response

import (
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type TrackedItemDetail struct {
	TrackedItemID     string       `json:"tracked_item_id"`
	ProductID         string       `json:"product_id"`
	SKUID             string       `json:"sku_id,omitempty"`
	Market            string       `json:"market"`
	ExternalProductID string       `json:"external_product_id"`
	OriginalURL       string       `json:"original_url"`
	Title             string       `json:"title"`
	MainImageURL      string       `json:"main_image_url"`
	Currency          string       `json:"currency"`
	ProductURL        string       `json:"product_url"`
	SKUs              []ProductSKU `json:"skus"`
}

func NewTrackedItemDetail(detail *domaintrackeditem.TrackedItemDetail) TrackedItemDetail {
	display := resolveTrackedItemDisplay(detail.Product, detail.Variant)

	return TrackedItemDetail{
		TrackedItemID:     detail.TrackedItem.ID,
		ProductID:         detail.Product.ID,
		SKUID:             detail.TrackedItem.SKUID,
		Market:            string(detail.Product.Market),
		ExternalProductID: detail.Product.ExternalProductID,
		OriginalURL:       detail.TrackedItem.OriginalURL,
		Title:             display.title,
		MainImageURL:      display.mainImageURL,
		Currency:          detail.TrackedItem.Currency,
		ProductURL:        display.productURL,
		SKUs:              NewProductSKUs(detail.SKUs),
	}
}
