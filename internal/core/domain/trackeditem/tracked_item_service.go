package trackeditem

import (
	"context"
	"fmt"
	"strings"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type AddInput struct {
	UserID      string
	ProductID   string
	OriginalURL string
	Currency    string
}

type AddResult struct {
	TrackedItem    TrackedItem
	AlreadyTracked bool
}

type AddTrackedItemInput struct {
	UserID            string
	ProviderCommerce  string
	ExternalProductID string
	OriginalURL       string
	Currency          string
}

type AddTrackedItemResult struct {
	TrackedItem    TrackedItem
	AlreadyTracked bool
}

type Service struct {
	finder          Finder
	writer          Writer
	idGenerator     IDGenerator
	clock           Clock
	productService  *domainproduct.Service
	productProvider domainproduct.ProductProvider
}

func NewService(
	finder Finder,
	writer Writer,
	idGenerator IDGenerator,
	clock Clock,
	productService *domainproduct.Service,
	productProvider domainproduct.ProductProvider,
) *Service {
	return &Service{
		finder:          finder,
		writer:          writer,
		idGenerator:     idGenerator,
		clock:           clock,
		productService:  productService,
		productProvider: productProvider,
	}
}

func (s *Service) AddTrackedItem(ctx context.Context, input AddTrackedItemInput) (*AddTrackedItemResult, error) {
	market := enum.Market(input.ProviderCommerce).Normalize()
	if !market.IsSupported() {
		return nil, coreerror.New(coreerror.UnsupportedMarket)
	}

	product, err := s.resolveProduct(ctx, market, input.ExternalProductID, input.OriginalURL)
	if err != nil {
		return nil, err
	}

	currency := input.Currency
	if currency == "" {
		currency = product.Currency
	}
	if currency == "" {
		currency = "KRW"
	}

	addResult, err := s.Add(ctx, AddInput{
		UserID:      input.UserID,
		ProductID:   product.ID,
		OriginalURL: input.OriginalURL,
		Currency:    currency,
	})
	if err != nil {
		return nil, err
	}

	return &AddTrackedItemResult{
		TrackedItem:    addResult.TrackedItem,
		AlreadyTracked: addResult.AlreadyTracked,
	}, nil
}

func (s *Service) resolveProduct(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*domainproduct.Product, error) {
	found, err := s.productService.FindByMarketAndExternalProductID(ctx, market, externalProductID)
	if err != nil {
		return nil, fmt.Errorf("find product by market and external product id: %w", err)
	}
	if found != nil {
		return found, nil
	}

	newProduct, err := s.productProvider.Provide(ctx, market, externalProductID, originalURL)
	if err != nil {
		return nil, fmt.Errorf("provide product: %w", err)
	}
	if newProduct == nil {
		return nil, coreerror.New(coreerror.ProductNotFound)
	}

	return s.productService.Create(ctx, *newProduct)
}

func (s *Service) Add(ctx context.Context, input AddInput) (*AddResult, error) {
	found, err := s.finder.FindByUserIDAndProductID(ctx, strings.TrimSpace(input.UserID), input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find tracked item by user id and product id: %w", err)
	}
	if found != nil {
		return &AddResult{
			TrackedItem:    *found,
			AlreadyTracked: true,
		}, nil
	}

	trackedItemID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate tracked item id: %w", err)
	}

	currency := input.Currency
	if currency == "" {
		currency = "KRW"
	}

	tracked := TrackedItem{
		ID:          trackedItemID,
		UserID:      strings.TrimSpace(input.UserID),
		ProductID:   input.ProductID,
		OriginalURL: strings.TrimSpace(input.OriginalURL),
		Currency:    currency,
		CreatedAt:   s.clock.Now(),
	}
	if err := s.writer.Create(ctx, tracked); err != nil {
		return nil, fmt.Errorf("create tracked item: %w", err)
	}

	return &AddResult{
		TrackedItem:    tracked,
		AlreadyTracked: false,
	}, nil
}

type TrackedItemWithProduct struct {
	TrackedItem TrackedItem
	Product     domainproduct.Product
}

func (s *Service) ListWithProducts(ctx context.Context, userID string) ([]TrackedItemWithProduct, error) {
	trackedItems, err := s.finder.ListByUserID(ctx, strings.TrimSpace(userID))
	if err != nil {
		return nil, fmt.Errorf("list tracked items by user id: %w", err)
	}

	if len(trackedItems) == 0 {
		return nil, nil
	}

	productIDs := make([]string, len(trackedItems))
	for i, tracked := range trackedItems {
		productIDs[i] = tracked.ProductID
	}

	products, err := s.productService.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("find products by ids: %w", err)
	}

	productMap := make(map[string]domainproduct.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	result := make([]TrackedItemWithProduct, 0, len(trackedItems))
	for _, tracked := range trackedItems {
		product, ok := productMap[tracked.ProductID]
		if !ok {
			continue
		}
		result = append(result, TrackedItemWithProduct{
			TrackedItem: tracked,
			Product:     product,
		})
	}

	return result, nil
}

type TrackedItemDetail struct {
	TrackedItem TrackedItem
	Product     domainproduct.Product
	SKUs        []domainproduct.SKU
}

func (s *Service) GetDetail(ctx context.Context, trackedItemID string, userID string) (*TrackedItemDetail, error) {
	found, err := s.finder.FindByIDAndUserID(ctx, strings.TrimSpace(trackedItemID), strings.TrimSpace(userID))
	if err != nil {
		return nil, fmt.Errorf("find tracked item by id and user id: %w", err)
	}
	if found == nil {
		return nil, coreerror.New(coreerror.TrackedItemNotFound)
	}

	product, err := s.productService.FindByID(ctx, found.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}
	if product == nil {
		return nil, coreerror.New(coreerror.ProductNotFound)
	}

	skus, err := s.productService.FindSKUsByProductID(ctx, found.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find product skus: %w", err)
	}

	return &TrackedItemDetail{
		TrackedItem: *found,
		Product:     *product,
		SKUs:        skus,
	}, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error) {
	return s.finder.ListByUserID(ctx, strings.TrimSpace(userID))
}

func (s *Service) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error) {
	return s.finder.FindByUserIDAndProductID(ctx, strings.TrimSpace(userID), strings.TrimSpace(productID))
}

func (s *Service) SelectSKU(ctx context.Context, trackedItemID string, userID string, skuID string) error {
	trackedItemID = strings.TrimSpace(trackedItemID)
	userID = strings.TrimSpace(userID)
	skuID = strings.TrimSpace(skuID)

	found, err := s.finder.FindByIDAndUserID(ctx, trackedItemID, userID)
	if err != nil {
		return fmt.Errorf("find tracked item by id and user id: %w", err)
	}
	if found == nil {
		return coreerror.New(coreerror.TrackedItemNotFound)
	}

	if err := s.writer.UpdateSKU(ctx, found.ID, found.UserID, skuID); err != nil {
		return fmt.Errorf("update tracked item sku: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, trackedItemID string, userID string) error {
	found, err := s.finder.FindByIDAndUserID(ctx, strings.TrimSpace(trackedItemID), strings.TrimSpace(userID))
	if err != nil {
		return fmt.Errorf("find tracked item by id and user id: %w", err)
	}
	if found == nil {
		return coreerror.New(coreerror.TrackedItemNotFound)
	}

	if err := s.writer.DeleteByIDAndUserID(ctx, found.ID, found.UserID); err != nil {
		return fmt.Errorf("delete tracked item by id and user id: %w", err)
	}

	return nil
}
