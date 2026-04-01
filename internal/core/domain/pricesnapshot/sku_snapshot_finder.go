package pricesnapshot

import (
	"context"
	"time"
)

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
