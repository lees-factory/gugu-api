package api

import (
	"database/sql"
	"fmt"

	"github.com/go-chi/chi/v5"
	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiintegration "github.com/ljj/gugu-api/internal/core/api/v1/integration"
	domainintegration "github.com/ljj/gugu-api/internal/core/domain/integration"
	dbcorealiexpress "github.com/ljj/gugu-api/internal/storage/dbcore/aliexpress"
	memoryaliexpress "github.com/ljj/gugu-api/internal/storage/memory/aliexpress"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
)

func registerAliExpressIntegrationModule(router chi.Router, cfg config.Config, db *sql.DB) error {
	aliExpressTokenStore := buildAliExpressTokenStore(db)
	aliExpressRecordIDGenerator := id.NewRandomHexGenerator(16)

	aliExpressClient, err := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})
	if err != nil {
		return fmt.Errorf("build aliexpress client: %w", err)
	}

	aliExpressConnectionService := domainintegration.NewAliExpressConnectionService(
		aliExpressClient,
		aliExpressTokenStore,
		aliExpressRecordIDGenerator,
	)

	aliExpressController := apiintegration.NewAliExpressController(aliExpressConnectionService)
	router.Route("/v1/integrations/aliexpress", func(r chi.Router) {
		r.Post("/authorize-url", apiadvice.Wrap(aliExpressController.BuildAuthorizationURL))
		r.Post("/exchange-code", apiadvice.Wrap(aliExpressController.ExchangeCode))
		r.Get("/connection-status", apiadvice.Wrap(aliExpressController.GetConnectionStatus))
	})

	return nil
}

func buildAliExpressTokenStore(db *sql.DB) clientaliexpress.TokenStore {
	if db == nil {
		return memoryaliexpress.NewRepository()
	}

	return dbcorealiexpress.NewSQLCRepository(db)
}
