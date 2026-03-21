package pricehistory

import (
	"context"
	"database/sql"

	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type PriceHistorySQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *PriceHistorySQLCRepository {
	return &PriceHistorySQLCRepository{queries: sqldb.New(db)}
}

func (r *PriceHistorySQLCRepository) Create(ctx context.Context, history domainpricehistory.PriceHistory) error {
	return r.queries.CreatePriceHistory(ctx, sqldb.CreatePriceHistoryParams{
		ProductID:   history.ProductID,
		Price:       history.Price,
		Currency:    history.Currency,
		RecordedAt:  history.RecordedAt,
		ChangeValue: history.ChangeValue,
	})
}

func (r *PriceHistorySQLCRepository) ListByProductID(ctx context.Context, productID string) ([]domainpricehistory.PriceHistory, error) {
	rows, err := r.queries.ListPriceHistoriesByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	items := make([]domainpricehistory.PriceHistory, 0, len(rows))
	for _, row := range rows {
		items = append(items, domainpricehistory.PriceHistory{
			ProductID:   row.ProductID,
			Price:       row.Price,
			Currency:    row.Currency,
			RecordedAt:  row.RecordedAt,
			ChangeValue: row.ChangeValue,
		})
	}
	return items, nil
}
