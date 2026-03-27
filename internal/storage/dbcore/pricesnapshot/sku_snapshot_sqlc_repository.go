package pricesnapshot

import (
	"context"
	"database/sql"
	"time"

	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type SKUSnapshotSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSKUSnapshotSQLCRepository(db *sql.DB) *SKUSnapshotSQLCRepository {
	return &SKUSnapshotSQLCRepository{queries: sqldb.New(db)}
}

func (r *SKUSnapshotSQLCRepository) Upsert(ctx context.Context, snapshot domainps.SKUPriceSnapshot) error {
	return r.queries.UpsertSKUPriceSnapshot(ctx, sqldb.UpsertSKUPriceSnapshotParams{
		SkuID:         snapshot.SKUID,
		SnapshotDate:  snapshot.SnapshotDate,
		Price:         snapshot.Price,
		OriginalPrice: snapshot.OriginalPrice,
		Currency:      snapshot.Currency,
	})
}

func (r *SKUSnapshotSQLCRepository) ListBySKUIDAndDateRange(ctx context.Context, skuID string, from time.Time, to time.Time) ([]domainps.SKUPriceSnapshot, error) {
	rows, err := r.queries.ListSKUPriceSnapshotsByDateRange(ctx, sqldb.ListSKUPriceSnapshotsByDateRangeParams{
		SkuID:          skuID,
		SnapshotDate:   from,
		SnapshotDate_2: to,
	})
	if err != nil {
		return nil, err
	}

	items := make([]domainps.SKUPriceSnapshot, 0, len(rows))
	for _, row := range rows {
		items = append(items, domainps.SKUPriceSnapshot{
			SKUID:         row.SkuID,
			SnapshotDate:  row.SnapshotDate,
			Price:         row.Price,
			OriginalPrice: row.OriginalPrice,
			Currency:      row.Currency,
		})
	}
	return items, nil
}
