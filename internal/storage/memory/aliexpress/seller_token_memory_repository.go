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
	byUserID   map[string]string
}

func NewRepository() *SellerTokenMemoryRepository {
	return &SellerTokenMemoryRepository{
		bySellerID: make(map[string]clientaliexpress.SellerTokenRecord),
		byUserID:   make(map[string]string),
	}
}

func (r *SellerTokenMemoryRepository) Upsert(_ context.Context, token clientaliexpress.SellerTokenRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.bySellerID[token.SellerID] = token
	r.byUserID[token.UserID] = token.SellerID
	return nil
}

func (r *SellerTokenMemoryRepository) FindByUserID(_ context.Context, userID string) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	sellerID, ok := r.byUserID[userID]
	if !ok {
		return nil, nil
	}

	token := r.bySellerID[sellerID]
	return &token, nil
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
