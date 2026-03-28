package aliexpress

import (
	"context"
	"fmt"
	"strings"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

const collectionSource = "AFFILIATE_API"

type ProductDetailClient interface {
	GetAffiliateProductDetail(ctx context.Context, input clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error)
}

type TokenProvider interface {
	GetAccessToken(ctx context.Context) (string, error)
}

type Provider struct {
	client         ProductDetailClient
	tokenProvider  TokenProvider
	targetCurrency string
	targetLanguage string
	shipToCountry  string
}

func NewProvider(client ProductDetailClient, tokenProvider TokenProvider, targetCurrency string, targetLanguage string, shipToCountry string) *Provider {
	return &Provider{
		client:         client,
		tokenProvider:  tokenProvider,
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

	return &domainproduct.NewProduct{
		Market:            enum.MarketAliExpress,
		ExternalProductID: externalProductID,
		OriginalURL:       firstNonEmpty(originalURL, detailProduct.ProductDetailURL),
		Title:             detailProduct.ProductTitle,
		MainImageURL:      detailProduct.ProductMainImageURL,
		ProductURL:        detailProduct.ProductDetailURL,
		CollectionSource:  collectionSource,
	}, nil
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
