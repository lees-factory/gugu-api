package product

import "context"

type Finder interface {
	FindByID(ctx context.Context, productID string) (*Product, error)
	FindByIDs(ctx context.Context, productIDs []string) ([]Product, error)
	FindByMarketAndExternalProductID(ctx context.Context, market Market, externalProductID string) (*Product, error)
}

type finder struct {
	repository Repository
}

func NewFinder(repository Repository) Finder {
	return &finder{repository: repository}
}

func (f *finder) FindByID(ctx context.Context, productID string) (*Product, error) {
	return f.repository.FindByID(ctx, productID)
}

func (f *finder) FindByIDs(ctx context.Context, productIDs []string) ([]Product, error) {
	return f.repository.FindByIDs(ctx, productIDs)
}

func (f *finder) FindByMarketAndExternalProductID(ctx context.Context, market Market, externalProductID string) (*Product, error) {
	return f.repository.FindByMarketAndExternalProductID(ctx, market, externalProductID)
}
