package product

import (
	"context"
	"database/sql"
	"errors"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type ProductVariantSQLCRepository struct {
	queries *sqldb.Queries
}

func NewVariantSQLCRepository(db *sql.DB) *ProductVariantSQLCRepository {
	return &ProductVariantSQLCRepository{queries: sqldb.New(db)}
}

func (r *ProductVariantSQLCRepository) FindByProductIDLanguageCurrency(ctx context.Context, productID string, language string, currency string) (*domainproduct.Variant, error) {
	row, err := r.queries.FindProductVariantByProductIDLanguageCurrency(ctx, sqldb.FindProductVariantByProductIDLanguageCurrencyParams{
		ProductID: productID,
		Language:  language,
		Currency:  currency,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainProductVariant(row)
	return &item, nil
}

func (r *ProductVariantSQLCRepository) FindByLookupKeys(ctx context.Context, keys []domainproduct.VariantLookupKey) ([]domainproduct.Variant, error) {
	productIDs := make([]string, 0, len(keys))
	languages := make([]string, 0, len(keys))
	currencies := make([]string, 0, len(keys))
	for _, key := range keys {
		productIDs = append(productIDs, key.ProductID)
		languages = append(languages, key.Language)
		currencies = append(currencies, key.Currency)
	}

	rows, err := r.queries.FindProductVariantsByLookupKeys(ctx, sqldb.FindProductVariantsByLookupKeysParams{
		Column1: productIDs,
		Column2: languages,
		Column3: currencies,
	})
	if err != nil {
		return nil, err
	}

	items := make([]domainproduct.Variant, 0, len(rows))
	for _, row := range rows {
		items = append(items, toDomainProductVariant(row))
	}
	return items, nil
}

func (r *ProductVariantSQLCRepository) Upsert(ctx context.Context, variant domainproduct.Variant) error {
	return r.queries.UpsertProductVariant(ctx, sqldb.UpsertProductVariantParams{
		ProductID:       variant.ProductID,
		Language:        variant.Language,
		Currency:        variant.Currency,
		Title:           variant.Title,
		MainImageUrl:    variant.MainImageURL,
		ProductUrl:      variant.ProductURL,
		CurrentPrice:    variant.CurrentPrice,
		LastCollectedAt: toNullTime(variant.LastCollectedAt),
		CreatedAt:       variant.CreatedAt,
		UpdatedAt:       variant.UpdatedAt,
	})
}

func toDomainProductVariant(row sqldb.GuguProductVariant) domainproduct.Variant {
	return domainproduct.Variant{
		ProductID:       row.ProductID,
		Language:        row.Language,
		Currency:        row.Currency,
		Title:           row.Title,
		MainImageURL:    row.MainImageUrl,
		ProductURL:      row.ProductUrl,
		CurrentPrice:    row.CurrentPrice,
		LastCollectedAt: fromNullTime(row.LastCollectedAt),
		CreatedAt:       row.CreatedAt,
		UpdatedAt:       row.UpdatedAt,
	}
}

func toNullTime(t *sql.NullTime) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return *t
}

func fromNullTime(t sql.NullTime) *sql.NullTime {
	if !t.Valid {
		return nil
	}
	copy := t
	return &copy
}
