package trackeditem

import (
	"context"
	"database/sql"
	"errors"

	domaintrackeditem "github.com/ljj/gugu-api/internal/core/domain/trackeditem"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

type TrackedItemSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *TrackedItemSQLCRepository {
	return &TrackedItemSQLCRepository{queries: sqldb.New(db)}
}

func (r *TrackedItemSQLCRepository) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*domaintrackeditem.TrackedItem, error) {
	row, err := r.queries.FindTrackedItemByUserIDAndProductID(ctx, sqldb.FindTrackedItemByUserIDAndProductIDParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainTrackedItem(row)
	return &item, nil
}

func (r *TrackedItemSQLCRepository) FindByIDAndUserID(ctx context.Context, trackedItemID string, userID string) (*domaintrackeditem.TrackedItem, error) {
	row, err := r.queries.FindTrackedItemByIDAndUserID(ctx, sqldb.FindTrackedItemByIDAndUserIDParams{
		ID:     trackedItemID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainTrackedItem(row)
	return &item, nil
}

func (r *TrackedItemSQLCRepository) ListByUserID(ctx context.Context, userID string) ([]domaintrackeditem.TrackedItem, error) {
	rows, err := r.queries.ListTrackedItemsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	items := make([]domaintrackeditem.TrackedItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainTrackedItem(row))
	}
	return items, nil
}

func (r *TrackedItemSQLCRepository) Create(ctx context.Context, trackedItem domaintrackeditem.TrackedItem) error {
	return r.queries.CreateTrackedItem(ctx, sqldb.CreateTrackedItemParams{
		ID:          trackedItem.ID,
		UserID:      trackedItem.UserID,
		ProductID:   trackedItem.ProductID,
		SkuID:       toNullString(trackedItem.SKUID),
		OriginalUrl: trackedItem.OriginalURL,
		CreatedAt:   trackedItem.CreatedAt,
	})
}

func (r *TrackedItemSQLCRepository) UpdateSKU(ctx context.Context, trackedItemID string, userID string, skuID string) error {
	affected, err := r.queries.UpdateTrackedItemSKU(ctx, sqldb.UpdateTrackedItemSKUParams{
		ID:     trackedItemID,
		UserID: userID,
		SkuID:  toNullString(skuID),
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TrackedItemSQLCRepository) DeleteByIDAndUserID(ctx context.Context, trackedItemID string, userID string) error {
	affected, err := r.queries.DeleteTrackedItemByIDAndUserID(ctx, sqldb.DeleteTrackedItemByIDAndUserIDParams{
		ID:     trackedItemID,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func toDomainTrackedItem(row sqldb.GuguUserTrackedItem) domaintrackeditem.TrackedItem {
	return domaintrackeditem.TrackedItem{
		ID:          row.ID,
		UserID:      row.UserID,
		ProductID:   row.ProductID,
		SKUID:       fromNullString(row.SkuID),
		OriginalURL: row.OriginalUrl,
		CreatedAt:   row.CreatedAt,
	}
}
