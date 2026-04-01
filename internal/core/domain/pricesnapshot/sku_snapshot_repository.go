package pricesnapshot

import (
	"context"
	"time"
)

type SKUSnapshotRepository interface {
	Upsert(ctx context.Context, snapshot SKUPriceSnapshot) error
	ListBySKUIDAndDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]SKUPriceSnapshot, error)
}
