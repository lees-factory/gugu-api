package response

import (
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
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
	CurrentPrice      string       `json:"current_price"`
	Currency          string       `json:"currency"`
	ProductURL        string       `json:"product_url"`
	SKUs              []ProductSKU `json:"skus"`
}

type PriceAlertState struct {
	Enabled bool   `json:"enabled"`
	Channel string `json:"channel,omitempty"`
}

func NewPriceAlertState(alert *domainpricealert.PriceAlert) PriceAlertState {
	if alert == nil {
		return PriceAlertState{Enabled: false}
	}
	return PriceAlertState{
		Enabled: alert.Enabled,
		Channel: alert.Channel,
	}
}

func NewTrackedItemDetail(detail *domaintrackeditem.TrackedItemDetail, currentPrice string, skus []ProductSKU) TrackedItemDetail {
	display := resolveTrackedItemDisplay(detail.Variant)

	return TrackedItemDetail{
		TrackedItemID:     detail.TrackedItem.ID,
		ProductID:         detail.Product.ID,
		SKUID:             detail.TrackedItem.SKUID,
		Market:            string(detail.Product.Market),
		ExternalProductID: detail.Product.ExternalProductID,
		OriginalURL:       detail.TrackedItem.OriginalURL,
		Title:             display.title,
		MainImageURL:      display.mainImageURL,
		CurrentPrice:      currentPrice,
		Currency:          detail.TrackedItem.Currency,
		ProductURL:        display.productURL,
		SKUs:              skus,
	}
}
