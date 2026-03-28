package aliexpress

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type TokenRefresher interface {
	RefreshAccessToken(ctx context.Context, input clientaliexpress.RefreshTokenInput) (*clientaliexpress.TokenSet, error)
}

type tokenStoreProvider struct {
	tokenStore clientaliexpress.TokenStore
	refresher  TokenRefresher
}

func NewTokenProvider(tokenStore clientaliexpress.TokenStore, refresher TokenRefresher) TokenProvider {
	return &tokenStoreProvider{tokenStore: tokenStore, refresher: refresher}
}

func (p *tokenStoreProvider) GetAccessToken(ctx context.Context) (string, error) {
	record, err := p.tokenStore.FindOne(ctx)
	if err != nil {
		return "", fmt.Errorf("find aliexpress token: %w", err)
	}
	if record == nil {
		return "", nil
	}

	if time.Now().After(record.AccessTokenExpiresAt) {
		refreshed, err := p.refresh(ctx, record)
		if err != nil {
			return "", fmt.Errorf("auto refresh aliexpress token: %w", err)
		}
		return refreshed, nil
	}

	return record.AccessToken, nil
}

func (p *tokenStoreProvider) refresh(ctx context.Context, record *clientaliexpress.SellerTokenRecord) (string, error) {
	slog.Info("aliexpress access token expired, refreshing automatically")

	tokenSet, err := p.refresher.RefreshAccessToken(ctx, clientaliexpress.RefreshTokenInput{
		RefreshToken: record.RefreshToken,
	})
	if err != nil {
		return "", err
	}

	now := time.Now()
	record.AccessToken = tokenSet.AccessToken
	record.RefreshToken = tokenSet.RefreshToken
	record.AccessTokenExpiresAt = time.UnixMilli(tokenSet.ExpireTime)
	record.LastRefreshedAt = now
	record.UpdatedAt = now
	if tokenSet.RefreshTokenValidTime > 0 {
		refreshExpiresAt := time.UnixMilli(tokenSet.RefreshTokenValidTime)
		record.RefreshTokenExpiresAt = &refreshExpiresAt
	}

	if err := p.tokenStore.Upsert(ctx, *record); err != nil {
		return "", fmt.Errorf("upsert refreshed token: %w", err)
	}

	slog.Info("aliexpress access token refreshed successfully")
	return record.AccessToken, nil
}
