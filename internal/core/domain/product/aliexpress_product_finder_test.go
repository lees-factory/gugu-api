package product

import (
	"context"
	"testing"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

func TestAliExpressProductFinderFind(t *testing.T) {
	finder := NewAliExpressProductFinder(&stubAliExpressProductClient{
		productDetailResult: &clientaliexpress.ProductDetailResult{
			CurrentRecordCount: 1,
			Products: []clientaliexpress.AffiliateProduct{
				{
					ProductID:                  1005001234567890,
					ProductTitle:               "Keyboard",
					TargetSalePrice:            "89.99",
					TargetSalePriceCurrency:    "USD",
					ProductMainImageURL:        "https://img.example/1.jpg",
					ProductDetailURL:           "https://www.aliexpress.com/item/1005001234567890.html",
					TargetOriginalPrice:        "99.99",
					TargetAppSalePrice:         "88.99",
					TargetAppSalePriceCurrency: "USD",
				},
			},
		},
		productSKUDetailResult: &clientaliexpress.ProductSKUDetailResult{
			Code:    0,
			Success: true,
			ItemInfo: clientaliexpress.AffiliateSKUItemInfo{
				ProductID:    "1005001234567890",
				Title:        "Keyboard",
				OriginalLink: "https://www.aliexpress.com/item/1005001234567890.html",
				ImageLink:    "https://img.example/1.jpg",
			},
			SKUInfos: []clientaliexpress.AffiliateSKUInfo{
				{
					SKUID:            20001,
					Currency:         "USD",
					SalePriceWithTax: "87.99",
				},
			},
		},
	}, "USD", "EN", "KR")

	result, err := finder.Find(context.Background(), CollectInput{
		Market:            MarketAliExpress,
		ExternalProductID: "1005001234567890",
		OriginalURL:       "https://www.aliexpress.com/item/1005001234567890.html",
	})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if result == nil {
		t.Fatal("Find() result = nil")
	}
	if result.Title != "Keyboard" {
		t.Fatalf("Title = %q", result.Title)
	}
	if result.CurrentPrice != "87.99" {
		t.Fatalf("CurrentPrice = %q", result.CurrentPrice)
	}
	if result.CollectionSource != affiliateCollectionSource {
		t.Fatalf("CollectionSource = %q", result.CollectionSource)
	}
}

func TestAliExpressProductFinderFindSkipsOtherMarket(t *testing.T) {
	finder := NewAliExpressProductFinder(&stubAliExpressProductClient{}, "USD", "EN", "KR")

	result, err := finder.Find(context.Background(), CollectInput{
		Market:            MarketEBay,
		ExternalProductID: "item-1",
		OriginalURL:       "https://www.ebay.com/itm/item-1",
	})
	if err != nil {
		t.Fatalf("Find() error = %v", err)
	}
	if result != nil {
		t.Fatalf("Find() result = %#v, want nil", result)
	}
}

type stubAliExpressProductClient struct {
	productDetailResult    *clientaliexpress.ProductDetailResult
	productDetailErr       error
	productSKUDetailResult *clientaliexpress.ProductSKUDetailResult
	productSKUDetailErr    error
}

func (c *stubAliExpressProductClient) GetAffiliateProductDetail(_ context.Context, _ clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error) {
	return c.productDetailResult, c.productDetailErr
}

func (c *stubAliExpressProductClient) GetAffiliateProductSKUDetail(_ context.Context, _ clientaliexpress.ProductSKUDetailInput) (*clientaliexpress.ProductSKUDetailResult, error) {
	return c.productSKUDetailResult, c.productSKUDetailErr
}
