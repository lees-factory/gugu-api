package trackeditem

import (
	"context"
	"fmt"
	"strings"

	product "github.com/ljj/gugu-api/internal/core/domain/product"
)

type AddInput struct {
	UserID            string
	OriginalURL       string
	Market            product.Market
	ExternalProductID string
}

type AddResult struct {
	Product        product.Product
	TrackedItem    TrackedItem
	AlreadyTracked bool
}

type Service struct {
	repository       Repository
	idGenerator      IDGenerator
	clock            Clock
	productService   *product.Service
	productCollector product.Collector
}

func NewService(
	repository Repository,
	idGenerator IDGenerator,
	clock Clock,
	productService *product.Service,
	productCollector product.Collector,
) *Service {
	return &Service{
		repository:       repository,
		idGenerator:      idGenerator,
		clock:            clock,
		productService:   productService,
		productCollector: productCollector,
	}
}

func (s *Service) Add(ctx context.Context, input AddInput) (*AddResult, error) {
	if !input.Market.IsSupported() {
		return nil, product.ErrUnsupportedMarket
	}

	item, err := s.productService.FindByMarketAndExternalProductID(ctx, input.Market, input.ExternalProductID)
	if err != nil {
		return nil, err
	}

	if item == nil {
		collected, err := s.productCollector.Collect(ctx, product.CollectInput{
			Market:            input.Market,
			ExternalProductID: input.ExternalProductID,
			OriginalURL:       input.OriginalURL,
		})
		if err != nil {
			return nil, fmt.Errorf("collect product: %w", err)
		}

		item, err = s.productService.Create(ctx, product.CreateInput{
			Market:            collected.Market,
			ExternalProductID: collected.ExternalProductID,
			OriginalURL:       collected.OriginalURL,
			Title:             collected.Title,
			MainImageURL:      collected.MainImageURL,
			CurrentPrice:      collected.CurrentPrice,
			Currency:          collected.Currency,
			ProductURL:        collected.ProductURL,
			CollectionSource:  collected.CollectionSource,
		})
		if err != nil {
			return nil, fmt.Errorf("create product: %w", err)
		}
	}

	found, err := s.repository.FindByUserIDAndProductID(ctx, strings.TrimSpace(input.UserID), item.ID)
	if err != nil {
		return nil, fmt.Errorf("find tracked item by user id and product id: %w", err)
	}
	if found != nil {
		return &AddResult{
			Product:        *item,
			TrackedItem:    *found,
			AlreadyTracked: true,
		}, nil
	}

	trackedItemID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate tracked item id: %w", err)
	}

	tracked := TrackedItem{
		ID:          trackedItemID,
		UserID:      strings.TrimSpace(input.UserID),
		ProductID:   item.ID,
		OriginalURL: strings.TrimSpace(input.OriginalURL),
		CreatedAt:   s.clock.Now(),
	}
	if err := s.repository.Create(ctx, tracked); err != nil {
		return nil, fmt.Errorf("create tracked item: %w", err)
	}

	return &AddResult{
		Product:        *item,
		TrackedItem:    tracked,
		AlreadyTracked: false,
	}, nil
}

func (s *Service) Delete(ctx context.Context, trackedItemID string, userID string) error {
	found, err := s.repository.FindByIDAndUserID(ctx, strings.TrimSpace(trackedItemID), strings.TrimSpace(userID))
	if err != nil {
		return fmt.Errorf("find tracked item by id and user id: %w", err)
	}
	if found == nil {
		return ErrTrackedItemNotFound
	}

	if err := s.repository.DeleteByIDAndUserID(ctx, found.ID, found.UserID); err != nil {
		return fmt.Errorf("delete tracked item by id and user id: %w", err)
	}

	return nil
}
