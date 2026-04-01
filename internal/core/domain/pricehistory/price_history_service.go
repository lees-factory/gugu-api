package pricehistory

import "context"

type Service struct {
	finder Finder
}

func NewService(finder Finder) *Service {
	return &Service{finder: finder}
}

func (s *Service) ListByProductID(ctx context.Context, productID string, currency string) ([]PriceHistory, error) {
	return s.finder.ListByProductID(ctx, productID, currency)
}
