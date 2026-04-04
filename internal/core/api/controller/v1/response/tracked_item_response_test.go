package response

import (
	"testing"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	"github.com/ljj/gugu-api/internal/core/enum"
)

func TestNewListTrackedItem_UsesVariantWhenPresent(t *testing.T) {
	item := NewListTrackedItem(domaintrackeditem.TrackedItemWithProduct{
		TrackedItem: domaintrackeditem.TrackedItem{
			ID:          "tracked-1",
			ProductID:   "product-1",
			Currency:    "USD",
			OriginalURL: "https://origin.example.com/item/1",
		},
		Product: domainproduct.Product{
			ID:                "product-1",
			Market:            enum.MarketAliExpress,
			ExternalProductID: "1001",
			Title:             "Base Title",
			MainImageURL:      "https://img.example.com/base.jpg",
			ProductURL:        "https://example.com/base",
		},
		Variant: &domainproduct.Variant{
			ProductID:    "product-1",
			Language:     "EN",
			Currency:     "USD",
			Title:        "Variant Title",
			MainImageURL: "https://img.example.com/variant.jpg",
			ProductURL:   "https://example.com/variant",
		},
	})

	if item.Title != "Variant Title" {
		t.Fatalf("Title = %q, want Variant Title", item.Title)
	}
	if item.MainImageURL != "https://img.example.com/variant.jpg" {
		t.Fatalf("MainImageURL = %q", item.MainImageURL)
	}
	if item.ProductURL != "https://example.com/variant" {
		t.Fatalf("ProductURL = %q", item.ProductURL)
	}
}

func TestNewTrackedItemDetail_FallsBackToProductWhenVariantMissing(t *testing.T) {
	item := NewTrackedItemDetail(&domaintrackeditem.TrackedItemDetail{
		TrackedItem: domaintrackeditem.TrackedItem{
			ID:          "tracked-1",
			ProductID:   "product-1",
			Currency:    "KRW",
			OriginalURL: "https://origin.example.com/item/1",
		},
		Product: domainproduct.Product{
			ID:                "product-1",
			Market:            enum.MarketAliExpress,
			ExternalProductID: "1001",
			Title:             "Base Title",
			MainImageURL:      "https://img.example.com/base.jpg",
			ProductURL:        "https://example.com/base",
		},
	})

	if item.Title != "Base Title" {
		t.Fatalf("Title = %q, want Base Title", item.Title)
	}
	if item.MainImageURL != "https://img.example.com/base.jpg" {
		t.Fatalf("MainImageURL = %q", item.MainImageURL)
	}
	if item.ProductURL != "https://example.com/base" {
		t.Fatalf("ProductURL = %q", item.ProductURL)
	}
}
