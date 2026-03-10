package aliexpress

import "context"

type ProductLookupInput struct {
	ProductID string
}

type ProductSnapshot struct {
	ProductID string
	Title     string
	Price     string
	Currency  string
}

type Client interface {
	GetProductSnapshot(ctx context.Context, input ProductLookupInput) (*ProductSnapshot, error)
}
