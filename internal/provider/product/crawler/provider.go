package crawler

import (
	"context"
	"fmt"

	clientcrawler "github.com/ljj/gugu-api/internal/clients/crawler"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
)

const collectionSource = "CRAWLER"

type Provider struct {
	client clientcrawler.Client
}

func NewProvider(client clientcrawler.Client) *Provider {
	return &Provider{client: client}
}

func (p *Provider) Provide(ctx context.Context, market enum.Market, externalProductID string, originalURL string) (*domainproduct.NewProduct, error) {
	if p.client == nil {
		return nil, nil
	}
	if originalURL == "" {
		return nil, nil
	}

	result, err := p.client.Crawl(ctx, clientcrawler.CrawlInput{URL: originalURL})
	if err != nil {
		return nil, fmt.Errorf("crawl product: %w", err)
	}
	if result == nil {
		return nil, nil
	}

	skus := make([]domainproduct.NewSKU, len(result.SKUs))
	for i, s := range result.SKUs {
		skus[i] = domainproduct.NewSKU{
			ExternalSKUID: s.ExternalSKUID,
			SKUName:       s.SKUName,
			Color:         s.Color,
			Size:          s.Size,
			Price:         s.Price,
			OriginalPrice: s.OriginalPrice,
			Currency:      s.Currency,
			ImageURL:      s.ImageURL,
			SKUProperties: s.SKUProperties,
		}
	}

	price := ""
	currency := ""
	if len(result.SKUs) > 0 {
		price = result.SKUs[0].Price
		currency = result.SKUs[0].Currency
	}

	return &domainproduct.NewProduct{
		Market:            market,
		ExternalProductID: externalProductID,
		OriginalURL:       originalURL,
		Title:             result.Title,
		MainImageURL:      result.MainImage,
		CurrentPrice:      price,
		Currency:          currency,
		ProductURL:        result.URL,
		CollectionSource:  collectionSource,
		SKUs:              skus,
	}, nil
}
