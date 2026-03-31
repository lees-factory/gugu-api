package aliexpress

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type TokenRefresher interface {
	RefreshAccessToken(ctx context.Context, input clientaliexpress.RefreshTokenInput) (*clientaliexpress.TokenSet, error)
}

type tokenStoreProvider struct {
	mu         sync.Mutex
	appType    string
	tokenStore clientaliexpress.TokenStore
	refresher  TokenRefresher
}

func NewTokenProvider(appType string, tokenStore clientaliexpress.TokenStore, refresher TokenRefresher) TokenProvider {
	return &tokenStoreProvider{appType: appType, tokenStore: tokenStore, refresher: refresher}
}

func (p *tokenStoreProvider) GetAccessToken(ctx context.Context) (string, error) {
	record, err := p.tokenStore.FindByAppType(ctx, p.appType)
	if err != nil {
		return "", fmt.Errorf("find aliexpress token: %w", err)
	}
	if record == nil {
		return "", nil
	}

	if time.Now().After(record.AccessTokenExpiresAt) {
		refreshed, err := p.refreshOnce(ctx, record)
		if err != nil {
			return "", fmt.Errorf("auto refresh aliexpress token: %w", err)
		}
		return refreshed, nil
	}

	return record.AccessToken, nil
}

func (p *tokenStoreProvider) refreshOnce(ctx context.Context, record *clientaliexpress.SellerTokenRecord) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// lock 잡은 후 다시 확인 — 다른 goroutine이 이미 refresh 했을 수 있음
	fresh, err := p.tokenStore.FindByAppType(ctx, p.appType)
	if err != nil {
		return "", err
	}
	if fresh != nil && time.Now().Before(fresh.AccessTokenExpiresAt) {
		return fresh.AccessToken, nil
	}

	slog.Info("aliexpress access token expired, refreshing", "app_type", p.appType)

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

	slog.Info("aliexpress access token refreshed successfully", "app_type", p.appType)
	return record.AccessToken, nil
}
