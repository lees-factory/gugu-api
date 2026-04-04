package aliexpress

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

const (
	collectionSourceAffiliate = "AFFILIATE_API"
	collectionSourceDS        = "DS_API"
)

type AffiliateClient interface {
	GetAffiliateProductDetail(ctx context.Context, input clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error)
}

type DSClient interface {
	GetDSProduct(ctx context.Context, input clientaliexpress.DSProductInput) (*clientaliexpress.DSProductResult, error)
}

type TokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type Provider struct {
	affiliateClient        AffiliateClient
	dsClient               DSClient
	affiliateTokenProvider TokenProvider
	dsTokenProvider        TokenProvider
	targetCurrency         string
	targetLanguage         string
	shipToCountry          string
}

func NewProvider(
	affiliateClient AffiliateClient,
	dsClient DSClient,
	affiliateTokenProvider TokenProvider,
	dsTokenProvider TokenProvider,
	targetCurrency string,
	targetLanguage string,
	shipToCountry string,
) *Provider {
	return &Provider{
		affiliateClient:        affiliateClient,
		dsClient:               dsClient,
		affiliateTokenProvider: affiliateTokenProvider,
		dsTokenProvider:        dsTokenProvider,
		targetCurrency:         strings.TrimSpace(targetCurrency),
		targetLanguage:         strings.TrimSpace(targetLanguage),
		shipToCountry:          strings.TrimSpace(shipToCountry),
	}
}

func (p *Provider) Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string, currency string, language string) (*domainproduct.NewProduct, error) {
	if market != enum.MarketAliExpress {
		return nil, nil
	}

	options := p.resolveRequestOptions(currency, language)

	// 1. Affiliate API로 상품 조회
	affiliateProduct, promotionLink := p.fetchAffiliate(ctx, externalProductID, options)

	// 2. DS API로 SKU 전체 조회 (항상)
	dsResult := p.fetchDS(ctx, externalProductID, options)

	// 3. 조합
	if affiliateProduct != nil {
		// Case 1: Affiliate 상품 + DS SKU
		return p.buildFromAffiliate(affiliateProduct, dsResult, promotionLink, externalProductID, originalURL, options.targetLanguage), nil
	}

	if dsResult != nil {
		// Case 2: Affiliate 없음 → DS fallback
		slog.Info("affiliate product not found, using DS fallback", "product_id", externalProductID)
		return p.buildFromDS(dsResult, externalProductID, originalURL, options.targetLanguage), nil
	}

	// Case 3: 둘 다 없음
	return nil, nil
}

func (p *Provider) fetchAffiliate(ctx context.Context, externalProductID string, options requestOptions) (*clientaliexpress.AffiliateProduct, string) {
	if p.affiliateClient == nil {
		return nil, ""
	}

	accessToken, err := p.resolveToken(ctx, p.affiliateTokenProvider)
	if err != nil {
		slog.Warn("failed to get affiliate access token", "error", err)
		return nil, ""
	}

	result, err := p.affiliateClient.GetAffiliateProductDetail(ctx, clientaliexpress.ProductDetailInput{
		ProductIDs:     []string{externalProductID},
		TargetCurrency: options.targetCurrency,
		TargetLanguage: options.targetLanguage,
		AccessToken:    accessToken,
	})
	if err != nil {
		slog.Warn("failed to fetch affiliate product detail", "product_id", externalProductID, "error", err)
		return nil, ""
	}
	if result == nil || len(result.Products) == 0 {
		return nil, ""
	}

	product := result.Products[0]
	return &product, product.PromotionLink
}

func (p *Provider) fetchDS(ctx context.Context, externalProductID string, options requestOptions) *clientaliexpress.DSProductResult {
	if p.dsClient == nil {
		return nil
	}

	accessToken, err := p.resolveToken(ctx, p.dsTokenProvider)
	if err != nil {
		slog.Warn("failed to get ds access token", "error", err)
		return nil
	}

	result, err := p.dsClient.GetDSProduct(ctx, clientaliexpress.DSProductInput{
		ProductID:      externalProductID,
		ShipToCountry:  options.shipToCountry,
		TargetCurrency: options.targetCurrency,
		TargetLanguage: strings.ToLower(options.targetLanguage),
		AccessToken:    accessToken,
	})
	if err != nil {
		slog.Warn("failed to fetch ds product", "product_id", externalProductID, "error", err)
		return nil
	}

	return result
}

func (p *Provider) buildFromAffiliate(ap *clientaliexpress.AffiliateProduct, ds *clientaliexpress.DSProductResult, promotionLink string, externalProductID string, originalURL string, language string) *domainproduct.NewProduct {
	currency := firstNonEmpty(ap.TargetSalePriceCurrency, ap.SalePriceCurrency, ap.TargetAppSalePriceCurrency, ap.AppSalePriceCurrency)

	product := &domainproduct.NewProduct{
		Market:            enum.MarketAliExpress,
		ExternalProductID: externalProductID,
		OriginalURL:       firstNonEmpty(originalURL, ap.ProductDetailURL),
		Language:          strings.ToUpper(strings.TrimSpace(language)),
		Title:             ap.ProductTitle,
		MainImageURL:      ap.ProductMainImageURL,
		CurrentPrice:      firstNonEmpty(ap.TargetSalePrice, ap.SalePrice, ap.TargetAppSalePrice, ap.AppSalePrice),
		Currency:          currency,
		ProductURL:        ap.ProductDetailURL,
		PromotionLink:     promotionLink,
		CollectionSource:  collectionSourceAffiliate,
	}

	// SKU는 DS API에서 가져옴 (전체 SKU)
	if ds != nil {
		product.SKUs = p.mapDSSKUs(ds)
	}

	return product
}

