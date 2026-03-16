package trackeditemlist

import (
	"context"
	"fmt"
	"strings"

	product "github.com/ljj/gugu-api/internal/core/domain/product"
	trackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type Item struct {
	TrackedItem trackeditem.TrackedItem
	Product     product.Product
}

type Service struct {
	trackedItemRepository trackeditem.Repository
	productRepository     product.Repository
}

func NewService(
	trackedItemRepository trackeditem.Repository,
	productRepository product.Repository,
) *Service {
	return &Service{
		trackedItemRepository: trackedItemRepository,
		productRepository:     productRepository,
	}
}

func (s *Service) List(ctx context.Context, userID string) ([]Item, error) {
	trackedItems, err := s.trackedItemRepository.ListByUserID(ctx, strings.TrimSpace(userID))
	if err != nil {
		return nil, fmt.Errorf("list tracked items by user id: %w", err)
	}

	items := make([]Item, 0, len(trackedItems))
	for _, trackedItem := range trackedItems {
		foundProduct, err := s.productRepository.FindByID(ctx, trackedItem.ProductID)
		if err != nil {
			return nil, fmt.Errorf("find product by id: %w", err)
		}
		if foundProduct == nil {
			continue
		}

		items = append(items, Item{
			TrackedItem: trackedItem,
			Product:     *foundProduct,
		})
	}

	return items, nil
}
