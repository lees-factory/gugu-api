package aliexpress

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func (c *HTTPClient) BuildAuthorizationURL() (string, error) {
	authURL, err := url.Parse(c.baseURL + "/oauth/authorize")
	if err != nil {
		return "", fmt.Errorf("parse authorization url: %w", err)
	}

	query := authURL.Query()
	query.Set("response_type", "code")
	query.Set("force_auth", "true")
	query.Set("redirect_uri", c.callbackURL)
	query.Set("client_id", c.appKey)
	authURL.RawQuery = query.Encode()

	return authURL.String(), nil
}

func (c *HTTPClient) ExchangeCode(ctx context.Context, input TokenExchangeInput) (*TokenSet, error) {
	code := strings.TrimSpace(input.Code)
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}

	return c.executeTokenRequest(ctx, "/auth/token/create", map[string]string{
		"code": code,
	})
}

func (c *HTTPClient) RefreshAccessToken(ctx context.Context, input RefreshTokenInput) (*TokenSet, error) {
	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, fmt.Errorf("refresh token is required")
	}

	return c.executeTokenRequest(ctx, "/auth/token/refresh", map[string]string{
		"refresh_token": refreshToken,
	})
}

func (c *HTTPClient) executeTokenRequest(ctx context.Context, apiName string, form map[string]string) (*TokenSet, error) {
	response, err := c.executeFormRequest(ctx, signedRequest{
		apiName: apiName,
		form:    form,
	})
	if err != nil {
		return nil, err
	}

	var payload tokenResponse
	if err := json.Unmarshal(response, &payload); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	if payload.AccessToken == "" {
		return nil, fmt.Errorf("aliexpress token response missing access token")
	}

	return &TokenSet{
		AccessToken:           payload.AccessToken,
		RefreshToken:          payload.RefreshToken,
		ExpiresIn:             payload.ExpiresIn,
		RefreshExpiresIn:      payload.RefreshExpiresIn,
		ExpireTime:            payload.ExpireTime,
		RefreshTokenValidTime: payload.RefreshTokenValidTime,
		HavanaID:              payload.HavanaID,
		UserID:                payload.UserID,
		SellerID:              payload.SellerID,
		UserNick:              payload.UserNick,
		Account:               payload.Account,
		Locale:                payload.Locale,
		AccountPlatform:       payload.AccountPlatform,
		SP:                    payload.SP,
		RequestID:             payload.RequestID,
		Code:                  payload.Code,
	}, nil
}