func (p *Provider) buildFromDS(ds *clientaliexpress.DSProductResult, externalProductID string, originalURL string, language string) *domainproduct.NewProduct {
	baseInfo := ds.BaseInfo

	var currency string
	if len(ds.SKUs) > 0 {
		currency = ds.SKUs[0].CurrencyCode
	}
	if currency == "" {
		currency = baseInfo.CurrencyCode
	}

	var imageURL string
	if ds.Multimedia.ImageURLs != "" {
		parts := strings.SplitN(ds.Multimedia.ImageURLs, ";", 2)
		imageURL = strings.TrimSpace(parts[0])
	}

	product := &domainproduct.NewProduct{
		Market:            enum.MarketAliExpress,
		ExternalProductID: externalProductID,
		OriginalURL:       originalURL,
		Language:          strings.ToUpper(strings.TrimSpace(language)),
		Title:             baseInfo.Subject,
		MainImageURL:      imageURL,
		CurrentPrice:      firstDSCurrentPrice(ds),
		Currency:          currency,
		ProductURL:        originalURL,
		CollectionSource:  collectionSourceDS,
		SKUs:              p.mapDSSKUs(ds),
	}

	return product
}

func (p *Provider) mapDSSKUs(ds *clientaliexpress.DSProductResult) []domainproduct.NewSKU {
	if ds == nil || len(ds.SKUs) == 0 {
		return nil
	}

	skus := make([]domainproduct.NewSKU, 0, len(ds.SKUs))
	for _, sku := range ds.SKUs {
		color, size, propStr := extractSKUProperties(sku.Properties)

		skus = append(skus, domainproduct.NewSKU{
			ExternalSKUID: sku.SKUID,
			OriginSKUID:   sku.ID,
			Color:         color,
			Size:          size,
			Price:         firstNonEmpty(sku.OfferSalePrice, sku.SKUPrice),
			OriginalPrice: sku.SKUPrice,
			Currency:      sku.CurrencyCode,
			ImageURL:      extractSKUImage(sku.Properties),
			SKUProperties: propStr,
		})
	}
	return skus
}

const (
	aePropertyIDColor = 14
	aePropertyIDSize  = 5
)

func extractSKUProperties(props []clientaliexpress.DSSKUPropertyDTO) (color, size, propJSON string) {
	var parts []string
	for _, prop := range props {
		value := firstNonEmpty(prop.PropertyValueDefinitionName, prop.SKUPropertyValue)

		switch prop.SKUPropertyID {
		case aePropertyIDColor:
			color = value
		case aePropertyIDSize:
			size = value
		}
		parts = append(parts, prop.SKUPropertyName+":"+value)
	}
	return color, size, strings.Join(parts, ";")
}

func extractSKUImage(props []clientaliexpress.DSSKUPropertyDTO) string {
	for _, prop := range props {
		if prop.SKUImage != "" {
			return prop.SKUImage
		}
	}
	return ""
}

func (p *Provider) resolveToken(ctx context.Context, tp TokenProvider) (string, error) {
	if tp == nil {
		return "", nil
	}
	token, err := tp.GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return token, nil
}

func defaultValue(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

type requestOptions struct {
	targetCurrency string
	targetLanguage string
	shipToCountry  string
}

func (p *Provider) resolveRequestOptions(currency string, language string) requestOptions {
	targetCurrency := defaultValue(currency, p.targetCurrency)
	targetCurrency = defaultValue(targetCurrency, "KRW")

	targetLanguage := defaultValue(strings.ToUpper(strings.TrimSpace(language)), languageForCurrency(targetCurrency))
	targetLanguage = defaultValue(targetLanguage, p.targetLanguage)
	targetLanguage = defaultValue(targetLanguage, "KO")

	shipToCountry := defaultValue(shipToCountryForCurrency(targetCurrency), p.shipToCountry)
	shipToCountry = defaultValue(shipToCountry, "KR")

	return requestOptions{
		targetCurrency: targetCurrency,
		targetLanguage: targetLanguage,
		shipToCountry:  shipToCountry,
	}
}

func languageForCurrency(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "KRW":
		return "KO"
	default:
		return "EN"
	}
}

func shipToCountryForCurrency(currency string) string {
	switch strings.ToUpper(strings.TrimSpace(currency)) {
	case "KRW":
		return "KR"
	default:
		return "US"
	}
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

func firstDSCurrentPrice(ds *clientaliexpress.DSProductResult) string {
	if ds == nil {
		return ""
	}
	for _, sku := range ds.SKUs {
		if price := firstNonEmpty(sku.OfferSalePrice, sku.SKUPrice); price != "" {
			return price
		}
	}
	return ""
}
