package aliexpress

import (
	"context"
	"fmt"
	"testing"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	"github.com/ljj/gugu-api/internal/core/enum"
)

func TestProvide_AffiliateExists_WithDSSKU(t *testing.T) {
	p := NewProvider(
		&stubAffiliateClient{
			result: &clientaliexpress.ProductDetailResult{
				Products: []clientaliexpress.AffiliateProduct{
					{
						ProductID:               1001,
						ProductTitle:            "Test Product",
						ProductMainImageURL:     "https://img.ae/1.jpg",
						ProductDetailURL:        "https://ae.com/item/1001",
						PromotionLink:           "https://s.click.ae/e/promo1001",
						TargetSalePrice:         "9900",
						TargetSalePriceCurrency: "KRW",
					},
				},
			},
		},
		&stubDSClient{
			result: &clientaliexpress.DSProductResult{
				SKUs: []clientaliexpress.DSItemSKUInfo{
					{
						SKUID:          "sku-001",
						ID:             "attr-001",
						SKUPrice:       "12000",
						OfferSalePrice: "9900",
						CurrencyCode:   "KRW",
						Properties: []clientaliexpress.DSSKUPropertyDTO{
							{SKUPropertyID: 14, SKUPropertyName: "Color", SKUPropertyValue: "Red", SKUImage: "https://img.ae/red.jpg"},
							{SKUPropertyID: 5, SKUPropertyName: "Size", SKUPropertyValue: "M"},
						},
					},
					{
						SKUID:          "sku-002",
						ID:             "attr-002",
						SKUPrice:       "12000",
						OfferSalePrice: "9900",
						CurrencyCode:   "KRW",
						Properties: []clientaliexpress.DSSKUPropertyDTO{
							{SKUPropertyID: 14, SKUPropertyName: "Color", SKUPropertyValue: "Blue", SKUImage: "https://img.ae/blue.jpg"},
							{SKUPropertyID: 5, SKUPropertyName: "Size", SKUPropertyValue: "L"},
						},
					},
				},
			},
		},
		&stubTokenProvider{token: "aff-token"},
		&stubTokenProvider{token: "ds-token"},
		"KRW", "KO", "KR",
	)

	result, err := p.Provide(context.Background(), enum.MarketAliExpress, "1001", "https://ae.com/item/1001", "KRW", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result == nil {
		t.Fatal("Provide() returned nil")
	}

	// 상품 정보는 Affiliate에서
	if result.Title != "Test Product" {
		t.Errorf("Title = %q, want %q", result.Title, "Test Product")
	}
	if result.PromotionLink != "https://s.click.ae/e/promo1001" {
		t.Errorf("PromotionLink = %q", result.PromotionLink)
	}
	if result.CollectionSource != "AFFILIATE_API" {
		t.Errorf("CollectionSource = %q, want AFFILIATE_API", result.CollectionSource)
	}
	if result.Language != "KO" {
		t.Errorf("Language = %q, want KO", result.Language)
	}

	// SKU는 DS에서
	if len(result.SKUs) != 2 {
		t.Fatalf("SKUs count = %d, want 2", len(result.SKUs))
	}
	if result.SKUs[0].ExternalSKUID != "sku-001" {
		t.Errorf("SKU[0].ExternalSKUID = %q", result.SKUs[0].ExternalSKUID)
	}
	if result.SKUs[0].Color != "Red" {
		t.Errorf("SKU[0].Color = %q, want Red", result.SKUs[0].Color)
	}
	if result.SKUs[0].Size != "M" {
		t.Errorf("SKU[0].Size = %q, want M", result.SKUs[0].Size)
	}
	if result.SKUs[0].ImageURL != "https://img.ae/red.jpg" {
		t.Errorf("SKU[0].ImageURL = %q", result.SKUs[0].ImageURL)
	}
	if result.SKUs[0].Price != "9900" {
		t.Errorf("SKU[0].Price = %q, want 9900", result.SKUs[0].Price)
	}
	if result.SKUs[0].OriginalPrice != "12000" {
		t.Errorf("SKU[0].OriginalPrice = %q, want 12000", result.SKUs[0].OriginalPrice)
	}
}

func TestProvide_AffiliateNotFound_DSFallback(t *testing.T) {
	p := NewProvider(
		&stubAffiliateClient{result: &clientaliexpress.ProductDetailResult{Products: nil}},
		&stubDSClient{
			result: &clientaliexpress.DSProductResult{
				BaseInfo: clientaliexpress.DSItemBaseInfo{
					ProductID:    2001,
					Subject:      "DS Only Product",
					CurrencyCode: "KRW",
				},
				Multimedia: clientaliexpress.DSMultimediaInfo{
					ImageURLs: "https://img.ae/ds1.jpg;https://img.ae/ds2.jpg",
				},
				SKUs: []clientaliexpress.DSItemSKUInfo{
					{
						SKUID:          "ds-sku-001",
						SKUPrice:       "15000",
						OfferSalePrice: "12000",
						CurrencyCode:   "KRW",
						Properties: []clientaliexpress.DSSKUPropertyDTO{
							{SKUPropertyID: 14, SKUPropertyName: "Color", SKUPropertyValue: "Black"},
						},
					},
				},
			},
		},
		&stubTokenProvider{token: "aff-token"},
		&stubTokenProvider{token: "ds-token"},
		"KRW", "KO", "KR",
	)

	result, err := p.Provide(context.Background(), enum.MarketAliExpress, "2001", "https://ae.com/item/2001", "KRW", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result == nil {
		t.Fatal("Provide() returned nil")
	}

	if result.Title != "DS Only Product" {
		t.Errorf("Title = %q, want DS Only Product", result.Title)
	}
	if result.MainImageURL != "https://img.ae/ds1.jpg" {
		t.Errorf("MainImageURL = %q", result.MainImageURL)
	}
	if result.PromotionLink != "" {
		t.Errorf("PromotionLink = %q, want empty", result.PromotionLink)
	}
	if result.CollectionSource != "DS_API" {
		t.Errorf("CollectionSource = %q, want DS_API", result.CollectionSource)
	}
	if len(result.SKUs) != 1 {
		t.Fatalf("SKUs count = %d, want 1", len(result.SKUs))
	}
	if result.SKUs[0].Color != "Black" {
		t.Errorf("SKU[0].Color = %q, want Black", result.SKUs[0].Color)
	}
}

func TestProvide_DSError_AffiliateAlone(t *testing.T) {
	p := NewProvider(
		&stubAffiliateClient{
			result: &clientaliexpress.ProductDetailResult{
				Products: []clientaliexpress.AffiliateProduct{
					{
						ProductID:               3001,
						ProductTitle:            "Affiliate Only",
						ProductDetailURL:        "https://ae.com/item/3001",
						TargetSalePrice:         "5000",
						TargetSalePriceCurrency: "KRW",
					},
				},
			},
		},
		&stubDSClient{err: fmt.Errorf("ds api timeout")},
		&stubTokenProvider{token: "aff-token"},
		&stubTokenProvider{token: "ds-token"},
		"KRW", "KO", "KR",
	)

	result, err := p.Provide(context.Background(), enum.MarketAliExpress, "3001", "", "KRW", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result == nil {
		t.Fatal("Provide() returned nil")
	}

	if result.Title != "Affiliate Only" {
		t.Errorf("Title = %q", result.Title)
	}
	if result.CollectionSource != "AFFILIATE_API" {
		t.Errorf("CollectionSource = %q", result.CollectionSource)
	}
	// DS 실패 → SKU 없음
	if len(result.SKUs) != 0 {
		t.Errorf("SKUs count = %d, want 0", len(result.SKUs))
	}
}

func TestProvide_BothFail_ReturnsNil(t *testing.T) {
	p := NewProvider(
		&stubAffiliateClient{err: fmt.Errorf("affiliate error")},
		&stubDSClient{err: fmt.Errorf("ds error")},
		&stubTokenProvider{token: "aff-token"},
		&stubTokenProvider{token: "ds-token"},
		"KRW", "KO", "KR",
	)

	result, err := p.Provide(context.Background(), enum.MarketAliExpress, "9999", "", "KRW", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result != nil {
		t.Errorf("Provide() = %+v, want nil", result)
	}
}

func TestProvide_UsesRequestedCurrencyForAffiliateAndDS(t *testing.T) {
	affiliate := &stubAffiliateClient{
		result: &clientaliexpress.ProductDetailResult{
			Products: []clientaliexpress.AffiliateProduct{
				{ProductID: 1001, ProductTitle: "USD Product"},
			},
		},
	}
	ds := &stubDSClient{
		result: &clientaliexpress.DSProductResult{},
	}

	p := NewProvider(
		affiliate,
		ds,
		&stubTokenProvider{token: "aff-token"},
		&stubTokenProvider{token: "ds-token"},
		"KRW", "KO", "KR",
	)

	_, err := p.Provide(context.Background(), enum.MarketAliExpress, "1001", "", "USD", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}

	if affiliate.lastInput.TargetCurrency != "USD" {
		t.Fatalf("affiliate target currency = %q, want USD", affiliate.lastInput.TargetCurrency)
	}
	if affiliate.lastInput.TargetLanguage != "EN" {
		t.Fatalf("affiliate target language = %q, want EN", affiliate.lastInput.TargetLanguage)
	}
	if ds.lastInput.TargetCurrency != "USD" {
		t.Fatalf("ds target currency = %q, want USD", ds.lastInput.TargetCurrency)
	}
	if ds.lastInput.ShipToCountry != "US" {
		t.Fatalf("ds ship_to_country = %q, want US", ds.lastInput.ShipToCountry)
	}
}

func TestProvide_UnsupportedMarket_ReturnsNil(t *testing.T) {
	p := NewProvider(nil, nil, nil, nil, "KRW", "KO", "KR")

	result, err := p.Provide(context.Background(), enum.Market("COUPANG"), "1001", "", "KRW", "")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result != nil {
		t.Errorf("Provide() = %+v, want nil", result)
	}
}

func TestProvide_UsesExplicitLanguageWhenProvided(t *testing.T) {
	affiliate := &stubAffiliateClient{
		result: &clientaliexpress.ProductDetailResult{
			Products: []clientaliexpress.AffiliateProduct{
				{ProductID: 1001, ProductTitle: "Explicit Language Product"},
			},
		},
	}

	p := NewProvider(
		affiliate,
		nil,
		&stubTokenProvider{token: "aff-token"},
		nil,
		"KRW", "KO", "KR",
	)

	result, err := p.Provide(context.Background(), enum.MarketAliExpress, "1001", "", "KRW", "EN")
	if err != nil {
		t.Fatalf("Provide() error = %v", err)
	}
	if result == nil {
		t.Fatal("Provide() returned nil")
	}
	if affiliate.lastInput.TargetLanguage != "EN" {
		t.Fatalf("affiliate target language = %q, want EN", affiliate.lastInput.TargetLanguage)
	}
	if result.Language != "EN" {
		t.Fatalf("result language = %q, want EN", result.Language)
	}
}

func TestExtractSKUProperties(t *testing.T) {
	props := []clientaliexpress.DSSKUPropertyDTO{
		{SKUPropertyID: 14, SKUPropertyName: "Color", SKUPropertyValue: "Red", PropertyValueDefinitionName: "Wine Red", SKUImage: "https://img/red.jpg"},
		{SKUPropertyID: 200, SKUPropertyName: "Ships From", SKUPropertyValue: "China"},
		{SKUPropertyID: 5, SKUPropertyName: "Size", SKUPropertyValue: "XL"},
	}

	color, size, propStr := extractSKUProperties(props)
	if color != "Wine Red" {
		t.Errorf("color = %q, want Wine Red", color)
	}
	if size != "XL" {
		t.Errorf("size = %q, want XL", size)
	}
	if propStr == "" {
		t.Error("propStr is empty")
	}
}

// --- stubs ---

type stubAffiliateClient struct {
	result    *clientaliexpress.ProductDetailResult
	err       error
	lastInput clientaliexpress.ProductDetailInput
}

func (c *stubAffiliateClient) GetAffiliateProductDetail(_ context.Context, input clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error) {
	c.lastInput = input
	return c.result, c.err
}

type stubDSClient struct {
	result    *clientaliexpress.DSProductResult
	err       error
	lastInput clientaliexpress.DSProductInput
}

func (c *stubDSClient) GetDSProduct(_ context.Context, input clientaliexpress.DSProductInput) (*clientaliexpress.DSProductResult, error) {
	c.lastInput = input
	return c.result, c.err
}

type stubTokenProvider struct {
	token string
}

func (p *stubTokenProvider) GetAccessToken(_ context.Context) (string, error) {
	return p.token, nil
}
