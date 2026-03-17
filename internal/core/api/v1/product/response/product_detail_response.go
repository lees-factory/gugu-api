package response

import (
	"time"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductDetail struct {
	ProductID         string             `json:"product_id"`
	Market            string             `json:"market"`
	ExternalProductID string             `json:"external_product_id"`
	OriginalURL       string             `json:"original_url"`
	Title             string             `json:"title"`
	MainImageURL      string             `json:"main_image_url"`
	CurrentPrice      string             `json:"current_price"`
	Currency          string             `json:"currency"`
	ProductURL        string             `json:"product_url"`
	IsTrackedByUser   bool               `json:"is_tracked_by_user"`
	TrackedItemID     string             `json:"tracked_item_id,omitempty"`
	PriceHistories    []PriceHistoryItem `json:"price_histories"`
}

type PriceHistoryItem struct {
	Price       string    `json:"price"`
	Currency    string    `json:"currency"`
	RecordedAt  time.Time `json:"recorded_at"`
	ChangeValue string    `json:"change_value"`
}

func NewProductDetail(product domainproduct.Product, histories []domainpricehistory.PriceHistory, isTrackedByUser bool, trackedItemID string) ProductDetail {
	items := make([]PriceHistoryItem, 0, len(histories))
	for _, h := range histories {
		items = append(items, PriceHistoryItem{
			Price:       h.Price,
			Currency:    h.Currency,
			RecordedAt:  h.RecordedAt,
			ChangeValue: h.ChangeValue,
		})
	}

	return ProductDetail{
		ProductID:         product.ID,
		Market:            string(product.Market),
		ExternalProductID: product.ExternalProductID,
		OriginalURL:       product.OriginalURL,
		Title:             product.Title,
		MainImageURL:      product.MainImageURL,
		CurrentPrice:      product.CurrentPrice,
		Currency:          product.Currency,
		ProductURL:        product.ProductURL,
		IsTrackedByUser:   isTrackedByUser,
		TrackedItemID:     trackedItemID,
		PriceHistories:    items,
	}
}
