package product

import (
	"context"
	"database/sql"
	"errors"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type ProductSKUSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSKUSQLCRepository(db *sql.DB) *ProductSKUSQLCRepository {
	return &ProductSKUSQLCRepository{queries: sqldb.New(db)}
}

func (r *ProductSKUSQLCRepository) Create(ctx context.Context, sku domainproduct.ProductSKU) error {
	return r.queries.CreateProductSKU(ctx, sqldb.CreateProductSKUParams{
		ID:            sku.ID,
		ProductID:     sku.ProductID,
		ExternalSkuID: sku.ExternalSKUID,
		SkuName:       sku.SKUName,
		Color:         sku.Color,
		Size:          sku.Size,
		Price:         sku.Price,
		OriginalPrice: sku.OriginalPrice,
		Currency:      sku.Currency,
		ImageUrl:      sku.ImageURL,
		SkuProperties: sku.SKUProperties,
		CreatedAt:     sku.CreatedAt,
		UpdatedAt:     sku.UpdatedAt,
	})
}

func (r *ProductSKUSQLCRepository) FindByID(ctx context.Context, skuID string) (*domainproduct.ProductSKU, error) {
	row, err := r.queries.FindProductSKUByID(ctx, skuID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainProductSKU(row)
	return &item, nil
}

func (r *ProductSKUSQLCRepository) FindByProductID(ctx context.Context, productID string) ([]domainproduct.ProductSKU, error) {
	rows, err := r.queries.FindProductSKUsByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	skus := make([]domainproduct.ProductSKU, len(rows))
	for i, row := range rows {
		skus[i] = toDomainProductSKU(row)
	}
	return skus, nil
}

func toDomainProductSKU(row sqldb.GuguProductSku) domainproduct.ProductSKU {
	return domainproduct.ProductSKU{
		ID:            row.ID,
		ProductID:     row.ProductID,
		ExternalSKUID: row.ExternalSkuID,
		SKUName:       row.SkuName,
		Color:         row.Color,
		Size:          row.Size,
		Price:         row.Price,
		OriginalPrice: row.OriginalPrice,
		Currency:      row.Currency,
		ImageURL:      row.ImageUrl,
		SKUProperties: row.SkuProperties,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}
