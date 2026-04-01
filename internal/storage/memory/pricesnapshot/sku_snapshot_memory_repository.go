package pricesnapshot

import (
	"context"
	"sync"
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
)

type SKUSnapshotMemoryRepository struct {
	mu      sync.RWMutex
	bySKUID map[string][]domainps.SKUPriceSnapshot
}

func NewSKUSnapshotRepository() *SKUSnapshotMemoryRepository {
	return &SKUSnapshotMemoryRepository{
		bySKUID: make(map[string][]domainps.SKUPriceSnapshot),
	}
}

func (r *SKUSnapshotMemoryRepository) Upsert(_ context.Context, snapshot domainps.SKUPriceSnapshot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	items := r.bySKUID[snapshot.SKUID]
	for i, item := range items {
		if item.SnapshotDate.Equal(snapshot.SnapshotDate) {
			items[i] = snapshot
			return nil
		}
	}
	r.bySKUID[snapshot.SKUID] = append(items, snapshot)
	return nil
}

func (r *SKUSnapshotMemoryRepository) ListBySKUIDAndDateRange(_ context.Context, skuID string, currency string, from time.Time, to time.Time) ([]domainps.SKUPriceSnapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainps.SKUPriceSnapshot
	for _, item := range r.bySKUID[skuID] {
		if item.Currency == currency && !item.SnapshotDate.Before(from) && !item.SnapshotDate.After(to) {
			result = append(result, item)
		}
	}
	return result, nil
}
