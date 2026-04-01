package skupricehistory

import (
	"context"
	"sync"

	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
)

type MemoryRepository struct {
	mu      sync.RWMutex
	bySKUID map[string][]domainsph.SKUPriceHistory
}

func NewRepository() *MemoryRepository {
	return &MemoryRepository{
		bySKUID: make(map[string][]domainsph.SKUPriceHistory),
	}
}

func (r *MemoryRepository) Create(_ context.Context, history domainsph.SKUPriceHistory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bySKUID[history.SKUID] = append(r.bySKUID[history.SKUID], history)
	return nil
}

func (r *MemoryRepository) ListBySKUID(_ context.Context, skuID string, currency string) ([]domainsph.SKUPriceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainsph.SKUPriceHistory
	for _, item := range r.bySKUID[skuID] {
		if item.Currency == currency {
			result = append(result, item)
		}
	}
	return result, nil
}
