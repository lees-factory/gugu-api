package product

import (
	"database/sql"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	dbcorepricehistory "github.com/ljj/gugu-api/internal/storage/dbcore/pricehistory"
	dbcoresnapshot "github.com/ljj/gugu-api/internal/storage/dbcore/pricesnapshot"
	memorypricehistory "github.com/ljj/gugu-api/internal/storage/memory/pricehistory"
	memorysnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
)

func Wire(db *sql.DB, productService *domainproduct.Service, trackedItemService *domaintrackeditem.Service) *Controller {
	priceHistoryRepository := buildPriceHistoryRepository(db)
	priceHistoryService := domainpricehistory.NewService(
		domainpricehistory.NewFinder(priceHistoryRepository),
	)

	snapshotService := buildSnapshotService(db)

	return NewController(productService, priceHistoryService, snapshotService, trackedItemService)
}

func buildPriceHistoryRepository(db *sql.DB) domainpricehistory.Repository {
	if db == nil {
		return memorypricehistory.NewRepository()
	}
	return dbcorepricehistory.NewSQLCRepository(db)
}

func buildSnapshotService(db *sql.DB) *domainps.Service {
	if db == nil {
		productRepo := memorysnapshot.NewProductSnapshotRepository()
		skuRepo := memorysnapshot.NewSKUSnapshotRepository()
		return domainps.NewService(
			domainps.NewProductSnapshotFinder(productRepo),
			domainps.NewSKUSnapshotFinder(skuRepo),
		)
	}
	productRepo := dbcoresnapshot.NewProductSnapshotSQLCRepository(db)
	skuRepo := dbcoresnapshot.NewSKUSnapshotSQLCRepository(db)
	return domainps.NewService(
		domainps.NewProductSnapshotFinder(productRepo),
		domainps.NewSKUSnapshotFinder(skuRepo),
	)
}
