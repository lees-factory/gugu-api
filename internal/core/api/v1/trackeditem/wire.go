package trackeditem

import (
	"database/sql"
	"fmt"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	dbcoreproduct "github.com/ljj/gugu-api/internal/storage/dbcore/product"
	dbcoretrackeditem "github.com/ljj/gugu-api/internal/storage/dbcore/trackeditem"
	memoryproduct "github.com/ljj/gugu-api/internal/storage/memory/product"
	memorytrackeditem "github.com/ljj/gugu-api/internal/storage/memory/trackeditem"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

func Wire(cfg config.Config, db *sql.DB) (*Controller, *domaintrackeditem.Service, *domainproduct.Service, error) {
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

	productService := domainproduct.NewService(
		domainproduct.NewFinder(productRepository),
		domainproduct.NewWriter(productRepository),
		id.NewRandomHexGenerator(16),
		clock,
	)
	productCollector := domainproduct.NewDefaultCollector(
		domainproduct.NewAliExpressProductFinder(aliExpressClient, "KRW", "KO", "KR"),
		nil,
	)

	trackedItemService := domaintrackeditem.NewService(
		domaintrackeditem.NewFinder(trackedItemRepository),
		domaintrackeditem.NewWriter(trackedItemRepository),
		id.NewRandomHexGenerator(16),
		clock,
	)

	controller := NewController(trackedItemService, productService, productCollector)
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
