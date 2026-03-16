package api

import (
	"database/sql"
	"fmt"

	"github.com/go-chi/chi/v5"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiproduct "github.com/ljj/gugu-api/internal/core/api/v1/product"
	apitrackeditem "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	productdetail "github.com/ljj/gugu-api/internal/core/domain/productdetail"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	trackeditemlist "github.com/ljj/gugu-api/internal/core/domain/trackeditemlist"
	dbcorepricehistory "github.com/ljj/gugu-api/internal/storage/dbcore/pricehistory"
	dbcoreproduct "github.com/ljj/gugu-api/internal/storage/dbcore/product"
	dbcoretrackeditem "github.com/ljj/gugu-api/internal/storage/dbcore/trackeditem"
	memorypricehistory "github.com/ljj/gugu-api/internal/storage/memory/pricehistory"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
	memorytrackeditem "github.com/ljj/gugu-api/internal/storage/memory/trackeditem"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

func registerProductModule(router chi.Router, cfg config.Config, db *sql.DB) error {
	productRepository := buildProductRepository(db)
	trackedItemRepository := buildTrackedItemRepository(db)
	priceHistoryRepository := buildPriceHistoryRepository(db)

	aliExpressClient, err := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})
	if err != nil {
		return fmt.Errorf("build aliexpress client: %w", err)
	}

	clock := timeutil.SystemClock{}
	productService := domainproduct.NewService(productRepository, id.NewRandomHexGenerator(16), clock)
	productCollector := domainproduct.NewDefaultCollector(
		domainproduct.NewAliExpressProductFinder(aliExpressClient, "KRW", "KO", "KR"),
		nil,
	)
	trackedItemService := domaintrackeditem.NewService(
		trackedItemRepository,
		id.NewRandomHexGenerator(16),
		clock,
		productService,
		productCollector,
	)
	trackedItemListService := trackeditemlist.NewService(trackedItemRepository, productRepository)
	productDetailService := productdetail.NewService(productRepository, priceHistoryRepository, trackedItemRepository)

	trackedItemController := apitrackeditem.NewController(trackedItemService, trackedItemListService)
	productController := apiproduct.NewController(productDetailService)

	router.Route("/v1/tracked-items", func(r chi.Router) {
		r.Get("/", apiadvice.Wrap(trackedItemController.List))
		r.Post("/", apiadvice.Wrap(trackedItemController.Add))
		r.Delete("/{trackedItemID}", apiadvice.Wrap(trackedItemController.Delete))
	})
	router.Route("/v1/products", func(r chi.Router) {
		r.Get("/{productID}", apiadvice.Wrap(productController.GetDetail))
	})

	return nil
}

func buildProductRepository(db *sql.DB) domainproduct.Repository {
	if db == nil {
		return memoryproduct.NewRepository()
	}
	return dbcoreproduct.NewSQLCRepository(db)
}

func buildTrackedItemRepository(db *sql.DB) domaintrackeditem.Repository {
	if db == nil {
		return memorytrackeditem.NewRepository()
	}
	return dbcoretrackeditem.NewSQLCRepository(db)
}

func buildPriceHistoryRepository(db *sql.DB) domainpricehistory.Repository {
	if db == nil {
		return memorypricehistory.NewRepository()
	}
	return dbcorepricehistory.NewSQLCRepository(db)
}
