package pricesnapshot

import (
	"context"
	"time"
)

type Service struct {
	productFinder ProductSnapshotFinder
	skuFinder     SKUSnapshotFinder
}

func NewService(productFinder ProductSnapshotFinder, skuFinder SKUSnapshotFinder) *Service {
	return &Service{
		productFinder: productFinder,
		skuFinder:     skuFinder,
	}
}

func (s *Service) ListProductSnapshotsByDateRange(ctx context.Context, productID string, currency string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error) {
	return s.productFinder.ListByProductIDAndDateRange(ctx, productID, currency, from, to)
}

func (s *Service) ListSKUSnapshotsByDateRange(ctx context.Context, skuID string, currency string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error) {
	return s.skuFinder.ListBySKUIDAndDateRange(ctx, skuID, currency, from, to)
}
