package response

import (
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
)

type PriceTrendPoint struct {
	Date          string `json:"date"`
	Price         string `json:"price"`
	OriginalPrice string `json:"original_price,omitempty"`
	Currency      string `json:"currency"`
}

type PriceTrendResponse struct {
	Points []PriceTrendPoint `json:"points"`
}

func NewSKUPriceTrend(snapshots []domainps.SKUPriceSnapshot) PriceTrendResponse {
	points := make([]PriceTrendPoint, 0, len(snapshots))
	for _, s := range snapshots {
		points = append(points, PriceTrendPoint{
			Date:          s.SnapshotDate.Format(time.DateOnly),
			Price:         s.Price,
			OriginalPrice: s.OriginalPrice,
			Currency:      s.Currency,
		})
	}
	return PriceTrendResponse{Points: points}
}
