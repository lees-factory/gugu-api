package product

import (
	"database/sql"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	dbcoresnapshot "github.com/ljj/gugu-api/internal/storage/dbcore/pricesnapshot"
	memorysnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
)

func Wire(db *sql.DB, productService *domainproduct.Service, trackedItemService *domaintrackeditem.Service) *Controller {
	snapshotService := buildSnapshotService(db)

	return NewController(productService, snapshotService, trackedItemService)
}

func buildSnapshotService(db *sql.DB) *domainps.Service {
	if db == nil {
		skuRepo := memorysnapshot.NewSKUSnapshotRepository()
		return domainps.NewService(domainps.NewSKUSnapshotFinder(skuRepo))
	}
	skuRepo := dbcoresnapshot.NewSKUSnapshotSQLCRepository(db)
	return domainps.NewService(domainps.NewSKUSnapshotFinder(skuRepo))
}
