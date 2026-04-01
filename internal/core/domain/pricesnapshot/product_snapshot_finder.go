package pricesnapshot

import (
	"context"
	"time"
)

type ProductSnapshotFinder interface {
	ListByProductIDAndDateRange(ctx context.Context, productID string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error)
}

type productSnapshotFinder struct {
	repository ProductSnapshotRepository
}

func NewProductSnapshotFinder(repository ProductSnapshotRepository) ProductSnapshotFinder {
	return &productSnapshotFinder{repository: repository}
}

func (f *productSnapshotFinder) ListByProductIDAndDateRange(ctx context.Context, productID string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error) {
	return f.repository.ListByProductIDAndDateRange(ctx, productID, from, to)
}
