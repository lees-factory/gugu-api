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

func (r *MemoryRepository) ListBySKUID(_ context.Context, skuID string) ([]domainsph.SKUPriceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := r.bySKUID[skuID]
	copied := make([]domainsph.SKUPriceHistory, len(items))
	copy(copied, items)
	return copied, nil
}
