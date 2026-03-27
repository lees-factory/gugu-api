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

type SKUSnapshotFinder interface {
	ListBySKUIDAndDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error)
}

type skuSnapshotFinder struct {
	repository SKUSnapshotRepository
}

func NewSKUSnapshotFinder(repository SKUSnapshotRepository) SKUSnapshotFinder {
	return &skuSnapshotFinder{repository: repository}
}

func (f *skuSnapshotFinder) ListBySKUIDAndDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error) {
	return f.repository.ListBySKUIDAndDateRange(ctx, skuID, from, to)
}
