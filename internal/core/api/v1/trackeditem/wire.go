package trackeditem

import (
	"database/sql"
	"fmt"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	clientcrawler "github.com/ljj/gugu-api/internal/clients/crawler"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	productprovider "github.com/ljj/gugu-api/internal/provider/product"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
	providercrawler "github.com/ljj/gugu-api/internal/provider/product/crawler"
	dbcoreproduct "github.com/ljj/gugu-api/internal/storage/dbcore/product"
	dbcoretrackeditem "github.com/ljj/gugu-api/internal/storage/dbcore/trackeditem"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
	memorytrackeditem "github.com/ljj/gugu-api/internal/storage/memory/trackeditem"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

func Wire(cfg config.Config, db *sql.DB, aliExpressTokenStore clientaliexpress.TokenStore) (*Controller, *domaintrackeditem.Service, *domainproduct.Service, error) {
	trackedItemRepository := buildTrackedItemRepository(db)
	productRepository := buildProductRepository(db)

	aliExpressClient, err := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("build aliexpress client: %w", err)
	}

	clock := timeutil.SystemClock{}
	skuRepository := buildSKURepository(db)

	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepository),
		domainproduct.NewWriter(productRepository),
		skuRepository,
		id.NewRandomHexGenerator(16),
		clock,
	)

	crawlerClient := clientcrawler.NewHTTPClient(clientcrawler.Config{
		BaseURL: cfg.CrawlerBaseURL,
	})

	tokenProvider := provideraliexpress.NewTokenProvider(aliExpressTokenStore)
	provider := productprovider.NewFallbackProvider(
		provideraliexpress.NewProvider(aliExpressClient, tokenProvider, "KRW", "KO", "KR"),
		providercrawler.NewProvider(crawlerClient),
	)

	trackedItemService := domaintrackeditem.NewService(
		domaintrackeditem.NewFinder(trackedItemRepository),
		domaintrackeditem.NewWriter(trackedItemRepository),
		id.NewRandomHexGenerator(16),
		clock,
		productService,
		provider,
	)

	controller := NewController(trackedItemService)
	return controller, trackedItemService, productService, nil
}

func buildTrackedItemRepository(db *sql.DB) domaintrackeditem.Repository {
	if db == nil {
		return memorytrackeditem.NewRepository()
	}
	return dbcoretrackeditem.NewSQLCRepository(db)
}

func buildProductRepository(db *sql.DB) domainproduct.Repository {
	if db == nil {
		return memoryproduct.NewRepository()
	}
	return dbcoreproduct.NewSQLCRepository(db)
}

func buildSKURepository(db *sql.DB) domainproduct.SKURepository {
	if db == nil {
		return memoryproduct.NewSKURepository()
	}
	return dbcoreproduct.NewSKUSQLCRepository(db)
}
