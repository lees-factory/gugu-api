package pricesnapshot

import (
	"context"
	"sync"
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
)

type ProductSnapshotMemoryRepository struct {
	mu          sync.RWMutex
	byProductID map[string][]domainps.ProductPriceSnapshot
}

func NewProductSnapshotRepository() *ProductSnapshotMemoryRepository {
	return &ProductSnapshotMemoryRepository{
		byProductID: make(map[string][]domainps.ProductPriceSnapshot),
	}
}

func (r *ProductSnapshotMemoryRepository) Upsert(_ context.Context, snapshot domainps.ProductPriceSnapshot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := r.byProductID[snapshot.ProductID]
	for i, item := range items {
		if item.SnapshotDate.Equal(snapshot.SnapshotDate) {
			items[i] = snapshot
			return nil
		}
	}
	r.byProductID[snapshot.ProductID] = append(items, snapshot)
	return nil
}

func (r *ProductSnapshotMemoryRepository) ListByProductIDAndDateRange(_ context.Context, productID string, currency string, from time.Time, to time.Time) ([]domainps.ProductPriceSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainps.ProductPriceSnapshot
	for _, item := range r.byProductID[productID] {
		if item.Currency == currency && !item.SnapshotDate.Before(from) && !item.SnapshotDate.After(to) {
			result = append(result, item)
		}
	}
	return result, nil
}
