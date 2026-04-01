package pricealert

import (
	"database/sql"

	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	dbcorepricealert "github.com/ljj/gugu-api/internal/storage/dbcore/pricealert"
	memorypricealert "github.com/ljj/gugu-api/internal/storage/memory/pricealert"
	"github.com/ljj/gugu-api/internal/support/id"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

func Wire(db *sql.DB) *Controller {
	repository := buildRepository(db)

	service := domainpricealert.NewService(
		domainpricealert.NewFinder(repository),
		domainpricealert.NewWriter(repository),
		id.NewRandomHexGenerator(16),
		timeutil.SystemClock{},
	)

	return NewController(service)
}

func buildRepository(db *sql.DB) domainpricealert.Repository {
	if db == nil {
		return memorypricealert.NewRepository()
	}
	return dbcorepricealert.NewSQLCRepository(db)
}
