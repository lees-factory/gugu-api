package product

import (
	"context"
	"fmt"
	"strings"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	"github.com/ljj/gugu-api/internal/core/enum"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type Service struct {
	finder              Finder
	writer              Writer
	skuRepository       SKURepository
	idGenerator         IDGenerator
	clock               Clock
	priceHistoryWriter  domainpricehistory.Writer
	skuPriceHistWriter  domainsph.Writer
}

func NewService(finder Finder, writer Writer, skuRepository SKURepository, idGenerator IDGenerator, clock Clock, priceHistoryWriter domainpricehistory.Writer, skuPriceHistWriter domainsph.Writer) *Service {
	return &Service{
		finder:              finder,
		writer:              writer,
		skuRepository:       skuRepository,
		idGenerator:         idGenerator,
		clock:               clock,
		priceHistoryWriter:  priceHistoryWriter,
		skuPriceHistWriter:  skuPriceHistWriter,
	}
}

func (s *Service) FindByID(ctx context.Context, productID string) (*Product, error) {
	found, err := s.finder.FindByID(ctx, strings.TrimSpace(productID))
	if err != nil {
		return nil, fmt.Errorf("find product: %w", err)
	}
	if found == nil {
		return nil, coreerror.New(coreerror.ProductNotFound)
	}
	return found, nil
}

func (s *Service) FindByIDs(ctx context.Context, productIDs []string) ([]Product, error) {
	trimmed := make([]string, len(productIDs))
	for i, id := range productIDs {
		trimmed[i] = strings.TrimSpace(id)
	}
	return s.finder.FindByIDs(ctx, trimmed)
}

func (s *Service) FindSKUsByProductID(ctx context.Context, productID string) ([]SKU, error) {
	return s.skuRepository.FindByProductID(ctx, strings.TrimSpace(productID))
}

func (s *Service) FindByMarketAndExternalProductID(ctx context.Context, market enum.Market, externalProductID string) (*Product, error) {
	return s.finder.FindByMarketAndExternalProductID(ctx, market, strings.TrimSpace(externalProductID))
}

func (s *Service) ListByMarket(ctx context.Context, market enum.Market) ([]Product, error) {
	return s.finder.ListByMarket(ctx, market)
}

func (s *Service) UpdatePrice(ctx context.Context, productID string, price string, currency string) error {
	found, err := s.finder.FindByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("find product: %w", err)
	}
	if found == nil {
		return fmt.Errorf("product not found: %s", productID)
	}

	now := s.clock.Now()
	found.CurrentPrice = strings.TrimSpace(price)
	found.Currency = strings.TrimSpace(currency)
	found.LastCollectedAt = now
	found.UpdatedAt = now

	if err := s.writer.Update(ctx, *found); err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

func (s *Service) EnrichSKUs(ctx context.Context, productID string, skuInputs []NewSKU) (int, error) {
	found, err := s.finder.FindByID(ctx, strings.TrimSpace(productID))
	if err != nil {
		return 0, fmt.Errorf("find product: %w", err)
	}
	if found == nil {
		return 0, coreerror.New(coreerror.ProductNotFound)
	}

	now := s.clock.Now()
	added := 0

	for _, skuInput := range skuInputs {
		externalSKUID := strings.TrimSpace(skuInput.ExternalSKUID)
		existing, err := s.skuRepository.FindByProductIDAndExternalSKUID(ctx, productID, externalSKUID)
		if err != nil {
			return added, fmt.Errorf("find existing sku: %w", err)
		}

		if existing != nil {
			existing.SKUName = strings.TrimSpace(skuInput.SKUName)
			existing.Color = strings.TrimSpace(skuInput.Color)
			existing.Size = strings.TrimSpace(skuInput.Size)
			existing.Price = strings.TrimSpace(skuInput.Price)
			existing.OriginalPrice = strings.TrimSpace(skuInput.OriginalPrice)
			existing.Currency = strings.TrimSpace(skuInput.Currency)
			existing.ImageURL = strings.TrimSpace(skuInput.ImageURL)
			existing.SKUProperties = strings.TrimSpace(skuInput.SKUProperties)
			existing.UpdatedAt = now
			if err := s.skuRepository.Upsert(ctx, *existing); err != nil {
				return added, fmt.Errorf("upsert existing sku: %w", err)
			}
			continue
		}

		skuID, err := s.idGenerator.New()
		if err != nil {
			return added, fmt.Errorf("generate sku id: %w", err)
		}

		sku := SKU{
			ID:            skuID,
			ProductID:     productID,
			ExternalSKUID: externalSKUID,
			OriginSKUID:   strings.TrimSpace(skuInput.OriginSKUID),
			SKUName:       strings.TrimSpace(skuInput.SKUName),
			Color:         strings.TrimSpace(skuInput.Color),
			Size:          strings.TrimSpace(skuInput.Size),
			Price:         strings.TrimSpace(skuInput.Price),
			OriginalPrice: strings.TrimSpace(skuInput.OriginalPrice),
			Currency:      strings.TrimSpace(skuInput.Currency),
			ImageURL:      strings.TrimSpace(skuInput.ImageURL),
			SKUProperties: strings.TrimSpace(skuInput.SKUProperties),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.skuRepository.Create(ctx, sku); err != nil {
			return added, fmt.Errorf("create new sku: %w", err)
		}

		if s.skuPriceHistWriter != nil && sku.Price != "" {
			if err := s.skuPriceHistWriter.Create(ctx, domainsph.SKUPriceHistory{
				SKUID:       sku.ID,
				Price:       sku.Price,
				Currency:    sku.Currency,
				RecordedAt:  now,
				ChangeValue: "0",
			}); err != nil {
				return added, fmt.Errorf("create sku price history: %w", err)
			}
		}
		added++
	}

	return added, nil
}

func (s *Service) Create(ctx context.Context, input NewProduct) (*Product, error) {
	productID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate product id: %w", err)
	}

	now := s.clock.Now()
	item := Product{
		ID:                productID,
		Market:            input.Market,
		ExternalProductID: strings.TrimSpace(input.ExternalProductID),
		OriginalURL:       strings.TrimSpace(input.OriginalURL),
		Title:             strings.TrimSpace(input.Title),
		MainImageURL:      strings.TrimSpace(input.MainImageURL),
		CurrentPrice:      strings.TrimSpace(input.CurrentPrice),
		Currency:          strings.TrimSpace(input.Currency),
		ProductURL:        strings.TrimSpace(input.ProductURL),
		PromotionLink:     strings.TrimSpace(input.PromotionLink),
		CollectionSource:  strings.TrimSpace(input.CollectionSource),
		LastCollectedAt:   now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.writer.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	if s.priceHistoryWriter != nil && item.CurrentPrice != "" {
		if err := s.priceHistoryWriter.Create(ctx, domainpricehistory.PriceHistory{
			ProductID:   item.ID,
			Price:       item.CurrentPrice,
			Currency:    item.Currency,
			RecordedAt:  now,
			ChangeValue: "0",
		}); err != nil {
			return nil, fmt.Errorf("create product price history: %w", err)
		}
	}

	for _, skuInput := range input.SKUs {
		skuID, err := s.idGenerator.New()
		if err != nil {
			return nil, fmt.Errorf("generate sku id: %w", err)
		}
		sku := SKU{
			ID:            skuID,
			ProductID:     item.ID,
			ExternalSKUID: strings.TrimSpace(skuInput.ExternalSKUID),
			OriginSKUID:   strings.TrimSpace(skuInput.OriginSKUID),
			SKUName:       strings.TrimSpace(skuInput.SKUName),
			Color:         strings.TrimSpace(skuInput.Color),
			Size:          strings.TrimSpace(skuInput.Size),
			Price:         strings.TrimSpace(skuInput.Price),
			OriginalPrice: strings.TrimSpace(skuInput.OriginalPrice),
			Currency:      strings.TrimSpace(skuInput.Currency),
			ImageURL:      strings.TrimSpace(skuInput.ImageURL),
			SKUProperties: strings.TrimSpace(skuInput.SKUProperties),
			CreatedAt:     now,
			UpdatedAt:     now,
		}
		if err := s.skuRepository.Create(ctx, sku); err != nil {
			return nil, fmt.Errorf("create product sku: %w", err)
		}

		if s.skuPriceHistWriter != nil && sku.Price != "" {
			if err := s.skuPriceHistWriter.Create(ctx, domainsph.SKUPriceHistory{
				SKUID:       sku.ID,
				Price:       sku.Price,
				Currency:    sku.Currency,
				RecordedAt:  now,
				ChangeValue: "0",
			}); err != nil {
				return nil, fmt.Errorf("create sku price history: %w", err)
			}
		}
	}

	return &item, nil
}
