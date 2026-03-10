package crawler

import "context"

type CrawlInput struct {
	URL string
}

type SKU struct {
	Name          string
	Price         string
	OriginalPrice string
	ImageURL      string
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
