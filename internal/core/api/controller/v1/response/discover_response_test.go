package response

import (
	"testing"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

func TestNewHotProductList(t *testing.T) {
	items := NewHotProductList([]domainproduct.Product{
		{
			ID:           "p1",
			Market:       enum.MarketAliExpress,
			Title:        "product-1",
			MainImageURL: "https://example.com/1.jpg",
			ProductURL:   "https://example.com/p1",
		},
		{
			ID:           "p2",
			Market:       enum.MarketAliExpress,
			Title:        "product-2",
			MainImageURL: "https://example.com/2.jpg",
			ProductURL:   "https://example.com/p2",
		},
	}, map[string]*domainproduct.Variant{
		"p2": {
			ProductID:    "p2",
			Language:     "EN",
			Currency:     "USD",
			Title:        "variant-product-2",
			MainImageURL: "https://example.com/2-variant.jpg",
			ProductURL:   "https://example.com/p2-en",
		},
	}, "KRW")

	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want 2", len(items))
	}
	if got := items[0].Currency; got != "KRW" {
		t.Fatalf("first item currency = %q, want KRW", got)
	}
	if got := items[1].ProductURL; got != "https://example.com/p2-en" {
		t.Fatalf("second item product_url = %q", got)
	}
	if got := items[1].Currency; got != "USD" {
		t.Fatalf("second item currency = %q, want USD", got)
	}
}
