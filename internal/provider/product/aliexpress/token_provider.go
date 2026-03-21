package aliexpress

import (
	"context"
	"fmt"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type tokenStoreProvider struct {
	tokenStore clientaliexpress.TokenStore
}

func NewTokenProvider(tokenStore clientaliexpress.TokenStore) TokenProvider {
	return &tokenStoreProvider{tokenStore: tokenStore}
}

func (p *tokenStoreProvider) GetAccessToken(ctx context.Context) (string, error) {
	record, err := p.tokenStore.FindOne(ctx)
	if err != nil {
		return "", fmt.Errorf("find aliexpress token: %w", err)
	}
	if record == nil {
		return "", nil
	}
	return record.AccessToken, nil
}
