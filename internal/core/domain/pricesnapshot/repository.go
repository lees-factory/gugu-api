package pricesnapshot

import (
	"context"
	"time"
)

type ProductSnapshotRepository interface {
	Upsert(ctx context.Context, snapshot ProductPriceSnapshot) error
	ListByProductIDAndDateRange(ctx context.Context, productID string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error)
}

type SKUSnapshotRepository interface {
	Upsert(ctx context.Context, snapshot SKUPriceSnapshot) error
	ListBySKUIDAndDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error)
}
