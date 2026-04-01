package response

import (
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	"github.com/ljj/gugu-api/internal/core/support/page"
)

type ListTrackedItem struct {
	TrackedItemID     string `json:"tracked_item_id"`
	ProductID         string `json:"product_id"`
	SKUID             string `json:"sku_id,omitempty"`
	Market            string `json:"market"`
	ExternalProductID string `json:"external_product_id"`
	OriginalURL       string `json:"original_url"`
	Title             string `json:"title"`
	MainImageURL      string `json:"main_image_url"`
	CurrentPrice      string `json:"current_price"`
	Currency          string `json:"currency"`
	ProductURL        string `json:"product_url"`
}

func NewListTrackedItem(tracked domaintrackeditem.TrackedItem, product domainproduct.Product) ListTrackedItem {
	return ListTrackedItem{
		TrackedItemID:     tracked.ID,
		ProductID:         product.ID,
		SKUID:             tracked.SKUID,
		Market:            string(product.Market),
		ExternalProductID: product.ExternalProductID,
		OriginalURL:       tracked.OriginalURL,
		Title:             product.Title,
		MainImageURL:      product.MainImageURL,
		CurrentPrice:      product.CurrentPrice,
		Currency:          tracked.Currency,
		ProductURL:        product.ProductURL,
	}
}

func NewListTrackedItems(items []domaintrackeditem.TrackedItemWithProduct) []ListTrackedItem {
	result := make([]ListTrackedItem, 0, len(items))
	for _, item := range items {
		result = append(result, NewListTrackedItem(item.TrackedItem, item.Product))
	}
	return result
}

type ListTrackedItemsPage struct {
	Items      []ListTrackedItem `json:"items"`
	NextCursor string            `json:"next_cursor,omitempty"`
	HasMore    bool              `json:"has_more"`
}

func NewListTrackedItemsPage(p *page.CursorPage[domaintrackeditem.TrackedItemWithProduct]) ListTrackedItemsPage {
	items := make([]ListTrackedItem, 0, len(p.Items))
	for _, item := range p.Items {
		items = append(items, NewListTrackedItem(item.TrackedItem, item.Product))
	}
	return ListTrackedItemsPage{
		Items:      items,
		NextCursor: p.NextCursor,
		HasMore:    p.HasMore,
	}
}
