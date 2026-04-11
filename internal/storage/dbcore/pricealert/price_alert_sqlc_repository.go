package pricealert

import (
	"context"
	"database/sql"
	"errors"

	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type SQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *SQLCRepository {
	return &SQLCRepository{queries: sqldb.New(db)}
}

func (r *SQLCRepository) FindByUserIDAndSKUID(ctx context.Context, userID string, skuID string) (*domainpricealert.PriceAlert, error) {
	row, err := r.queries.FindPriceAlertByUserIDAndSKUID(ctx, sqldb.FindPriceAlertByUserIDAndSKUIDParams{
		UserID: userID,
		SkuID:  skuID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	alert := toDomain(row)
	return &alert, nil
}

func (r *SQLCRepository) ListBySKUID(ctx context.Context, skuID string) ([]domainpricealert.PriceAlert, error) {
	rows, err := r.queries.ListPriceAlertsBySKUID(ctx, skuID)
	if err != nil {
		return nil, err
	}
	return toDomainList(rows), nil
}

func (r *SQLCRepository) ListByProductID(ctx context.Context, productID string) ([]domainpricealert.PriceAlert, error) {
	rows, err := r.queries.ListPriceAlertsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	return toDomainList(rows), nil
}

func (r *SQLCRepository) ListByProductIDs(ctx context.Context, productIDs []string) ([]domainpricealert.PriceAlert, error) {
	rows, err := r.queries.ListPriceAlertsByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	return toDomainList(rows), nil
}

func (r *SQLCRepository) ListByUserID(ctx context.Context, userID string) ([]domainpricealert.PriceAlert, error) {
	rows, err := r.queries.ListPriceAlertsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return toDomainList(rows), nil
}

func (r *SQLCRepository) Create(ctx context.Context, alert domainpricealert.PriceAlert) error {
	return r.queries.CreatePriceAlert(ctx, sqldb.CreatePriceAlertParams{
		ID:        alert.ID,
		UserID:    alert.UserID,
		SkuID:     alert.SKUID,
		Channel:   alert.Channel,
		Enabled:   alert.Enabled,
		CreatedAt: alert.CreatedAt,
	})
}

func (r *SQLCRepository) UpdateEnabled(ctx context.Context, alertID string, enabled bool) error {
	_, err := r.queries.UpdatePriceAlertEnabled(ctx, sqldb.UpdatePriceAlertEnabledParams{
		ID:      alertID,
		Enabled: enabled,
	})
	return err
}

func (r *SQLCRepository) UpdateSettings(ctx context.Context, alertID string, channel string, enabled bool) error {
	_, err := r.queries.UpdatePriceAlertSettings(ctx, sqldb.UpdatePriceAlertSettingsParams{
		ID:      alertID,
		Channel: channel,
		Enabled: enabled,
	})
	return err
}

func toDomain(row sqldb.GuguPriceAlert) domainpricealert.PriceAlert {
	return domainpricealert.PriceAlert{
		ID:        row.ID,
		UserID:    row.UserID,
		SKUID:     row.SkuID,
		Channel:   row.Channel,
		Enabled:   row.Enabled,
		CreatedAt: row.CreatedAt,
	}
}

func toDomainList(rows []sqldb.GuguPriceAlert) []domainpricealert.PriceAlert {
	alerts := make([]domainpricealert.PriceAlert, len(rows))
	for i, row := range rows {
		alerts[i] = toDomain(row)
	}
	return alerts
}
