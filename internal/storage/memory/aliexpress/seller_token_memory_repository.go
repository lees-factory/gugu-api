package aliexpress

import (
	"context"
	"sync"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type SellerTokenMemoryRepository struct {
	mu         sync.RWMutex
	bySellerID map[string]clientaliexpress.SellerTokenRecord
}

func NewRepository() *SellerTokenMemoryRepository {
	return &SellerTokenMemoryRepository{
		bySellerID: make(map[string]clientaliexpress.SellerTokenRecord),
	}
}

func (r *SellerTokenMemoryRepository) Upsert(_ context.Context, token clientaliexpress.SellerTokenRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bySellerID[token.SellerID] = token
	return nil
}

func (r *SellerTokenMemoryRepository) FindOne(_ context.Context) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.bySellerID {
		return &token, nil
	}
	return nil, nil
}

func (r *SellerTokenMemoryRepository) FindBySellerID(_ context.Context, sellerID string) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, ok := r.bySellerID[sellerID]
	if !ok {
		return nil, nil
	}

	return &token, nil
}

func (r *SellerTokenMemoryRepository) ListExpiringBefore(_ context.Context, expiresBefore time.Time) ([]clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]clientaliexpress.SellerTokenRecord, 0)
	for _, token := range r.bySellerID {
		if token.AccessTokenExpiresAt.After(expiresBefore) {
			continue
		}
		items = append(items, token)
	}

	return items, nil
}
