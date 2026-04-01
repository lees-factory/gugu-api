package pricealert

import (
	"context"
	"sync"

	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
)

type MemoryRepository struct {
	mu   sync.RWMutex
	byID map[string]domainpricealert.PriceAlert
}

func NewRepository() *MemoryRepository {
	return &MemoryRepository{
		byID: make(map[string]domainpricealert.PriceAlert),
	}
}

func (r *MemoryRepository) FindByUserIDAndProductID(_ context.Context, userID string, productID string) (*domainpricealert.PriceAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, a := range r.byID {
		if a.UserID == userID && a.ProductID == productID {
			found := a
			return &found, nil
		}
	}
	return nil, nil
}

func (r *MemoryRepository) ListByProductID(_ context.Context, productID string) ([]domainpricealert.PriceAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainpricealert.PriceAlert
	for _, a := range r.byID {
		if a.ProductID == productID && a.Enabled {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *MemoryRepository) ListByProductIDs(_ context.Context, productIDs []string) ([]domainpricealert.PriceAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	idSet := make(map[string]struct{}, len(productIDs))
	for _, id := range productIDs {
		idSet[id] = struct{}{}
	}

	var result []domainpricealert.PriceAlert
	for _, a := range r.byID {
		if _, ok := idSet[a.ProductID]; ok && a.Enabled {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *MemoryRepository) ListByUserID(_ context.Context, userID string) ([]domainpricealert.PriceAlert, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []domainpricealert.PriceAlert
	for _, a := range r.byID {
		if a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (r *MemoryRepository) Create(_ context.Context, alert domainpricealert.PriceAlert) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[alert.ID] = alert
	return nil
}

func (r *MemoryRepository) UpdateEnabled(_ context.Context, alertID string, enabled bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if a, ok := r.byID[alertID]; ok {
		a.Enabled = enabled
		r.byID[alertID] = a
	}
	return nil
}

func (r *MemoryRepository) DeleteByUserIDAndProductID(_ context.Context, userID string, productID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, a := range r.byID {
		if a.UserID == userID && a.ProductID == productID {
			delete(r.byID, id)
		}
	}
	return nil
}
