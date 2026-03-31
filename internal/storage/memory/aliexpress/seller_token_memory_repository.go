package aliexpress

import (
	"context"
	"sync"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type SellerTokenMemoryRepository struct {
	mu     sync.RWMutex
	tokens map[string]clientaliexpress.SellerTokenRecord // key: appType + ":" + sellerID
}

func NewRepository() *SellerTokenMemoryRepository {
	return &SellerTokenMemoryRepository{
		tokens: make(map[string]clientaliexpress.SellerTokenRecord),
	}
}

func tokenKey(token clientaliexpress.SellerTokenRecord) string {
	return token.AppType + ":" + token.SellerID
}

func (r *SellerTokenMemoryRepository) Upsert(_ context.Context, token clientaliexpress.SellerTokenRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[tokenKey(token)] = token
	return nil
}

func (r *SellerTokenMemoryRepository) FindOne(_ context.Context) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.tokens {
		return &token, nil
	}
	return nil, nil
}

func (r *SellerTokenMemoryRepository) FindByAppType(_ context.Context, appType string) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.tokens {
		if token.AppType == appType {
			return &token, nil
		}
	}
	return nil, nil
}

func (r *SellerTokenMemoryRepository) FindBySellerID(_ context.Context, sellerID string) (*clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, token := range r.tokens {
		if token.SellerID == sellerID {
			return &token, nil
		}
	}
	return nil, nil
}

func (r *SellerTokenMemoryRepository) ListExpiringBefore(_ context.Context, expiresBefore time.Time) ([]clientaliexpress.SellerTokenRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]clientaliexpress.SellerTokenRecord, 0)
	for _, token := range r.tokens {
		if token.AccessTokenExpiresAt.After(expiresBefore) {
			continue
		}
		items = append(items, token)
	}

	return items, nil
}
