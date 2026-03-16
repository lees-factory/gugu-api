package pricehistory

import (
	"context"
	"sync"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
)

type PriceHistoryMemoryRepository struct {
	mu            sync.RWMutex
	byProductID   map[string][]domainpricehistory.PriceHistory
}

func NewRepository() *PriceHistoryMemoryRepository {
	return &PriceHistoryMemoryRepository{
		byProductID: make(map[string][]domainpricehistory.PriceHistory),
	}
}

func (r *PriceHistoryMemoryRepository) ListByProductID(_ context.Context, productID string) ([]domainpricehistory.PriceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.byProductID[productID]
	copied := make([]domainpricehistory.PriceHistory, len(items))
	copy(copied, items)
	return copied, nil
}
