package product

import (
	"database/sql"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	dbcorepricehistory "github.com/ljj/gugu-api/internal/storage/dbcore/pricehistory"
	memorypricehistory "github.com/ljj/gugu-api/internal/storage/memory/pricehistory"
)

func Wire(db *sql.DB, productService *domainproduct.Service, trackedItemService *domaintrackeditem.Service) *Controller {
	priceHistoryRepository := buildPriceHistoryRepository(db)
	priceHistoryService := domainpricehistory.NewService(
		domainpricehistory.NewFinder(priceHistoryRepository),
	)

	return NewController(productService, priceHistoryService, trackedItemService)
}

func buildPriceHistoryRepository(db *sql.DB) domainpricehistory.Repository {
	if db == nil {
		return memorypricehistory.NewRepository()
	}
	return dbcorepricehistory.NewSQLCRepository(db)
}
