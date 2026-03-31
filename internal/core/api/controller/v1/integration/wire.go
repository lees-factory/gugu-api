package integration

import (
	"database/sql"
	"fmt"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
	domainintegration "github.com/ljj/gugu-api/internal/core/domain/integration"
	dbcorealiexpress "github.com/ljj/gugu-api/internal/storage/dbcore/aliexpress"
	memoryaliexpress "github.com/ljj/gugu-api/internal/storage/memory/aliexpress"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
)

func Wire(cfg config.Config, db *sql.DB) (*AliExpressController, clientaliexpress.TokenStore, error) {
	tokenStore := buildAliExpressTokenStore(db)
	recordIDGenerator := id.NewRandomHexGenerator(16)

	services := make(map[string]*domainintegration.AliExpressConnectionService)

	affiliateClient, err := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
		BaseURL:     cfg.AliExpressBaseURL,
		AppKey:      cfg.AliExpressAppKey,
		AppSecret:   cfg.AliExpressAppSecret,
		CallbackURL: cfg.AliExpressCallbackURL,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("build aliexpress affiliate client: %w", err)
	}
	services["AFFILIATE"] = domainintegration.NewAliExpressConnectionService(
		"AFFILIATE", affiliateClient, tokenStore, recordIDGenerator,
	)

	if cfg.AliExpressDSAppKey != "" {
		dsClient, err := clientaliexpress.NewHTTPClient(clientaliexpress.Config{
			BaseURL:     cfg.AliExpressDSBaseURL,
			AppKey:      cfg.AliExpressDSAppKey,
			AppSecret:   cfg.AliExpressDSAppSecret,
			CallbackURL: cfg.AliExpressDSCallbackURL,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("build aliexpress ds client: %w", err)
		}
		services["DROPSHIPPING"] = domainintegration.NewAliExpressConnectionService(
			"DROPSHIPPING", dsClient, tokenStore, recordIDGenerator,
		)
	}

	controller := NewAliExpressController(services, affiliateClient, tokenStore)
	return controller, tokenStore, nil
}

func buildAliExpressTokenStore(db *sql.DB) clientaliexpress.TokenStore {
	if db == nil {
		return memoryaliexpress.NewRepository()
	}
	return dbcorealiexpress.NewSQLCRepository(db)
}
