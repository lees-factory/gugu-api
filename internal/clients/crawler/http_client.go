package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

type Config struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewHTTPClient(cfg Config) *HTTPClient {
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:8000"
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}

	return &HTTPClient{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

type crawlRequest struct {
	URL string `json:"url"`
}

type crawlResponse struct {
	Success bool            `json:"success"`
	Data    *crawlProduct   `json:"data"`
	Error   *string         `json:"error"`
}

type crawlProduct struct {
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Source    string     `json:"source"`
	MainImage *string    `json:"main_image"`
	Images    []string   `json:"images"`
	SKUs      []crawlSKU `json:"skus"`
}

type crawlSKU struct {
	ExternalSKUID string  `json:"external_sku_id"`
	SKUName       string  `json:"sku_name"`
	Color         *string `json:"color"`
	Size          *string `json:"size"`
	Price         string  `json:"price"`
	OriginalPrice *string `json:"original_price"`
	Currency      string  `json:"currency"`
	ImageURL      *string `json:"image_url"`
	SKUProperties string  `json:"sku_properties"`
}

func (c *HTTPClient) Crawl(ctx context.Context, input CrawlInput) (*Product, error) {
	body, err := json.Marshal(crawlRequest{URL: input.URL})
	if err != nil {
		return nil, fmt.Errorf("marshal crawl request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/crawl", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("build crawl request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request crawler: %w", err)
	}
	defer resp.Body.Close()

	var result crawlResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode crawl response: %w", err)
	}

	if !result.Success || result.Data == nil {
		errMsg := "unknown crawl error"
		if result.Error != nil {
			errMsg = *result.Error
		}
		return nil, fmt.Errorf("crawl failed: %s", errMsg)
	}

	product := &Product{
		Title:     result.Data.Title,
		URL:       result.Data.URL,
		Source:    result.Data.Source,
		Images:    result.Data.Images,
		SKUs:      make([]SKU, len(result.Data.SKUs)),
	}
	if result.Data.MainImage != nil {
		product.MainImage = *result.Data.MainImage
	}

	for i, s := range result.Data.SKUs {
		product.SKUs[i] = SKU{
			ExternalSKUID: s.ExternalSKUID,
			SKUName:       s.SKUName,
			Color:         derefString(s.Color),
			Size:          derefString(s.Size),
			Price:         s.Price,
			OriginalPrice: derefString(s.OriginalPrice),
			Currency:      s.Currency,
			ImageURL:      derefString(s.ImageURL),
			SKUProperties: s.SKUProperties,
		}
	}

	return product, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
