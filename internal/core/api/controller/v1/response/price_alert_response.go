package response

import (
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
)

type PriceAlertItem struct {
	ID        string `json:"id"`
	SKUID     string `json:"sku_id"`
	Channel   string `json:"channel"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"created_at"`
}

func NewPriceAlertItem(alert *domainpricealert.PriceAlert) PriceAlertItem {
	return PriceAlertItem{
		ID:        alert.ID,
		SKUID:     alert.SKUID,
		Channel:   alert.Channel,
		Enabled:   alert.Enabled,
		CreatedAt: alert.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func NewPriceAlertList(alerts []domainpricealert.PriceAlert) []PriceAlertItem {
	items := make([]PriceAlertItem, 0, len(alerts))
	for _, a := range alerts {
		items = append(items, PriceAlertItem{
			ID:        a.ID,
			SKUID:     a.SKUID,
			Channel:   a.Channel,
			Enabled:   a.Enabled,
			CreatedAt: a.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}
	return items
}
