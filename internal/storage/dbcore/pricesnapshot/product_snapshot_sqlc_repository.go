package pricesnapshot

import (
	"context"
	"database/sql"
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type ProductSnapshotSQLCRepository struct {
	queries *sqldb.Queries
}

func NewProductSnapshotSQLCRepository(db *sql.DB) *ProductSnapshotSQLCRepository {
	return &ProductSnapshotSQLCRepository{queries: sqldb.New(db)}
}

func (r *ProductSnapshotSQLCRepository) Upsert(ctx context.Context, snapshot domainps.ProductPriceSnapshot) error {
	return r.queries.UpsertProductPriceSnapshot(ctx, sqldb.UpsertProductPriceSnapshotParams{
		ProductID:    snapshot.ProductID,
		SnapshotDate: snapshot.SnapshotDate,
		Price:        snapshot.Price,
		Currency:     snapshot.Currency,
	})
}

func (r *ProductSnapshotSQLCRepository) ListByProductIDAndDateRange(ctx context.Context, productID string, currency string, from time.Time, to time.Time) ([]domainps.ProductPriceSnapshot, error) {
	rows, err := r.queries.ListProductPriceSnapshotsByDateRange(ctx, sqldb.ListProductPriceSnapshotsByDateRangeParams{
		ProductID:      productID,
		Currency:       currency,
		SnapshotDate:   from,
		SnapshotDate_2: to,
	})
	if err != nil {
		return nil, err
	}

	items := make([]domainps.ProductPriceSnapshot, 0, len(rows))
	for _, row := range rows {
		items = append(items, domainps.ProductPriceSnapshot{
			ProductID:    row.ProductID,
			SnapshotDate: row.SnapshotDate,
			Price:        row.Price,
			Currency:     row.Currency,
		})
	}
	return items, nil
}
