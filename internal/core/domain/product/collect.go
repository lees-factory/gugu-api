package product

import (
	"context"
	"fmt"

	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type CollectInput struct {
	Market            Market
	ExternalProductID string
	OriginalURL       string
}

type CollectedProduct struct {
	Market            Market
	ExternalProductID string
	OriginalURL       string
	Title             string
	MainImageURL      string
	CurrentPrice      string
	Currency          string
	ProductURL        string
	CollectionSource  string
	SKUs              []CollectedSKU
}

type CollectedSKU struct {
	ExternalSKUID string
	SKUName       string
	Color         string
	Size          string
	Price         string
	OriginalPrice string
	Currency      string
	ImageURL      string
	SKUProperties string
}

type AffiliateProductFinder interface {
	Find(ctx context.Context, input CollectInput) (*CollectedProduct, error)
}

type CrawlerProductFinder interface {
	Find(ctx context.Context, input CollectInput) (*CollectedProduct, error)
}

type Collector interface {
	Collect(ctx context.Context, input CollectInput) (*CollectedProduct, error)
}

type DefaultCollector struct {
	affiliateProductFinder AffiliateProductFinder
	crawlerProductFinder   CrawlerProductFinder
}

func NewDefaultCollector(affiliateProductFinder AffiliateProductFinder, crawlerProductFinder CrawlerProductFinder) Collector {
	return &DefaultCollector{
		affiliateProductFinder: affiliateProductFinder,
		crawlerProductFinder:   crawlerProductFinder,
	}
}

func (c *DefaultCollector) Collect(ctx context.Context, input CollectInput) (*CollectedProduct, error) {
	if c.affiliateProductFinder != nil {
		item, err := c.affiliateProductFinder.Find(ctx, input)
		if err == nil && item != nil {
			return item, nil
		}
	}

	if c.crawlerProductFinder != nil {
		item, err := c.crawlerProductFinder.Find(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("find crawler product: %w", err)
		}
		if item != nil {
			return item, nil
		}
	}

	return nil, coreerror.New(coreerror.ProductNotFound)
}
