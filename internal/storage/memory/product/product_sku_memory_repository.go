package product

import (
	"context"
	"sync"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductSKUMemoryRepository struct {
	mu          sync.RWMutex
	byID        map[string]domainproduct.SKU
	byProductID map[string][]string
}

func NewSKURepository() *ProductSKUMemoryRepository {
	return &ProductSKUMemoryRepository{
		byID:        make(map[string]domainproduct.SKU),
		byProductID: make(map[string][]string),
	}
}

func (r *ProductSKUMemoryRepository) Create(_ context.Context, sku domainproduct.SKU) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byID[sku.ID] = sku
	r.byProductID[sku.ProductID] = append(r.byProductID[sku.ProductID], sku.ID)
	return nil
}

func (r *ProductSKUMemoryRepository) Upsert(_ context.Context, sku domainproduct.SKU) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, existing := range r.byID {
		if existing.ProductID == sku.ProductID && existing.ExternalSKUID == sku.ExternalSKUID {
			sku.ID = id
			r.byID[id] = sku
			return nil
		}
	}
	r.byID[sku.ID] = sku
	r.byProductID[sku.ProductID] = append(r.byProductID[sku.ProductID], sku.ID)
	return nil
}

func (r *ProductSKUMemoryRepository) FindByProductIDAndExternalSKUID(_ context.Context, productID string, externalSKUID string) (*domainproduct.SKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, sku := range r.byID {
		if sku.ProductID == productID && sku.ExternalSKUID == externalSKUID {
			found := sku
			return &found, nil
		}
	}
	return nil, nil
}

func (r *ProductSKUMemoryRepository) FindByID(_ context.Context, skuID string) (*domainproduct.SKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, ok := r.byID[skuID]
	if !ok {
		return nil, nil
	}
	found := item
	return &found, nil
}

func (r *ProductSKUMemoryRepository) FindByProductID(_ context.Context, productID string) ([]domainproduct.SKU, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byProductID[productID]
	skus := make([]domainproduct.SKU, 0, len(ids))
	for _, id := range ids {
		if sku, ok := r.byID[id]; ok {
			skus = append(skus, sku)
		}
	}
	return skus, nil
}
