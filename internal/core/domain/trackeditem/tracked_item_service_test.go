package trackeditem

import (
	"context"
	"testing"
	"time"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
)

func TestNormalizeCurrency_DefaultsToKRW(t *testing.T) {
	if got := normalizeCurrency(""); got != "KRW" {
		t.Fatalf("normalizeCurrency(\"\") = %q, want KRW", got)
	}
}

func TestNormalizeLanguage_UsesExplicitValue(t *testing.T) {
	if got := normalizeLanguage("en", "KRW"); got != "EN" {
		t.Fatalf("normalizeLanguage(\"en\", \"KRW\") = %q, want EN", got)
	}
}

func TestNormalizeLanguage_FallsBackFromCurrency(t *testing.T) {
	if got := normalizeLanguage("", "KRW"); got != "KO" {
		t.Fatalf("normalizeLanguage(\"\", \"KRW\") = %q, want KO", got)
	}
	if got := normalizeLanguage("", "USD"); got != "EN" {
		t.Fatalf("normalizeLanguage(\"\", \"USD\") = %q, want EN", got)
	}
}

func TestResolveProduct_CreatesMissingVariantForExistingProduct(t *testing.T) {
	productRepo := memoryproduct.NewRepository()
	variantRepo := memoryproduct.NewVariantRepository()
	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepo),
		domainproduct.NewWriter(productRepo),
		variantRepo,
		memoryproduct.NewSKURepository(),
		testIDGenerator{},
		testClock{now: time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC)},
		nil,
		nil,
	)

	product := domainproduct.Product{
		ID:                "product-1",
		Market:            enum.MarketAliExpress,
		ExternalProductID: "1001",
		OriginalURL:       "https://example.com/item/1001",
		Title:             "Base Product",
		CreatedAt:         time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		LastCollectedAt:   time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
	}
	if err := productRepo.Create(context.Background(), product); err != nil {
		t.Fatalf("create product: %v", err)
	}

	provider := &stubProductProvider{
		result: &domainproduct.NewProduct{
			Market:            enum.MarketAliExpress,
			ExternalProductID: "1001",
			OriginalURL:       "https://example.com/item/1001",
			Language:          "EN",
			Title:             "English Product",
			MainImageURL:      "https://img.example.com/1001.jpg",
			CurrentPrice:      "15.99",
			Currency:          "USD",
			ProductURL:        "https://example.com/en/item/1001",
		},
	}

	service := NewService(nil, nil, testIDGenerator{}, testClock{}, productService, provider)

	found, err := service.resolveProduct(context.Background(), enum.MarketAliExpress, "1001", "https://example.com/item/1001", "USD", "EN")
	if err != nil {
		t.Fatalf("resolveProduct() error = %v", err)
	}
	if found == nil {
		t.Fatal("resolveProduct() returned nil")
	}
	if provider.calls != 1 {
		t.Fatalf("provider calls = %d, want 1", provider.calls)
	}

	variant, err := productService.FindVariant(context.Background(), "product-1", "EN", "USD")
	if err != nil {
		t.Fatalf("FindVariant() error = %v", err)
	}
	if variant == nil {
		t.Fatal("variant = nil, want saved variant")
	}
	if variant.Title != "English Product" {
		t.Fatalf("variant.Title = %q, want English Product", variant.Title)
	}
}

func TestResolveProduct_SkipsProviderWhenVariantAlreadyExists(t *testing.T) {
	productRepo := memoryproduct.NewRepository()
	variantRepo := memoryproduct.NewVariantRepository()
	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepo),
		domainproduct.NewWriter(productRepo),
		variantRepo,
		memoryproduct.NewSKURepository(),
		testIDGenerator{},
		testClock{now: time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC)},
		nil,
		nil,
	)

	product := domainproduct.Product{
		ID:                "product-1",
		Market:            enum.MarketAliExpress,
		ExternalProductID: "1001",
		CreatedAt:         time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		LastCollectedAt:   time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
	}
	if err := productRepo.Create(context.Background(), product); err != nil {
		t.Fatalf("create product: %v", err)
	}
	if err := variantRepo.Upsert(context.Background(), domainproduct.Variant{
		ProductID:    "product-1",
		Language:     "EN",
		Currency:     "USD",
		Title:        "Existing Variant",
		CreatedAt:    time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2026, 4, 4, 12, 0, 0, 0, time.UTC),
		CurrentPrice: "15.99",
	}); err != nil {
		t.Fatalf("upsert variant: %v", err)
	}

	provider := &stubProductProvider{
		result: &domainproduct.NewProduct{Language: "EN", Currency: "USD"},
	}

	service := NewService(nil, nil, testIDGenerator{}, testClock{}, productService, provider)

	_, err := service.resolveProduct(context.Background(), enum.MarketAliExpress, "1001", "https://example.com/item/1001", "USD", "EN")
	if err != nil {
		t.Fatalf("resolveProduct() error = %v", err)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls = %d, want 0", provider.calls)
	}
}

type stubProductProvider struct {
	result *domainproduct.NewProduct
	err    error
	calls  int
}

func (p *stubProductProvider) Provide(_ context.Context, _ enum.Market, _ string, _ string, _ string, _ string) (*domainproduct.NewProduct, error) {
	p.calls++
	return p.result, p.err
}

type testIDGenerator struct{}

func (testIDGenerator) New() (string, error) {
	return "generated-id", nil
}

type testClock struct {
	now time.Time
}

func (c testClock) Now() time.Time {
	return c.now
}
