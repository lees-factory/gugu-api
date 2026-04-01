package pricealert

import "context"

type Finder interface {
	FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*PriceAlert, error)
	ListByProductID(ctx context.Context, productID string) ([]PriceAlert, error)
	ListByProductIDs(ctx context.Context, productIDs []string) ([]PriceAlert, error)
	ListByUserID(ctx context.Context, userID string) ([]PriceAlert, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*PriceAlert, error) {
	return f.repository.FindByUserIDAndProductID(ctx, userID, productID)
}

func (f *finder) ListByProductID(ctx context.Context, productID string) ([]PriceAlert, error) {
	return f.repository.ListByProductID(ctx, productID)
}

func (f *finder) ListByProductIDs(ctx context.Context, productIDs []string) ([]PriceAlert, error) {
	return f.repository.ListByProductIDs(ctx, productIDs)
}

func (f *finder) ListByUserID(ctx context.Context, userID string) ([]PriceAlert, error) {
	return f.repository.ListByUserID(ctx, userID)
}
