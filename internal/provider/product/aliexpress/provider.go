package aliexpress

import (
	"context"
	"fmt"
	"strings"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	clientcrawler "github.com/ljj/gugu-api/internal/clients/crawler"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

const (
	collectionSource = "AFFILIATE_API"
	skuBatchSize     = 20
)

type ProductDetailClient interface {
	GetAffiliateProductDetail(ctx context.Context, input clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error)
	GetAffiliateProductSKUDetail(ctx context.Context, input clientaliexpress.ProductSKUDetailInput) (*clientaliexpress.ProductSKUDetailResult, error)
}

type TokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type Provider struct {
	client         ProductDetailClient
	tokenProvider  TokenProvider
	crawler        clientcrawler.Client
	targetCurrency string
	targetLanguage string
	shipToCountry  string
}

func NewProvider(client ProductDetailClient, tokenProvider TokenProvider, crawler clientcrawler.Client, targetCurrency string, targetLanguage string, shipToCountry string) *Provider {
	return &Provider{
		client:         client,
		tokenProvider:  tokenProvider,
		crawler:        crawler,
		targetCurrency: strings.TrimSpace(targetCurrency),
		targetLanguage: strings.TrimSpace(targetLanguage),
		shipToCountry:  strings.TrimSpace(shipToCountry),
	}
}

func (p *Provider) Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*domainproduct.NewProduct, error) {
	if market != enum.MarketAliExpress {
		return nil, nil
	}
	if p.client == nil {
		return nil, fmt.Errorf("aliexpress client is required")
	}

	accessToken, err := p.resolveAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	detailResult, err := p.client.GetAffiliateProductDetail(ctx, clientaliexpress.ProductDetailInput{
		ProductIDs:     []string{externalProductID},
		TargetCurrency: defaultValue(p.targetCurrency, "USD"),
		TargetLanguage: defaultValue(p.targetLanguage, "EN"),
		AccessToken:    accessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("get aliexpress affiliate product detail: %w", err)
	}
	if detailResult == nil || len(detailResult.Products) == 0 {
		return nil, nil
	}

	detailProduct := detailResult.Products[0]

	allSKUIDs := p.crawlSKUIDs(ctx, originalURL, externalProductID)

	skuResult, err := p.fetchAllSKUs(ctx, externalProductID, accessToken, allSKUIDs)
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

	return &domainproduct.NewProduct{
		Market:            enum.MarketAliExpress,
		ExternalProductID: externalProductID,
		OriginalURL:       firstNonEmpty(originalURL, productURL),
		Title:             title,
		MainImageURL:      mainImageURL,
		CurrentPrice:      price,
		Currency:          currency,
		ProductURL:        productURL,
		CollectionSource:  collectionSource,
		SKUs:              buildSKUs(skuResult),
	}, nil
}

func (p *Provider) crawlSKUIDs(ctx context.Context, originalURL string, externalProductID string) []string {
	if p.crawler == nil || originalURL == "" {
		return nil
	}

	crawlURL := originalURL
	if crawlURL == "" {
		crawlURL = fmt.Sprintf("https://www.aliexpress.com/item/%s.html", externalProductID)
	}

	result, err := p.crawler.Crawl(ctx, clientcrawler.CrawlInput{URL: crawlURL})
	if err != nil || result == nil {
		return nil
	}

	skuIDs := make([]string, 0, len(result.SKUs))
	seen := make(map[string]bool)
	for _, sku := range result.SKUs {
		id := strings.TrimSpace(sku.ExternalSKUID)
		if id != "" && !seen[id] {
			seen[id] = true
			skuIDs = append(skuIDs, id)
		}
	}
	return skuIDs
}

func (p *Provider) fetchAllSKUs(ctx context.Context, externalProductID string, accessToken string, allSKUIDs []string) (*clientaliexpress.ProductSKUDetailResult, error) {
	baseInput := clientaliexpress.ProductSKUDetailInput{
		ProductID:      externalProductID,
		ShipToCountry:  defaultValue(p.shipToCountry, "US"),
		TargetCurrency: defaultValue(p.targetCurrency, "USD"),
		TargetLanguage: defaultValue(p.targetLanguage, "EN"),
		AccessToken:    accessToken,
	}

	// SKU ID 목록이 없거나 20개 이하면 한번에 호출
	if len(allSKUIDs) <= skuBatchSize {
		baseInput.SKUIDs = allSKUIDs
		return p.client.GetAffiliateProductSKUDetail(ctx, baseInput)
	}

	// 20개 초과면 배치 호출 후 합치기
	var combined *clientaliexpress.ProductSKUDetailResult
	for i := 0; i < len(allSKUIDs); i += skuBatchSize {
		end := i + skuBatchSize
		if end > len(allSKUIDs) {
			end = len(allSKUIDs)
		}

		batchInput := baseInput
		batchInput.SKUIDs = allSKUIDs[i:end]

		result, err := p.client.GetAffiliateProductSKUDetail(ctx, batchInput)
		if err != nil {
			return nil, err
		}

		if combined == nil {
			combined = result
		} else {
			combined.SKUInfos = append(combined.SKUInfos, result.SKUInfos...)
		}
	}
	return combined, nil
}

func (p *Provider) resolveAccessToken(ctx context.Context) (string, error) {
	if p.tokenProvider == nil {
		return "", nil
	}
	token, err := p.tokenProvider.GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get aliexpress access token: %w", err)
	}
	return token, nil
}

func buildSKUs(result *clientaliexpress.ProductSKUDetailResult) []domainproduct.NewSKU {
	if result == nil || len(result.SKUInfos) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	skus := make([]domainproduct.NewSKU, 0, len(result.SKUInfos))
	for _, info := range result.SKUInfos {
		originSKUID := fmt.Sprintf("%d", info.SKUID)
		skuProps := strings.TrimSpace(info.SKUProperties)
		externalID := domainproduct.GenerateExternalSKUID(originSKUID, skuProps)
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

		skus = append(skus, domainproduct.NewSKU{
			ExternalSKUID: externalID,
			OriginSKUID:   originSKUID,
			SKUName:       skuName,
			Color:         strings.TrimSpace(info.Color),
			Size:          strings.TrimSpace(info.Size),
			Price:         firstNonEmpty(info.SalePriceWithTax, info.PriceWithTax),
			OriginalPrice: strings.TrimSpace(info.PriceWithTax),
			Currency:      strings.TrimSpace(info.Currency),
			ImageURL:      strings.TrimSpace(info.SKUImageLink),
			SKUProperties: skuProps,
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
