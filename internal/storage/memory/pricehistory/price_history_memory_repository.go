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

func (r *PriceHistoryMemoryRepository) Create(_ context.Context, history domainpricehistory.PriceHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byProductID[history.ProductID] = append(r.byProductID[history.ProductID], history)
	return nil
}

func (r *PriceHistoryMemoryRepository) ListByProductID(_ context.Context, productID string, currency string) ([]domainpricehistory.PriceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainpricehistory.PriceHistory
	for _, item := range r.byProductID[productID] {
		if item.Currency == currency {
			result = append(result, item)
		}
	}
	return result, nil
}
