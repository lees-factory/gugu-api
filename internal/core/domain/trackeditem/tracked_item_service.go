package trackeditem

import (
	"context"
	"fmt"
	"strings"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	"github.com/ljj/gugu-api/internal/core/support/page"
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
	Language          string
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

	currency := normalizeCurrency(input.Currency)
	language := normalizeLanguage(input.Language, currency)

	product, err := s.resolveProduct(ctx, market, input.ExternalProductID, input.OriginalURL, currency, language)
	if err != nil {
		return nil, err
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

func (s *Service) resolveProduct(ctx context.Context, market enum.Market, externalProductID string, originalURL string, currency string, language string) (*domainproduct.Product, error) {
	found, err := s.productService.FindByMarketAndExternalProductID(ctx, market, externalProductID)
	if err != nil {
		return nil, fmt.Errorf("find product by market and external product id: %w", err)
	}
	if found != nil {
		if err := s.ensureVariant(ctx, *found, market, externalProductID, originalURL, currency, language); err != nil {
			return nil, err
		}
		return found, nil
	}

	newProduct, err := s.productProvider.Provide(ctx, market, externalProductID, originalURL, currency, language)
	if err != nil {
		return nil, fmt.Errorf("provide product: %w", err)
	}
	if newProduct == nil {
		return nil, coreerror.New(coreerror.ProductNotFound)
	}

	return s.productService.Create(ctx, *newProduct)
}

func (s *Service) ensureVariant(ctx context.Context, product domainproduct.Product, market enum.Market, externalProductID string, originalURL string, currency string, language string) error {
	variant, err := s.productService.FindVariant(ctx, product.ID, language, currency)
	if err != nil {
		return fmt.Errorf("find product variant: %w", err)
	}
	if variant != nil {
		return nil
	}

	newProduct, err := s.productProvider.Provide(ctx, market, externalProductID, originalURL, currency, language)
	if err != nil {
		return fmt.Errorf("provide product variant: %w", err)
	}
	if newProduct == nil {
		return coreerror.New(coreerror.ProductNotFound)
	}

	if err := s.productService.UpsertVariant(ctx, product.ID, *newProduct); err != nil {
		return fmt.Errorf("upsert product variant: %w", err)
	}
	return nil
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
	Variant     *domainproduct.Variant
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
	variantMap, err := s.findVariantsForTrackedItems(ctx, trackedItems)
	if err != nil {
		return nil, fmt.Errorf("find product variants: %w", err)
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
			Variant:     variantMap[variantMapKey(product.ID, tracked.Currency)],
		})
	}

	return result, nil
}

func (s *Service) ListWithProductsCursor(ctx context.Context, userID string, cursor page.CursorRequest) (*page.CursorPage[TrackedItemWithProduct], error) {
	userID = strings.TrimSpace(userID)
	fetchSize := cursor.EffectiveSize() + 1

	var trackedItems []TrackedItem
	var err error

	if cursor.Cursor == "" {
		trackedItems, err = s.finder.ListByUserIDFirstPage(ctx, userID, fetchSize)
	} else {
		cursorCreatedAt, cursorID, decodeErr := page.DecodeCursor(cursor.Cursor)
		if decodeErr != nil {
			return nil, fmt.Errorf("invalid cursor: %w", decodeErr)
		}
		trackedItems, err = s.finder.ListByUserIDWithCursor(ctx, userID, cursorCreatedAt, cursorID, fetchSize)
	}
	if err != nil {
		return nil, fmt.Errorf("list tracked items: %w", err)
	}

	hasMore := len(trackedItems) > cursor.EffectiveSize()
	if hasMore {
		trackedItems = trackedItems[:cursor.EffectiveSize()]
	}

	if len(trackedItems) == 0 {
		return &page.CursorPage[TrackedItemWithProduct]{
			Items:   []TrackedItemWithProduct{},
			HasMore: false,
		}, nil
	}

	productIDs := make([]string, len(trackedItems))
	for i, tracked := range trackedItems {
		productIDs[i] = tracked.ProductID
	}

	products, err := s.productService.FindByIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("find products by ids: %w", err)
	}
	variantMap, err := s.findVariantsForTrackedItems(ctx, trackedItems)
	if err != nil {
		return nil, fmt.Errorf("find product variants: %w", err)
	}

	productMap := make(map[string]domainproduct.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	items := make([]TrackedItemWithProduct, 0, len(trackedItems))
	for _, tracked := range trackedItems {
		product, ok := productMap[tracked.ProductID]
		if !ok {
			continue
		}
		items = append(items, TrackedItemWithProduct{
			TrackedItem: tracked,
			Product:     product,
			Variant:     variantMap[variantMapKey(product.ID, tracked.Currency)],
		})
	}

	var nextCursor string
	if hasMore && len(items) > 0 {
		last := items[len(items)-1].TrackedItem
		nextCursor = page.EncodeCursor(last.CreatedAt, last.ID)
	}

	return &page.CursorPage[TrackedItemWithProduct]{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

type TrackedItemDetail struct {
	TrackedItem TrackedItem
	Product     domainproduct.Product
	Variant     *domainproduct.Variant
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
		Variant:     s.findVariantForTrackedItem(ctx, product.ID, found.Currency),
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

func normalizeCurrency(currency string) string {
	currency = strings.ToUpper(strings.TrimSpace(currency))
	if currency == "" {
		return "KRW"
	}
	return currency
}

func normalizeLanguage(language string, currency string) string {
	language = strings.ToUpper(strings.TrimSpace(language))
	if language != "" {
		return language
	}

	switch normalizeCurrency(currency) {
	case "KRW":
		return "KO"
	default:
		return "EN"
	}
}

func (s *Service) findVariantForTrackedItem(ctx context.Context, productID string, currency string) *domainproduct.Variant {
	language := normalizeLanguage("", currency)
	variant, err := s.productService.FindVariant(ctx, productID, language, currency)
	if err != nil {
		return nil
	}
	return variant
}

func (s *Service) findVariantsForTrackedItems(ctx context.Context, trackedItems []TrackedItem) (map[string]*domainproduct.Variant, error) {
	keys := make([]domainproduct.VariantLookupKey, 0, len(trackedItems))
	for _, trackedItem := range trackedItems {
		keys = append(keys, domainproduct.VariantLookupKey{
			ProductID: trackedItem.ProductID,
			Language:  normalizeLanguage("", trackedItem.Currency),
			Currency:  normalizeCurrency(trackedItem.Currency),
		})
	}

	variants, err := s.productService.FindVariants(ctx, keys)
	if err != nil {
		return nil, err
	}

	variantMap := make(map[string]*domainproduct.Variant, len(variants))
	for i := range variants {
		variant := variants[i]
		variantMap[variantMapKey(variant.ProductID, variant.Currency)] = &variant
	}
	return variantMap, nil
}

func variantMapKey(productID string, currency string) string {
	return strings.TrimSpace(productID) + ":" + normalizeCurrency(currency)
}
