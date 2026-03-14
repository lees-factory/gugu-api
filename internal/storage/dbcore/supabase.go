package dbcore

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const simpleProtocolParam = "default_query_exec_mode"
const searchPathParam = "search_path"

func OpenSupabase(ctx context.Context, databaseURL string) (*sql.DB, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("database url is empty")
	}

	dsn, err := normalizeSupabaseTransactionPoolerURL(databaseURL)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

func normalizeSupabaseTransactionPoolerURL(databaseURL string) (string, error) {
	parsedURL, err := url.Parse(databaseURL)
	if err != nil {
		return "", fmt.Errorf("parse database url: %w", err)
	}

	query := parsedURL.Query()
	if query.Get("sslmode") == "" {
		query.Set("sslmode", "require")
	}
	if query.Get(simpleProtocolParam) == "" {
		query.Set(simpleProtocolParam, "simple_protocol")
	}
	if query.Get(searchPathParam) == "" {
		query.Set(searchPathParam, "gugu")
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String(), nil
}
