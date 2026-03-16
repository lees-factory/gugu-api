package response

import (
	"time"

	productdetail "github.com/ljj/gugu-api/internal/core/domain/productdetail"
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

func NewProductDetail(source productdetail.Detail) ProductDetail {
	items := make([]PriceHistoryItem, 0, len(source.PriceHistories))
	for _, item := range source.PriceHistories {
		items = append(items, PriceHistoryItem{
			Price:       item.Price,
			Currency:    item.Currency,
			RecordedAt:  item.RecordedAt,
			ChangeValue: item.ChangeValue,
		})
	}

	return ProductDetail{
		ProductID:         source.Product.ID,
		Market:            string(source.Product.Market),
		ExternalProductID: source.Product.ExternalProductID,
		OriginalURL:       source.Product.OriginalURL,
		Title:             source.Product.Title,
		MainImageURL:      source.Product.MainImageURL,
		CurrentPrice:      source.Product.CurrentPrice,
		Currency:          source.Product.Currency,
		ProductURL:        source.Product.ProductURL,
		IsTrackedByUser:   source.IsTrackedByUser,
		TrackedItemID:     source.TrackedItemID,
		PriceHistories:    items,
	}
}
