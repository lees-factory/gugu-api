package product

import (
	"context"
	"database/sql"
	"errors"
	"time"

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
	item := toDomainProductFromFindByID(row)
	return &item, nil
}

func (r *ProductSQLCRepository) FindByIDs(ctx context.Context, productIDs []string) ([]domainproduct.Product, error) {
	rows, err := r.queries.FindProductsByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	}
	products := make([]domainproduct.Product, len(rows))
	for i, row := range rows {
		products[i] = toDomainProductFromFindByIDs(row)
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
	item := toDomainProductFromFindByMarketAndExternalProductID(row)
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
		products[i] = toDomainProductFromListByMarket(row)
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
		products[i] = toDomainProductFromListByCollectionSource(row)
	}
	return products, nil
}

func toDomainProductValues(
	id string,
	market string,
	externalProductID string,
	originalURL string,
	title string,
	mainImageURL string,
	productURL string,
	promotionLink string,
	collectionSource string,
	lastCollectedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) domainproduct.Product {
	return domainproduct.Product{
		ID:                id,
		Market:            enum.Market(market),
		ExternalProductID: externalProductID,
		OriginalURL:       originalURL,
		Title:             title,
		MainImageURL:      mainImageURL,
		ProductURL:        productURL,
		PromotionLink:     promotionLink,
		CollectionSource:  collectionSource,
		LastCollectedAt:   lastCollectedAt,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}
}

func toDomainProductFromFindByID(row sqldb.FindProductByIDRow) domainproduct.Product {
	return toDomainProductValues(
		row.ID,
		row.Market,
		row.ExternalProductID,
		row.OriginalUrl,
		row.Title,
		row.MainImageUrl,
		row.ProductUrl,
		row.PromotionLink,
		row.CollectionSource,
		row.LastCollectedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func toDomainProductFromFindByIDs(row sqldb.FindProductsByIDsRow) domainproduct.Product {
	return toDomainProductValues(
		row.ID,
		row.Market,
		row.ExternalProductID,
		row.OriginalUrl,
		row.Title,
		row.MainImageUrl,
		row.ProductUrl,
		row.PromotionLink,
		row.CollectionSource,
		row.LastCollectedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func toDomainProductFromFindByMarketAndExternalProductID(row sqldb.FindProductByMarketAndExternalProductIDRow) domainproduct.Product {
	return toDomainProductValues(
		row.ID,
		row.Market,
		row.ExternalProductID,
		row.OriginalUrl,
		row.Title,
		row.MainImageUrl,
		row.ProductUrl,
		row.PromotionLink,
		row.CollectionSource,
		row.LastCollectedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func toDomainProductFromListByMarket(row sqldb.ListProductsByMarketRow) domainproduct.Product {
	return toDomainProductValues(
		row.ID,
		row.Market,
		row.ExternalProductID,
		row.OriginalUrl,
		row.Title,
		row.MainImageUrl,
		row.ProductUrl,
		row.PromotionLink,
		row.CollectionSource,
		row.LastCollectedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func toDomainProductFromListByCollectionSource(row sqldb.ListProductsByCollectionSourceRow) domainproduct.Product {
	return toDomainProductValues(
		row.ID,
		row.Market,
		row.ExternalProductID,
		row.OriginalUrl,
		row.Title,
		row.MainImageUrl,
		row.ProductUrl,
		row.PromotionLink,
		row.CollectionSource,
		row.LastCollectedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}
