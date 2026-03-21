package product

import (
	"context"
	"fmt"
	"strings"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

const (
	affiliateCollectionSource = "AFFILIATE_API"
)

type aliExpressProductDetailClient interface {
	GetAffiliateProductDetail(ctx context.Context, input clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error)
	GetAffiliateProductSKUDetail(ctx context.Context, input clientaliexpress.ProductSKUDetailInput) (*clientaliexpress.ProductSKUDetailResult, error)
}

type AccessTokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type AliExpressProductFinder struct {
	client         aliExpressProductDetailClient
	tokenProvider  AccessTokenProvider
	targetCurrency string
	targetLanguage string
	shipToCountry  string
}

func NewAliExpressProductFinder(client aliExpressProductDetailClient, tokenProvider AccessTokenProvider, targetCurrency string, targetLanguage string, shipToCountry string) *AliExpressProductFinder {
	return &AliExpressProductFinder{
		client:         client,
		tokenProvider:  tokenProvider,
		targetCurrency: strings.TrimSpace(targetCurrency),
		targetLanguage: strings.TrimSpace(targetLanguage),
		shipToCountry:  strings.TrimSpace(shipToCountry),
	}
}

func (f *AliExpressProductFinder) Find(ctx context.Context, input CollectInput) (*CollectedProduct, error) {
	if input.Market != MarketAliExpress {
		return nil, nil
	}
	if f.client == nil {
		return nil, fmt.Errorf("aliexpress client is required")
	}

	accessToken := ""
	if f.tokenProvider != nil {
		token, err := f.tokenProvider.GetAccessToken(ctx)
		if err != nil {
			return nil, fmt.Errorf("get aliexpress access token: %w", err)
		}
		accessToken = token
	}

	detailResult, err := f.client.GetAffiliateProductDetail(ctx, clientaliexpress.ProductDetailInput{
		ProductIDs:     []string{input.ExternalProductID},
		TargetCurrency: defaultValue(f.targetCurrency, "USD"),
		TargetLanguage: defaultValue(f.targetLanguage, "EN"),
		AccessToken:    accessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("get aliexpress affiliate product detail: %w", err)
	}
	if detailResult == nil || len(detailResult.Products) == 0 {
		return nil, nil
	}

	detailProduct := detailResult.Products[0]

	skuResult, err := f.client.GetAffiliateProductSKUDetail(ctx, clientaliexpress.ProductSKUDetailInput{
		ProductID:      input.ExternalProductID,
		ShipToCountry:  defaultValue(f.shipToCountry, "US"),
		TargetCurrency: defaultValue(f.targetCurrency, "USD"),
		TargetLanguage: defaultValue(f.targetLanguage, "EN"),
		AccessToken:    accessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("get aliexpress affiliate product sku detail: %w", err)
	}

	title := firstNonEmpty(detailProduct.ProductTitle, skuResult.ItemInfo.Title, skuResult.ItemInfo.EnTitle)
	mainImageURL := firstNonEmpty(detailProduct.ProductMainImageURL, skuResult.ItemInfo.ImageLink, skuResult.ItemInfo.ImageWhite)
	productURL := firstNonEmpty(detailProduct.ProductDetailURL, skuResult.ItemInfo.OriginalLink)
	price := firstNonEmpty(
		firstSKUPrice(skuResult),
		detailProduct.TargetSalePrice,
		detailProduct.SalePrice,
		detailProduct.TargetAppSalePrice,
		detailProduct.AppSalePrice,
	)
	currency := firstNonEmpty(
		firstSKUCurrency(skuResult),
		detailProduct.TargetSalePriceCurrency,
		detailProduct.SalePriceCurrency,
		detailProduct.TargetAppSalePriceCurrency,
		detailProduct.AppSalePriceCurrency,
	)

	return &CollectedProduct{
		Market:            MarketAliExpress,
		ExternalProductID: input.ExternalProductID,
		OriginalURL:       firstNonEmpty(input.OriginalURL, productURL),
		Title:             title,
		MainImageURL:      mainImageURL,
		CurrentPrice:      price,
		Currency:          currency,
		ProductURL:        productURL,
		CollectionSource:  affiliateCollectionSource,
		SKUs:              collectSKUs(skuResult),
	}, nil
}

func collectSKUs(result *clientaliexpress.ProductSKUDetailResult) []CollectedSKU {
	if result == nil || len(result.SKUInfos) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	skus := make([]CollectedSKU, 0, len(result.SKUInfos))
	for _, info := range result.SKUInfos {
		externalID := fmt.Sprintf("%d", info.SKUID)
		if seen[externalID] {
			continue
		}
		seen[externalID] = true

		skuName := strings.TrimSpace(info.Color)
		if size := strings.TrimSpace(info.Size); size != "" {
			if skuName != "" {
				skuName += " / " + size
			} else {
				skuName = size
			}
		}

		skus = append(skus, CollectedSKU{
			ExternalSKUID: externalID,
			SKUName:       skuName,
			Color:         strings.TrimSpace(info.Color),
			Size:          strings.TrimSpace(info.Size),
			Price:         firstNonEmpty(info.SalePriceWithTax, info.PriceWithTax),
			OriginalPrice: strings.TrimSpace(info.PriceWithTax),
			Currency:      strings.TrimSpace(info.Currency),
			ImageURL:      strings.TrimSpace(info.SKUImageLink),
			SKUProperties: strings.TrimSpace(info.SKUProperties),
		})
	}
	return skus
}

func firstSKUPrice(result *clientaliexpress.ProductSKUDetailResult) string {
	if result == nil || len(result.SKUInfos) == 0 {
		return ""
	}
	return firstNonEmpty(result.SKUInfos[0].SalePriceWithTax, result.SKUInfos[0].PriceWithTax)
}

func firstSKUCurrency(result *clientaliexpress.ProductSKUDetailResult) string {
	if result == nil || len(result.SKUInfos) == 0 {
		return ""
	}
	return result.SKUInfos[0].Currency
}

func defaultValue(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
