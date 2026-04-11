package pricesnapshot

import (
	"context"
	"time"
)

type Service struct {
	skuFinder SKUSnapshotFinder
}

func NewService(skuFinder SKUSnapshotFinder) *Service {
	return &Service{
		skuFinder: skuFinder,
	}
}

func (s *Service) ListSKUSnapshotsByDateRange(ctx context.Context, skuID string, currency string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error) {
	return s.skuFinder.ListBySKUIDAndDateRange(ctx, skuID, currency, from, to)
}
