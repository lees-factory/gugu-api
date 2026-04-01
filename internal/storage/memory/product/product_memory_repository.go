package product

import (
	"context"
	"sync"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

type ProductMemoryRepository struct {
	mu                   sync.RWMutex
	byID                 map[string]domainproduct.Product
	byMarketExternalID   map[string]string
}

func NewRepository() *ProductMemoryRepository {
	return &ProductMemoryRepository{
		byID:               make(map[string]domainproduct.Product),
		byMarketExternalID: make(map[string]string),
	}
}

func (r *ProductMemoryRepository) FindByID(_ context.Context, productID string) (*domainproduct.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[productID]
	if !ok {
		return nil, nil
	}
	found := item
	return &found, nil
}

func (r *ProductMemoryRepository) FindByIDs(_ context.Context, productIDs []string) ([]domainproduct.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []domainproduct.Product
	for _, id := range productIDs {
		if item, ok := r.byID[id]; ok {
			products = append(products, item)
		}
	}
	return products, nil
}

func (r *ProductMemoryRepository) FindByMarketAndExternalProductID(_ context.Context, market enum.Market, externalProductID string) (*domainproduct.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := string(market) + ":" + externalProductID
	productID, ok := r.byMarketExternalID[key]
	if !ok {
		return nil, nil
	}
	item := r.byID[productID]
	found := item
	return &found, nil
}

func (r *ProductMemoryRepository) ListByMarket(_ context.Context, market enum.Market) ([]domainproduct.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var products []domainproduct.Product
	for _, p := range r.byID {
		if p.Market == market {
			products = append(products, p)
		}
	}
	return products, nil
}

func (r *ProductMemoryRepository) ListByCollectionSource(_ context.Context, source string, limit int, offset int) ([]domainproduct.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var matched []domainproduct.Product
	for _, p := range r.byID {
		if p.CollectionSource == source {
			matched = append(matched, p)
		}
	}

	if offset >= len(matched) {
		return nil, nil
	}
	end := offset + limit
	if end > len(matched) {
		end = len(matched)
	}
	return matched[offset:end], nil
}

func (r *ProductMemoryRepository) Create(_ context.Context, product domainproduct.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[product.ID] = product
	r.byMarketExternalID[string(product.Market)+":"+product.ExternalProductID] = product.ID
	return nil
}

func (r *ProductMemoryRepository) Update(_ context.Context, product domainproduct.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[product.ID] = product
	r.byMarketExternalID[string(product.Market)+":"+product.ExternalProductID] = product.ID
	return nil
}
