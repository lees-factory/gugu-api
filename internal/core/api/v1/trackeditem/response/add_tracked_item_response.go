package response

import (
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type AddTrackedItem struct {
	TrackedItemID     string `json:"tracked_item_id"`
	ProductID         string `json:"product_id"`
	Market            string `json:"market"`
	ExternalProductID string `json:"external_product_id"`
	OriginalURL       string `json:"original_url"`
	Title             string `json:"title"`
	MainImageURL      string `json:"main_image_url"`
	CurrentPrice      string `json:"current_price"`
	Currency          string `json:"currency"`
	ProductURL        string `json:"product_url"`
	AlreadyTracked    bool   `json:"already_tracked"`
}

func NewAddTrackedItem(tracked domaintrackeditem.TrackedItem, product domainproduct.Product, alreadyTracked bool) AddTrackedItem {
	return AddTrackedItem{
		TrackedItemID:     tracked.ID,
		ProductID:         product.ID,
		Market:            string(product.Market),
		ExternalProductID: product.ExternalProductID,
		OriginalURL:       tracked.OriginalURL,
		Title:             product.Title,
		MainImageURL:      product.MainImageURL,
		CurrentPrice:      product.CurrentPrice,
		Currency:          product.Currency,
		ProductURL:        product.ProductURL,
		AlreadyTracked:    alreadyTracked,
	}
}
