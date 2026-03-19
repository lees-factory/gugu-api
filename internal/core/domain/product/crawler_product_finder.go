package product

import (
	"context"
	"fmt"

	clientcrawler "github.com/ljj/gugu-api/internal/clients/crawler"
)

const crawlerCollectionSource = "CRAWLER"

type crawlerProductFinder struct {
	client clientcrawler.Client
}

func NewCrawlerProductFinder(client clientcrawler.Client) CrawlerProductFinder {
	return &crawlerProductFinder{client: client}
}

func (f *crawlerProductFinder) Find(ctx context.Context, input CollectInput) (*CollectedProduct, error) {
	if f.client == nil {
		return nil, nil
	}

	url := input.OriginalURL
	if url == "" {
		return nil, nil
	}

	result, err := f.client.Crawl(ctx, clientcrawler.CrawlInput{URL: url})
	if err != nil {
		return nil, fmt.Errorf("crawl product: %w", err)
	}
	if result == nil {
		return nil, nil
	}

	skus := make([]CollectedSKU, len(result.SKUs))
	for i, s := range result.SKUs {
		skus[i] = CollectedSKU{
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

	return &CollectedProduct{
		Market:            input.Market,
		ExternalProductID: input.ExternalProductID,
		OriginalURL:       input.OriginalURL,
		Title:             result.Title,
		MainImageURL:      result.MainImage,
		CurrentPrice:      price,
		Currency:          currency,
		ProductURL:        result.URL,
		CollectionSource:  crawlerCollectionSource,
		SKUs:              skus,
	}, nil
}
