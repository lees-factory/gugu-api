package trackeditem

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

func (r *TrackedItemSQLCRepository) ListByUserIDWithCursor(ctx context.Context, userID string, cursorCreatedAt time.Time, cursorID string, limit int) ([]domaintrackeditem.TrackedItem, error) {
	rows, err := r.queries.ListTrackedItemsByUserIDWithCursor(ctx, sqldb.ListTrackedItemsByUserIDWithCursorParams{
		UserID:    userID,
		CreatedAt: cursorCreatedAt,
		ID:        cursorID,
		Limit:     int32(limit),
	})
	if err != nil {
		return nil, err
	}

	items := make([]domaintrackeditem.TrackedItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainTrackedItem(row))
	}
	return items, nil
}

func (r *TrackedItemSQLCRepository) ListByUserIDFirstPage(ctx context.Context, userID string, limit int) ([]domaintrackeditem.TrackedItem, error) {
	rows, err := r.queries.ListTrackedItemsByUserIDFirstPage(ctx, sqldb.ListTrackedItemsByUserIDFirstPageParams{
		UserID: userID,
		Limit:  int32(limit),
	})
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
		ID:                    trackedItem.ID,
		UserID:                trackedItem.UserID,
		ProductID:             trackedItem.ProductID,
		SkuID:                 toNullString(trackedItem.SKUID),
		OriginalUrl:           trackedItem.OriginalURL,
		ViewExternalProductID: trackedItem.ViewExternalProductID,
		PreferredLanguage:     trackedItem.PreferredLanguage,
		TrackingScope:         trackedItem.TrackingScope,
		Currency:              trackedItem.Currency,
		CreatedAt:             trackedItem.CreatedAt,
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

func (r *TrackedItemSQLCRepository) UpdatePreferredLanguage(ctx context.Context, trackedItemID string, userID string, preferredLanguage string) error {
	affected, err := r.queries.UpdateTrackedItemPreferredLanguage(ctx, sqldb.UpdateTrackedItemPreferredLanguageParams{
		ID:                trackedItemID,
		UserID:            userID,
		PreferredLanguage: preferredLanguage,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TrackedItemSQLCRepository) UpdateTrackingScope(ctx context.Context, trackedItemID string, userID string, trackingScope string) error {
	affected, err := r.queries.UpdateTrackedItemTrackingScope(ctx, sqldb.UpdateTrackedItemTrackingScopeParams{
		ID:            trackedItemID,
		UserID:        userID,
		TrackingScope: trackingScope,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *TrackedItemSQLCRepository) ReplaceWatchSKUs(ctx context.Context, trackedItemID string, skuIDs []string) error {
	if err := r.queries.DeleteTrackedItemWatchSKUs(ctx, trackedItemID); err != nil {
		return fmt.Errorf("delete tracked item watch skus: %w", err)
	}
	for _, skuID := range skuIDs {
		if skuID == "" {
			continue
		}
		if err := r.queries.CreateTrackedItemWatchSKU(ctx, sqldb.CreateTrackedItemWatchSKUParams{
			TrackedItemID: trackedItemID,
			SkuID:         skuID,
		}); err != nil {
			return fmt.Errorf("create tracked item watch sku: %w", err)
		}
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
		ID:                    row.ID,
		UserID:                row.UserID,
		ProductID:             row.ProductID,
		SKUID:                 fromNullString(row.SkuID),
		OriginalURL:           row.OriginalUrl,
		ViewExternalProductID: row.ViewExternalProductID,
		PreferredLanguage:     row.PreferredLanguage,
		TrackingScope:         row.TrackingScope,
		Currency:              row.Currency,
		CreatedAt:             row.CreatedAt,
	}
}
