package response

import (
	"time"

	domainintegration "github.com/ljj/gugu-api/internal/core/domain/integration"
)

type AliExpressAuthorizationURL struct {
	UserID           string `json:"user_id"`
	AuthorizationURL string `json:"authorization_url"`
}

type AliExpressConnectionStatus struct {
	UserID                  string     `json:"user_id"`
	SellerID                string     `json:"seller_id,omitempty"`
	Account                 string     `json:"account,omitempty"`
	UserNick                string     `json:"user_nick,omitempty"`
	Connected               bool       `json:"connected"`
	ReauthorizationRequired bool       `json:"reauthorization_required"`
	AuthorizedAt            *time.Time `json:"authorized_at,omitempty"`
	AccessTokenExpiresAt    *time.Time `json:"access_token_expires_at,omitempty"`
	RefreshTokenExpiresAt   *time.Time `json:"refresh_token_expires_at,omitempty"`
	LastRefreshedAt         *time.Time `json:"last_refreshed_at,omitempty"`
}

func NewAliExpressAuthorizationURL(source domainintegration.BuildAliExpressAuthorizationURLResult) AliExpressAuthorizationURL {
	return AliExpressAuthorizationURL{
		UserID:           source.UserID,
		AuthorizationURL: source.AuthorizationURL,
	}
}

func NewAliExpressConnectionStatus(source domainintegration.AliExpressConnectionStatus) AliExpressConnectionStatus {
	result := AliExpressConnectionStatus{
		UserID:                  source.UserID,
		SellerID:                source.SellerID,
		Account:                 source.Account,
		UserNick:                source.UserNick,
		Connected:               source.Connected,
		ReauthorizationRequired: source.ReauthorizationRequired,
	}
	if !source.AuthorizedAt.IsZero() {
		result.AuthorizedAt = &source.AuthorizedAt
	}
	if !source.AccessTokenExpiresAt.IsZero() {
		result.AccessTokenExpiresAt = &source.AccessTokenExpiresAt
	}
	if source.RefreshTokenExpiresAt != nil {
		result.RefreshTokenExpiresAt = source.RefreshTokenExpiresAt
	}
	if !source.LastRefreshedAt.IsZero() {
		result.LastRefreshedAt = &source.LastRefreshedAt
	}
	return result
}
