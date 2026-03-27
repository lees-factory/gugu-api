package batch

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
)

const hotProductCollectionSource = "HOT_PRODUCT_QUERY"

type HotProductLoadInput struct {
	CategoryIDs    string
	Keywords       string
	PageNo         string
	PageSize       string
	Sort           string
	MinSalePrice   string
	MaxSalePrice   string
	ShipToCountry  string
	TargetCurrency string
	TargetLanguage string
}

type HotProductLoadResult struct {
	TotalFetched   int `json:"total_fetched"`
	NewlyCreated   int `json:"newly_created"`
	AlreadyExisted int `json:"already_existed"`
}

type HotProductClient interface {
	GetAffiliateHotProducts(ctx context.Context, input clientaliexpress.ProductQueryInput) (*clientaliexpress.ProductQueryResult, error)
}

type HotProductLoader struct {
	client         HotProductClient
	tokenProvider  provideraliexpress.TokenProvider
	productService *domainproduct.Service
}

func NewHotProductLoader(
	client HotProductClient,
	tokenProvider provideraliexpress.TokenProvider,
	productService *domainproduct.Service,
) *HotProductLoader {
	return &HotProductLoader{
		client:         client,
		tokenProvider:  tokenProvider,
		productService: productService,
	}
}

func (l *HotProductLoader) LoadHotProducts(ctx context.Context, input HotProductLoadInput) (*HotProductLoadResult, error) {
	accessToken, err := l.resolveAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	result, err := l.client.GetAffiliateHotProducts(ctx, clientaliexpress.ProductQueryInput{
		CategoryIDs:    input.CategoryIDs,
		Keywords:       input.Keywords,
		PageNo:         input.PageNo,
		PageSize:       input.PageSize,
		Sort:           input.Sort,
		MinSalePrice:   input.MinSalePrice,
		MaxSalePrice:   input.MaxSalePrice,
		ShipToCountry:  input.ShipToCountry,
		TargetCurrency: defaultString(input.TargetCurrency, "KRW"),
		TargetLanguage: defaultString(input.TargetLanguage, "KO"),
		AccessToken:    accessToken,
	})
	if err != nil {
		return nil, fmt.Errorf("fetch hot products: %w", err)
	}

	loadResult := &HotProductLoadResult{
		TotalFetched: len(result.Products),
	}

	for _, p := range result.Products {
		externalID := strconv.FormatInt(p.ProductID, 10)

		existing, err := l.productService.FindByMarketAndExternalProductID(ctx, enum.MarketAliExpress, externalID)
		if err != nil {
			log.Printf("failed to check existing product %s: %v", externalID, err)
			continue
		}
		if existing != nil {
			loadResult.AlreadyExisted++
			continue
		}

		price := firstNonEmpty(p.TargetSalePrice, p.SalePrice, p.TargetAppSalePrice, p.AppSalePrice)
		currency := firstNonEmpty(p.TargetSalePriceCurrency, p.SalePriceCurrency, p.TargetAppSalePriceCurrency, p.AppSalePriceCurrency)

		_, err = l.productService.Create(ctx, domainproduct.NewProduct{
			Market:            enum.MarketAliExpress,
			ExternalProductID: externalID,
			OriginalURL:       strings.TrimSpace(p.ProductDetailURL),
			Title:             strings.TrimSpace(p.ProductTitle),
			MainImageURL:      strings.TrimSpace(p.ProductMainImageURL),
			CurrentPrice:      price,
			Currency:          currency,
			ProductURL:        strings.TrimSpace(p.ProductDetailURL),
			CollectionSource:  hotProductCollectionSource,
		})
		if err != nil {
			log.Printf("failed to create hot product %s: %v", externalID, err)
			continue
		}
		loadResult.NewlyCreated++
	}

	log.Printf("hot product load: fetched=%d created=%d existed=%d",
		loadResult.TotalFetched, loadResult.NewlyCreated, loadResult.AlreadyExisted)
	return loadResult, nil
}

func (l *HotProductLoader) resolveAccessToken(ctx context.Context) (string, error) {
	if l.tokenProvider == nil {
		return "", nil
	}
	token, err := l.tokenProvider.GetAccessToken(ctx)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return token, nil
}

func defaultString(value string, fallback string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return fallback
	}
	return v
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v != "" {
			return v
		}
	}
	return ""
}
