package response

import trackeditemlist "github.com/ljj/gugu-api/internal/core/domain/trackeditemlist"

type ListTrackedItem struct {
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
}

func NewListTrackedItems(items []trackeditemlist.Item) []ListTrackedItem {
	result := make([]ListTrackedItem, 0, len(items))
	for _, item := range items {
		result = append(result, NewListTrackedItem(item))
	}
	return result
}

func NewListTrackedItem(item trackeditemlist.Item) ListTrackedItem {
	return ListTrackedItem{
		TrackedItemID:     item.TrackedItem.ID,
		ProductID:         item.Product.ID,
		Market:            string(item.Product.Market),
		ExternalProductID: item.Product.ExternalProductID,
		OriginalURL:       item.TrackedItem.OriginalURL,
		Title:             item.Product.Title,
		MainImageURL:      item.Product.MainImageURL,
		CurrentPrice:      item.Product.CurrentPrice,
		Currency:          item.Product.Currency,
		ProductURL:        item.Product.ProductURL,
	}
}
