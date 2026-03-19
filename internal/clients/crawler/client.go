package crawler

import "context"

type CrawlInput struct {
	URL string
}

type SKU struct {
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

type Product struct {
	Title     string
	URL       string
	Source    string
	MainImage string
	Images    []string
	SKUs      []SKU
}

type Client interface {
	Crawl(ctx context.Context, input CrawlInput) (*Product, error)
}
