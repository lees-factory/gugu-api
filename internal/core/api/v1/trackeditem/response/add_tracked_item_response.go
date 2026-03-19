package response

import (
	productresponse "github.com/ljj/gugu-api/internal/core/api/v1/product/response"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type AddTrackedItem struct {
	TrackedItemID     string                      `json:"tracked_item_id"`
	ProductID         string                      `json:"product_id"`
	Market            string                      `json:"market"`
	ExternalProductID string                      `json:"external_product_id"`
	OriginalURL       string                      `json:"original_url"`
	Title             string                      `json:"title"`
	MainImageURL      string                      `json:"main_image_url"`
	CurrentPrice      string                      `json:"current_price"`
	Currency          string                      `json:"currency"`
	ProductURL        string                      `json:"product_url"`
	AlreadyTracked    bool                        `json:"already_tracked"`
	SKUs              []productresponse.ProductSKU `json:"skus"`
}

func NewAddTrackedItemFromResult(result *domaintrackeditem.AddTrackedItemResult) AddTrackedItem {
	return AddTrackedItem{
		TrackedItemID:     result.TrackedItem.ID,
		ProductID:         result.Product.ID,
		Market:            string(result.Product.Market),
		ExternalProductID: result.Product.ExternalProductID,
		OriginalURL:       result.TrackedItem.OriginalURL,
		Title:             result.Product.Title,
		MainImageURL:      result.Product.MainImageURL,
		CurrentPrice:      result.Product.CurrentPrice,
		Currency:          result.Product.Currency,
		ProductURL:        result.Product.ProductURL,
		AlreadyTracked:    result.AlreadyTracked,
		SKUs:              productresponse.NewProductSKUs(result.SKUs),
	}
}
