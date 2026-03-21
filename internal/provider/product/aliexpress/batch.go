package aliexpress

import (
	"context"
	"fmt"
	"strconv"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

const batchSize = 20

type PriceResult struct {
	ExternalProductID string
	Price             string
	Currency          string
}

type BatchFetcher struct {
	client        ProductDetailClient
	tokenProvider TokenProvider
}

func NewBatchFetcher(client ProductDetailClient, tokenProvider TokenProvider) *BatchFetcher {
	return &BatchFetcher{client: client, tokenProvider: tokenProvider}
}

func (f *BatchFetcher) FetchPrices(ctx context.Context, externalProductIDs []string) ([]PriceResult, error) {
	accessToken, err := f.resolveAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var results []PriceResult
	for i := 0; i < len(externalProductIDs); i += batchSize {
		end := i + batchSize
		if end > len(externalProductIDs) {
			end = len(externalProductIDs)
		}
		chunk := externalProductIDs[i:end]

		detail, err := f.client.GetAffiliateProductDetail(ctx, clientaliexpress.ProductDetailInput{
			ProductIDs:     chunk,
			TargetCurrency: "KRW",
			TargetLanguage: "KO",
			AccessToken:    accessToken,
		})
		if err != nil {
			return nil, fmt.Errorf("fetch product detail batch: %w", err)
		}
		if detail == nil {
			continue
		}

		for _, p := range detail.Products {
			price := firstNonEmpty(p.TargetSalePrice, p.SalePrice, p.TargetAppSalePrice, p.AppSalePrice)
			currency := firstNonEmpty(p.TargetSalePriceCurrency, p.SalePriceCurrency, p.TargetAppSalePriceCurrency, p.AppSalePriceCurrency)
			results = append(results, PriceResult{
				ExternalProductID: strconv.FormatInt(p.ProductID, 10),
				Price:             price,
				Currency:          currency,
			})
		}
	}

	return results, nil
}

func (f *BatchFetcher) resolveAccessToken(ctx context.Context) (string, error) {
	if f.tokenProvider == nil {
		return "", nil
	}
	token, err := f.tokenProvider.GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return token, nil
}
