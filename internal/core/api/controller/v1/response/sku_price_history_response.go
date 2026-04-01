package response

import (
	"time"

	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
)

type SKUPriceHistoryItem struct {
	SKUID       string    `json:"sku_id"`
	Price       string    `json:"price"`
	Currency    string    `json:"currency"`
	RecordedAt  time.Time `json:"recorded_at"`
	ChangeValue string    `json:"change_value"`
}

func NewSKUPriceHistories(histories []domainsph.SKUPriceHistory) []SKUPriceHistoryItem {
	items := make([]SKUPriceHistoryItem, 0, len(histories))
	for _, h := range histories {
		items = append(items, SKUPriceHistoryItem{
			SKUID:       h.SKUID,
			Price:       h.Price,
			Currency:    h.Currency,
			RecordedAt:  h.RecordedAt,
			ChangeValue: h.ChangeValue,
		})
	}
	return items
}
