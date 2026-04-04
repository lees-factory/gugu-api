package response

import (
	"testing"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

func TestNewHotProductListTrimsPriceFraction(t *testing.T) {
	items := NewHotProductList([]domainproduct.Product{
		{
			ID:           "p1",
			Market:       enum.MarketAliExpress,
			Title:        "product-1",
			MainImageURL: "https://example.com/1.jpg",
			CurrentPrice: "8335.00",
			Currency:     "KRW",
			ProductURL:   "https://example.com/p1",
		},
		{
			ID:           "p2",
			Market:       enum.MarketAliExpress,
			Title:        "product-2",
			MainImageURL: "https://example.com/2.jpg",
			CurrentPrice: "12999.95",
			Currency:     "KRW",
			ProductURL:   "https://example.com/p2",
		},
		{
			ID:           "p3",
			Market:       enum.MarketAliExpress,
			Title:        "product-3",
			MainImageURL: "https://example.com/3.jpg",
			CurrentPrice: "5400",
			Currency:     "KRW",
			ProductURL:   "https://example.com/p3",
		},
	})

	if got := items[0].CurrentPrice; got != "8335" {
		t.Fatalf("first item current_price = %q, want %q", got, "8335")
	}
	if got := items[1].CurrentPrice; got != "12999" {
		t.Fatalf("second item current_price = %q, want %q", got, "12999")
	}
	if got := items[2].CurrentPrice; got != "5400" {
		t.Fatalf("third item current_price = %q, want %q", got, "5400")
	}
}
