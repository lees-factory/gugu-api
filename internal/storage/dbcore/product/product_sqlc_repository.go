package product

import (
	"context"
	"database/sql"
	"errors"

	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	"github.com/ljj/gugu-api/internal/core/enum"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type ProductSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *ProductSQLCRepository {
	return &ProductSQLCRepository{queries: sqldb.New(db)}
}

func (r *ProductSQLCRepository) FindByID(ctx context.Context, productID string) (*domainproduct.Product, error) {
	row, err := r.queries.FindProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainProduct(row)
	return &item, nil
}

func (r *ProductSQLCRepository) FindByIDs(ctx context.Context, productIDs []string) ([]domainproduct.Product, error) {
	rows, err := r.queries.FindProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	products := make([]domainproduct.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomainProduct(row)
	}
	return products, nil
}

func (r *ProductSQLCRepository) FindByMarketAndExternalProductID(ctx context.Context, market enum.Market, externalProductID string) (*domainproduct.Product, error) {
	row, err := r.queries.FindProductByMarketAndExternalProductID(ctx, sqldb.FindProductByMarketAndExternalProductIDParams{
		Market:            string(market),
		ExternalProductID: externalProductID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	item := toDomainProduct(row)
	return &item, nil
}

func (r *ProductSQLCRepository) Create(ctx context.Context, product domainproduct.Product) error {
	return r.queries.CreateProduct(ctx, sqldb.CreateProductParams{
		ID:                product.ID,
		Market:            string(product.Market),
		ExternalProductID: product.ExternalProductID,
		OriginalUrl:       product.OriginalURL,
		Title:             product.Title,
		MainImageUrl:      product.MainImageURL,
		ProductUrl:        product.ProductURL,
		PromotionLink:     product.PromotionLink,
		CollectionSource:  product.CollectionSource,
		LastCollectedAt:   product.LastCollectedAt,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	})
}

func (r *ProductSQLCRepository) Update(ctx context.Context, product domainproduct.Product) error {
	affected, err := r.queries.UpdateProduct(ctx, sqldb.UpdateProductParams{
		ID:               product.ID,
		OriginalUrl:      product.OriginalURL,
		Title:            product.Title,
		MainImageUrl:     product.MainImageURL,
		ProductUrl:       product.ProductURL,
		PromotionLink:    product.PromotionLink,
		CollectionSource: product.CollectionSource,
		LastCollectedAt:  product.LastCollectedAt,
		UpdatedAt:        product.UpdatedAt,
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductSQLCRepository) ListByMarket(ctx context.Context, market enum.Market) ([]domainproduct.Product, error) {
	rows, err := r.queries.ListProductsByMarket(ctx, string(market))
	if err != nil {
		return nil, err
	}
	products := make([]domainproduct.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomainProduct(row)
	}
	return products, nil
}

func (r *ProductSQLCRepository) ListByCollectionSource(ctx context.Context, source string, limit int, offset int) ([]domainproduct.Product, error) {
	rows, err := r.queries.ListProductsByCollectionSource(ctx, sqldb.ListProductsByCollectionSourceParams{
		CollectionSource: source,
		Limit:            int32(limit),
		Offset:           int32(offset),
	})
	if err != nil {
		return nil, err
	}
	products := make([]domainproduct.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomainProduct(row)
	}
	return products, nil
}

func toDomainProduct(row sqldb.GuguProduct) domainproduct.Product {
	return domainproduct.Product{
		ID:                row.ID,
		Market:            enum.Market(row.Market),
		ExternalProductID: row.ExternalProductID,
		OriginalURL:       row.OriginalUrl,
		Title:             row.Title,
		MainImageURL:      row.MainImageUrl,
		ProductURL:        row.ProductUrl,
		PromotionLink:     row.PromotionLink,
		CollectionSource:  row.CollectionSource,
		LastCollectedAt:   row.LastCollectedAt,
		CreatedAt:         row.CreatedAt,
		UpdatedAt:         row.UpdatedAt,
	}
}
