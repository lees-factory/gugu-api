package aliexpress

import (
	"context"
	"time"
)

type SellerTokenRecord struct {
	ID                    string
	AppType               string
	SellerID              string
	HavanaID              string
	AppUserID             string
	UserNick              string
	Account               string
	AccountPlatform       string
	Locale                string
	SP                    string
	AccessToken           string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt *time.Time
	LastRefreshedAt       time.Time
	AuthorizedAt          time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type TokenStore interface {
	Upsert(ctx context.Context, token SellerTokenRecord) error
	FindOne(ctx context.Context) (*SellerTokenRecord, error)
	FindByAppType(ctx context.Context, appType string) (*SellerTokenRecord, error)
	FindBySellerID(ctx context.Context, sellerID string) (*SellerTokenRecord, error)
	ListExpiringBefore(ctx context.Context, expiresBefore time.Time) ([]SellerTokenRecord, error)
}
