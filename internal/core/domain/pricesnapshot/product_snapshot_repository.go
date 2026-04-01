package pricesnapshot

import (
	"context"
	"time"
)

type ProductSnapshotRepository interface {
	Upsert(ctx context.Context, snapshot ProductPriceSnapshot) error
	ListByProductIDAndDateRange(ctx context.Context, productID string, currency string, from time.Time, to time.Time) ([]ProductPriceSnapshot, error)
}
