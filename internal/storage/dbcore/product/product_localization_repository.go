package product

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
)

type ProductLocalizationRepository struct {
	db *sql.DB
}

func NewLocalizationRepository(db *sql.DB) *ProductLocalizationRepository {
	return &ProductLocalizationRepository{db: db}
}

func (r *ProductLocalizationRepository) FindByProductIDLanguageCurrency(ctx context.Context, productID string, language string, currency string) (*domainproduct.Variant, error) {
	const query = `
SELECT
    product_id,
    language,
    title,
    main_image_url,
    product_url,
    updated_at
FROM gugu.product_localization
WHERE product_id = $1
  AND language = $2
`

	var (
		foundProductID string
		foundLanguage  string
		title          string
		mainImageURL   string
		productURL     string
		updatedAt      time.Time
	)

	err := r.db.QueryRowContext(
		ctx,
		query,
		strings.TrimSpace(productID),
		strings.ToUpper(strings.TrimSpace(language)),
	).Scan(
		&foundProductID,
		&foundLanguage,
		&title,
		&mainImageURL,
		&productURL,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &domainproduct.Variant{
		ProductID:    foundProductID,
		Language:     foundLanguage,
		Currency:     strings.ToUpper(strings.TrimSpace(currency)),
		Title:        title,
		MainImageURL: mainImageURL,
		ProductURL:   productURL,
		CurrentPrice: "",
		CreatedAt:    updatedAt,
		UpdatedAt:    updatedAt,
	}, nil
}

func (r *ProductLocalizationRepository) FindByLookupKeys(ctx context.Context, keys []domainproduct.VariantLookupKey) ([]domainproduct.Variant, error) {
	items := make([]domainproduct.Variant, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		mapKey := strings.TrimSpace(key.ProductID) + ":" + strings.ToUpper(strings.TrimSpace(key.Language)) + ":" + strings.ToUpper(strings.TrimSpace(key.Currency))
		if _, ok := seen[mapKey]; ok {
			continue
		}
		seen[mapKey] = struct{}{}

		found, err := r.FindByProductIDLanguageCurrency(ctx, key.ProductID, key.Language, key.Currency)
		if err != nil {
			return nil, err
		}
		if found == nil {
			continue
		}
		items = append(items, *found)
	}

	return items, nil
}

func (r *ProductLocalizationRepository) Upsert(ctx context.Context, variant domainproduct.Variant) error {
	const query = `
INSERT INTO gugu.product_localization (
    product_id,
    language,
    title,
    main_image_url,
    product_url,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, NOW()
)
ON CONFLICT (product_id, language)
DO UPDATE SET
    title = EXCLUDED.title,
    main_image_url = EXCLUDED.main_image_url,
    product_url = EXCLUDED.product_url,
    updated_at = NOW()
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		strings.TrimSpace(variant.ProductID),
		strings.ToUpper(strings.TrimSpace(variant.Language)),
		strings.TrimSpace(variant.Title),
		strings.TrimSpace(variant.MainImageURL),
		strings.TrimSpace(variant.ProductURL),
	)
	if err != nil {
		return fmt.Errorf("upsert product localization: %w", err)
	}
	return nil
}
