package skupricehistory

import (
	"context"
	"database/sql"

	domainsph "github.com/ljj/gugu-api/internal/core/domain/skupricehistory"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type SQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{queries: sqldb.New(db)}
}

func (r *SQLCRepository) Create(ctx context.Context, history domainsph.SKUPriceHistory) error {
	return r.queries.CreateSKUPriceHistory(ctx, sqldb.CreateSKUPriceHistoryParams{
		SkuID:       history.SKUID,
		Price:       history.Price,
		Currency:    history.Currency,
		RecordedAt:  history.RecordedAt,
		ChangeValue: history.ChangeValue,
	})
}

func (r *SQLCRepository) ListBySKUID(ctx context.Context, skuID string) ([]domainsph.SKUPriceHistory, error) {
	rows, err := r.queries.ListSKUPriceHistoriesBySKUID(ctx, skuID)
	if err != nil {
		return nil, err
	}

	items := make([]domainsph.SKUPriceHistory, 0, len(rows))
	for _, row := range rows {
		items = append(items, domainsph.SKUPriceHistory{
			SKUID:       row.SkuID,
			Price:       row.Price,
			Currency:    row.Currency,
			RecordedAt:  row.RecordedAt,
			ChangeValue: row.ChangeValue,
		})
	}
	return items, nil
}
