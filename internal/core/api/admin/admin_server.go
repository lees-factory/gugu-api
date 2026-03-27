package admin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/controller/advice"
	apirequest "github.com/ljj/gugu-api/internal/core/api/controller/v1/request"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainps "github.com/ljj/gugu-api/internal/core/domain/pricesnapshot"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	apiresponse "github.com/ljj/gugu-api/internal/core/support/response"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
	"github.com/ljj/gugu-api/internal/provider/product/batch"
	dbcorepricehistory "github.com/ljj/gugu-api/internal/storage/dbcore/pricehistory"
	dbcoresnapshot "github.com/ljj/gugu-api/internal/storage/dbcore/pricesnapshot"
	memorypricehistory "github.com/ljj/gugu-api/internal/storage/memory/pricehistory"
	memorysnapshot "github.com/ljj/gugu-api/internal/storage/memory/pricesnapshot"
	"github.com/ljj/gugu-api/internal/support/config"
)

func RegisterBatchRoutes(
	router chi.Router,
	cfg config.Config,
	db *sql.DB,
	tokenStore clientaliexpress.TokenStore,
	productService *domainproduct.Service,
) {
	priceUpdater := wirePriceUpdater(cfg, db, tokenStore, productService)
	snapshotRecorder := wireSnapshotRecorder(db, productService)
	hotProductLoader := wireHotProductLoader(cfg, tokenStore, productService)

	router.Post("/v1/batch/update-prices", apiadvice.Wrap(func(r *stdhttp.Request) (int, any, error) {
		if err := priceUpdater.UpdateAll(r.Context()); err != nil {
			return 0, nil, err
		}
		return stdhttp.StatusOK, apiresponse.Success(), nil
	}))

	router.Post("/v1/batch/record-snapshots", apiadvice.Wrap(func(r *stdhttp.Request) (int, any, error) {
		if err := snapshotRecorder.RecordDailySnapshots(r.Context()); err != nil {
			return 0, nil, err
		}
		return stdhttp.StatusOK, apiresponse.Success(), nil
	}))

	router.Post("/v1/batch/load-hot-products", apiadvice.Wrap(func(r *stdhttp.Request) (int, any, error) {
		var req apirequest.LoadHotProductsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return stdhttp.StatusBadRequest, nil, fmt.Errorf("decode request: %w", err)
		}
		result, err := hotProductLoader.LoadHotProducts(r.Context(), batch.HotProductLoadInput{
			CategoryIDs:    req.CategoryIDs,
			Keywords:       req.Keywords,
			PageNo:         req.PageNo,
			PageSize:       req.PageSize,
			Sort:           req.Sort,
			MinSalePrice:   req.MinSalePrice,
			MaxSalePrice:   req.MaxSalePrice,
			ShipToCountry:  req.ShipToCountry,
			TargetCurrency: req.TargetCurrency,
			TargetLanguage: req.TargetLanguage,
		})
		if err != nil {
			return 0, nil, err
		}
		return stdhttp.StatusOK, apiresponse.SuccessWithData(result), nil
	}))
}

func wirePriceUpdater(cfg config.Config, db *sql.DB, tokenStore clientaliexpress.TokenStore, productService *domainproduct.Service) *batch.PriceUpdater {
	aliExpressClient, _ := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})

	tokenProvider := provideraliexpress.NewTokenProvider(tokenStore)
	fetcher := provideraliexpress.NewBatchFetcher(aliExpressClient, tokenProvider)

	var priceHistoryWriter domainpricehistory.Writer
	if db != nil {
		priceHistoryWriter = domainpricehistory.NewWriter(dbcorepricehistory.NewSQLCRepository(db))
	} else {
		priceHistoryWriter = domainpricehistory.NewWriter(memorypricehistory.NewRepository())
	}

	return batch.NewPriceUpdater(productService, priceHistoryWriter, fetcher)
}

func wireSnapshotRecorder(db *sql.DB, productService *domainproduct.Service) *batch.SnapshotRecorder {
	var productSnapshotWriter domainps.ProductSnapshotWriter
	var skuSnapshotWriter domainps.SKUSnapshotWriter
	if db != nil {
		productSnapshotWriter = domainps.NewProductSnapshotWriter(dbcoresnapshot.NewProductSnapshotSQLCRepository(db))
		skuSnapshotWriter = domainps.NewSKUSnapshotWriter(dbcoresnapshot.NewSKUSnapshotSQLCRepository(db))
	} else {
		productSnapshotWriter = domainps.NewProductSnapshotWriter(memorysnapshot.NewProductSnapshotRepository())
		skuSnapshotWriter = domainps.NewSKUSnapshotWriter(memorysnapshot.NewSKUSnapshotRepository())
	}
	return batch.NewSnapshotRecorder(productService, productSnapshotWriter, skuSnapshotWriter)
}

func wireHotProductLoader(cfg config.Config, tokenStore clientaliexpress.TokenStore, productService *domainproduct.Service) *batch.HotProductLoader {
	aliExpressClient, _ := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})
	tokenProvider := provideraliexpress.NewTokenProvider(tokenStore)
	return batch.NewHotProductLoader(aliExpressClient, tokenProvider, productService)
}
