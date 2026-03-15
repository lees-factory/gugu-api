package aliexpress

import (
	"context"
	"database/sql"
	"errors"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type SellerTokenSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *SellerTokenSQLCRepository {
	return &SellerTokenSQLCRepository{queries: sqldb.New(db)}
}

func (r *SellerTokenSQLCRepository) Upsert(ctx context.Context, token clientaliexpress.SellerTokenRecord) error {
	return r.queries.UpsertAliExpressSellerToken(ctx, sqldb.UpsertAliExpressSellerTokenParams{
		ID:                    token.ID,
		UserID:                token.UserID,
		SellerID:              token.SellerID,
		HavanaID:              token.HavanaID,
		AppUserID:             token.AppUserID,
		UserNick:              token.UserNick,
		Account:               token.Account,
		AccountPlatform:       token.AccountPlatform,
		Locale:                token.Locale,
		Sp:                    token.SP,
		AccessToken:           token.AccessToken,
		RefreshToken:          token.RefreshToken,
		AccessTokenExpiresAt:  token.AccessTokenExpiresAt,
		RefreshTokenExpiresAt: nullTime(token.RefreshTokenExpiresAt),
		LastRefreshedAt:       token.LastRefreshedAt,
		AuthorizedAt:          token.AuthorizedAt,
		CreatedAt:             token.CreatedAt,
		UpdatedAt:             token.UpdatedAt,
	})
}

func (r *SellerTokenSQLCRepository) FindByUserID(ctx context.Context, userID string) (*clientaliexpress.SellerTokenRecord, error) {
	row, err := r.queries.FindAliExpressSellerTokenByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	record := toSellerTokenRecord(row)
	return &record, nil
}

func (r *SellerTokenSQLCRepository) FindBySellerID(ctx context.Context, sellerID string) (*clientaliexpress.SellerTokenRecord, error) {
	row, err := r.queries.FindAliExpressSellerTokenBySellerID(ctx, sellerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	record := toSellerTokenRecord(row)
	return &record, nil
}

func (r *SellerTokenSQLCRepository) ListExpiringBefore(ctx context.Context, expiresBefore time.Time) ([]clientaliexpress.SellerTokenRecord, error) {
	rows, err := r.queries.ListAliExpressSellerTokensExpiringBefore(ctx, expiresBefore)
	if err != nil {
		return nil, err
	}

	items := make([]clientaliexpress.SellerTokenRecord, 0, len(rows))
	for _, row := range rows {
		items = append(items, toSellerTokenRecord(row))
	}

	return items, nil
}

func toSellerTokenRecord(row sqldb.GuguAliexpressSellerToken) clientaliexpress.SellerTokenRecord {
	token := clientaliexpress.SellerTokenRecord{
		ID:                   row.ID,
		UserID:               row.UserID,
		SellerID:             row.SellerID,
		HavanaID:             row.HavanaID,
		AppUserID:            row.AppUserID,
		UserNick:             row.UserNick,
		Account:              row.Account,
		AccountPlatform:      row.AccountPlatform,
		Locale:               row.Locale,
		SP:                   row.Sp,
		AccessToken:          row.AccessToken,
		RefreshToken:         row.RefreshToken,
		AccessTokenExpiresAt: row.AccessTokenExpiresAt,
		LastRefreshedAt:      row.LastRefreshedAt,
		AuthorizedAt:         row.AuthorizedAt,
		CreatedAt:            row.CreatedAt,
		UpdatedAt:            row.UpdatedAt,
	}
	if row.RefreshTokenExpiresAt.Valid {
		token.RefreshTokenExpiresAt = &row.RefreshTokenExpiresAt.Time
	}
	return token
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *value, Valid: true}
}
