package api

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	apiauth "github.com/ljj/gugu-api/internal/core/api/v1/auth"
	apiintegration "github.com/ljj/gugu-api/internal/core/api/v1/integration"
	apiproduct "github.com/ljj/gugu-api/internal/core/api/v1/product"
	apitrackeditem "github.com/ljj/gugu-api/internal/core/api/v1/trackeditem"
	domainpricehistory "github.com/ljj/gugu-api/internal/core/domain/pricehistory"
	domainproduct "github.com/ljj/gugu-api/internal/core/domain/product"
	provideraliexpress "github.com/ljj/gugu-api/internal/provider/product/aliexpress"
	"github.com/ljj/gugu-api/internal/provider/product/batch"
	dbcorepricehistory "github.com/ljj/gugu-api/internal/storage/dbcore/pricehistory"
	memorypricehistory "github.com/ljj/gugu-api/internal/storage/memory/pricehistory"
	"github.com/ljj/gugu-api/internal/support/config"
)

type Server struct {
	router *chi.Mux
}

func NewServer(cfg config.Config, db *sql.DB) (*Server, error) {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(CORSMiddleware(cfg.CORSAllowedOrigins))

	registerHealthRoute(router)

	authControllers, err := apiauth.Wire(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("wire auth: %w", err)
	}
	authControllers.Auth.RegisterRoutes(router)
	authControllers.User.RegisterRoutes(router)

	aliExpressController, aliExpressTokenStore, err := apiintegration.Wire(cfg, db)
	if err != nil {
		return nil, fmt.Errorf("wire aliexpress integration: %w", err)
	}
	aliExpressController.RegisterRoutes(router)

	trackedItemController, trackedItemService, productService, err := apitrackeditem.Wire(cfg, db, aliExpressTokenStore)
	if err != nil {
		return nil, fmt.Errorf("wire tracked item: %w", err)
	}
	trackedItemController.RegisterRoutes(router)

	productController := apiproduct.Wire(db, productService, trackedItemService)
	productController.RegisterRoutes(router)

	priceUpdater := wirePriceUpdater(cfg, db, aliExpressTokenStore, productService)
	registerBatchRoutes(router, priceUpdater)

	return &Server{router: router}, nil
}

func (s *Server) Handler() stdhttp.Handler {
	return s.router
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

func registerBatchRoutes(router chi.Router, priceUpdater *batch.PriceUpdater) {
	router.Post("/v1/batch/update-prices", apiadvice.Wrap(func(r *stdhttp.Request) (int, any, error) {
		if err := priceUpdater.UpdateAll(r.Context()); err != nil {
			return 0, nil, err
		}
		return stdhttp.StatusOK, apiresponse.Success(), nil
	}))
}

func registerHealthRoute(router chi.Router) {
	router.Get("/health", apiadvice.Wrap(func(_ *stdhttp.Request) (int, any, error) {
		return stdhttp.StatusOK, apiresponse.SuccessWithData(map[string]string{"status": "ok"}), nil
	}))
}
