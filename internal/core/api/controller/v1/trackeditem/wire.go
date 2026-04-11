package trackeditem

import (
	"database/sql"
	"fmt"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
	dbcorepricealert "github.com/ljj/gugu-api/internal/storage/dbcore/pricealert"
	dbcoresnapshot "github.com/ljj/gugu-api/internal/storage/dbcore/pricesnapshot"
	dbcoreproduct "github.com/ljj/gugu-api/internal/storage/dbcore/product"
	dbcoreskuph "github.com/ljj/gugu-api/internal/storage/dbcore/skupricehistory"
	dbcoretrackeditem "github.com/ljj/gugu-api/internal/storage/dbcore/trackeditem"
	memorypricealert "github.com/ljj/gugu-api/internal/storage/memory/pricealert"
	memorysnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
	memoryskuph "github.com/ljj/gugu-api/internal/storage/memory/skupricehistory"
	memorytrackeditem "github.com/ljj/gugu-api/internal/storage/memory/trackeditem"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

func Wire(cfg config.Config, db *sql.DB, aliExpressTokenStore clientaliexpress.TokenStore) (*Controller, *domaintrackeditem.Service, *domainproduct.Service, error) {
	trackedItemRepository := buildTrackedItemRepository(db)
	productRepository := buildProductRepository(db)
	productVariantRepository := buildProductVariantRepository(db)

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
	skuPriceHistRepo := buildSKUPriceHistoryRepository(db)
	priceAlertRepo := buildPriceAlertRepository(db)

	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepository),
		domainproduct.NewWriter(productRepository),
		productVariantRepository,
		skuRepository,
		id.NewRandomHexGenerator(16),
		clock,
		nil,
		domainsph.NewWriter(skuPriceHistRepo),
	)

	affiliateTokenProvider := provideraliexpress.NewTokenProvider("AFFILIATE", aliExpressTokenStore, aliExpressClient)

	var dsClient *clientaliexpress.HTTPClient
	var dsTokenProvider provideraliexpress.TokenProvider
	if cfg.AliExpressDSAppKey != "" {
		dsClient, _ = clientaliexpress.NewHTTPClient(clientaliexpress.Config{
			BaseURL:     cfg.AliExpressDSBaseURL,
			AppKey:      cfg.AliExpressDSAppKey,
			AppSecret:   cfg.AliExpressDSAppSecret,
			CallbackURL: cfg.AliExpressDSCallbackURL,
		})
		dsTokenProvider = provideraliexpress.NewTokenProvider("DROPSHIPPING", aliExpressTokenStore, dsClient)
	}

	provider := provideraliexpress.NewProvider(
		aliExpressClient,
		dsClient,
		affiliateTokenProvider,
		dsTokenProvider,
		"KRW", "KO", "KR",
	)

	trackedItemService := domaintrackeditem.NewService(
		domaintrackeditem.NewFinder(trackedItemRepository),
		domaintrackeditem.NewWriter(trackedItemRepository),
		id.NewRandomHexGenerator(16),
		clock,
		productService,
		provider,
	)

	skuPriceHistoryService := domainsph.NewService(domainsph.NewFinder(skuPriceHistRepo))
	snapshotService := buildSnapshotService(db)
	priceAlertService := domainpricealert.NewService(
		domainpricealert.NewFinder(priceAlertRepo),
		domainpricealert.NewWriter(priceAlertRepo),
		id.NewRandomHexGenerator(16),
		clock,
	)

	controller := NewController(trackedItemService, skuPriceHistoryService, snapshotService, priceAlertService)
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

func buildProductVariantRepository(db *sql.DB) domainproduct.VariantRepository {
	if db == nil {
		return memoryproduct.NewVariantRepository()
	}
	return dbcoreproduct.NewLocalizationRepository(db)
}

func buildSKUPriceHistoryRepository(db *sql.DB) domainsph.Repository {
	if db == nil {
		return memoryskuph.NewRepository()
	}
	return dbcoreskuph.NewSQLCRepository(db)
}

func buildPriceAlertRepository(db *sql.DB) domainpricealert.Repository {
	if db == nil {
		return memorypricealert.NewRepository()
	}
	return dbcorepricealert.NewSQLCRepository(db)
}

func buildSnapshotService(db *sql.DB) *domainps.Service {
	if db == nil {
		skuRepo := memorysnapshot.NewSKUSnapshotRepository()
		return domainps.NewService(domainps.NewSKUSnapshotFinder(skuRepo))
	}
	skuRepo := dbcoresnapshot.NewSKUSnapshotSQLCRepository(db)
	return domainps.NewService(domainps.NewSKUSnapshotFinder(skuRepo))
}
