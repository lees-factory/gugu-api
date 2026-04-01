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

func (s *Service) ListProductSnapshotsByDateRange(ctx context.Context, productID string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error) {
	return s.productFinder.ListByProductIDAndDateRange(ctx, productID, from, to)
}

func (s *Service) ListSKUSnapshotsByDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error) {
	return s.skuFinder.ListBySKUIDAndDateRange(ctx, skuID, from, to)
}
