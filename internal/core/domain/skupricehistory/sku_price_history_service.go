package skupricehistory

import "context"

type Service struct {
	finder Finder
}

func NewService(finder Finder) *Service {
	return &Service{finder: finder}
}

func (s *Service) ListBySKUID(ctx context.Context, skuID string) ([]SKUPriceHistory, error) {
	return s.finder.ListBySKUID(ctx, skuID)
}
