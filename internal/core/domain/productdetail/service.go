package productdetail

import (
	"context"
	"fmt"
	"strings"

	pricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	product "github.com/ljj/gugu-api/internal/core/domain/product"
	trackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
)

type Detail struct {
	Product         product.Product
	PriceHistories  []pricehistory.PriceHistory
	IsTrackedByUser bool
	TrackedItemID   string
}

type Service struct {
	productRepository      product.Repository
	priceHistoryRepository pricehistory.Repository
	trackedItemRepository  trackeditem.Repository
}

func NewService(
	productRepository product.Repository,
	priceHistoryRepository pricehistory.Repository,
	trackedItemRepository trackeditem.Repository,
) *Service {
	return &Service{
		productRepository:      productRepository,
		priceHistoryRepository: priceHistoryRepository,
		trackedItemRepository:  trackedItemRepository,
	}
}

func (s *Service) Get(ctx context.Context, productID string, userID string) (*Detail, error) {
	item, err := s.productRepository.FindByID(ctx, strings.TrimSpace(productID))
	if err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}
	if item == nil {
		return nil, product.ErrProductNotFound
	}

	histories, err := s.priceHistoryRepository.ListByProductID(ctx, item.ID)
	if err != nil {
		return nil, fmt.Errorf("list price histories by product id: %w", err)
	}

	isTrackedByUser := false
	trackedItemID := ""
	if strings.TrimSpace(userID) != "" {
		tracked, err := s.trackedItemRepository.FindByUserIDAndProductID(ctx, strings.TrimSpace(userID), item.ID)
		if err != nil {
			return nil, fmt.Errorf("find tracked item by user id and product id: %w", err)
		}
		if tracked != nil {
			isTrackedByUser = true
			trackedItemID = tracked.ID
		}
	}

	return &Detail{
		Product:         *item,
		PriceHistories:  histories,
		IsTrackedByUser: isTrackedByUser,
		TrackedItemID:   trackedItemID,
	}, nil
}
