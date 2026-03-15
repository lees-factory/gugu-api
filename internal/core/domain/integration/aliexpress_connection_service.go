package integration

import (
	"context"
	"fmt"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

type IDGenerator interface {
	New() (string, error)
}

type AliExpressConnectionService struct {
	client      clientaliexpress.Client
	tokenStore  clientaliexpress.TokenStore
	idGenerator IDGenerator
	now         func() time.Time
}

type BuildAliExpressAuthorizationURLInput struct {
	UserID string
}

type BuildAliExpressAuthorizationURLResult struct {
	UserID           string
	AuthorizationURL string
}

type ExchangeAliExpressCodeInput struct {
	UserID string
	Code   string
}

type AliExpressConnectionStatus struct {
	UserID                  string
	SellerID                string
	Account                 string
	UserNick                string
	Connected               bool
	ReauthorizationRequired bool
	AuthorizedAt            time.Time
	AccessTokenExpiresAt    time.Time
	RefreshTokenExpiresAt   *time.Time
	LastRefreshedAt         time.Time
}

func NewAliExpressConnectionService(client clientaliexpress.Client, tokenStore clientaliexpress.TokenStore, idGenerator IDGenerator) *AliExpressConnectionService {
	return &AliExpressConnectionService{
		client:      client,
		tokenStore:  tokenStore,
		idGenerator: idGenerator,
		now:         time.Now,
	}
}

func (s *AliExpressConnectionService) BuildAuthorizationURL(_ context.Context, input BuildAliExpressAuthorizationURLInput) (*BuildAliExpressAuthorizationURLResult, error) {
	url, err := s.client.BuildAuthorizationURL()
	if err != nil {
		return nil, err
	}

	return &BuildAliExpressAuthorizationURLResult{
		UserID:           input.UserID,
		AuthorizationURL: url,
	}, nil
}

func (s *AliExpressConnectionService) ExchangeCode(ctx context.Context, input ExchangeAliExpressCodeInput) (*AliExpressConnectionStatus, error) {
	tokenSet, err := s.client.ExchangeCode(ctx, clientaliexpress.TokenExchangeInput{Code: input.Code})
	if err != nil {
		return nil, err
	}

	now := s.now()
	recordID, err := s.resolveRecordID(ctx, input.UserID, tokenSet.SellerID)
	if err != nil {
		return nil, err
	}

	record := clientaliexpress.SellerTokenRecord{
		ID:                   recordID,
		UserID:               input.UserID,
		SellerID:             tokenSet.SellerID,
		HavanaID:             tokenSet.HavanaID,
		AppUserID:            tokenSet.UserID,
		UserNick:             tokenSet.UserNick,
		Account:              tokenSet.Account,
		AccountPlatform:      tokenSet.AccountPlatform,
		Locale:               tokenSet.Locale,
		SP:                   tokenSet.SP,
		AccessToken:          tokenSet.AccessToken,
		RefreshToken:         tokenSet.RefreshToken,
		AccessTokenExpiresAt: time.UnixMilli(tokenSet.ExpireTime),
		LastRefreshedAt:      now,
		AuthorizedAt:         now,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
	if tokenSet.RefreshTokenValidTime > 0 {
		refreshExpiresAt := time.UnixMilli(tokenSet.RefreshTokenValidTime)
		record.RefreshTokenExpiresAt = &refreshExpiresAt
	}

	existingRecord, err := s.tokenStore.FindBySellerID(ctx, tokenSet.SellerID)
	if err != nil {
		return nil, fmt.Errorf("find existing aliexpress token: %w", err)
	}
	if existingRecord != nil {
		record.CreatedAt = existingRecord.CreatedAt
		record.AuthorizedAt = existingRecord.AuthorizedAt
	}

	if err := s.tokenStore.Upsert(ctx, record); err != nil {
		return nil, fmt.Errorf("upsert aliexpress token: %w", err)
	}

	return s.buildStatus(record, now), nil
}

func (s *AliExpressConnectionService) GetConnectionStatus(ctx context.Context, userID string) (*AliExpressConnectionStatus, error) {
	record, err := s.tokenStore.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return &AliExpressConnectionStatus{
			UserID:                  userID,
			Connected:               false,
			ReauthorizationRequired: false,
		}, nil
	}

	return s.buildStatus(*record, s.now()), nil
}

func (s *AliExpressConnectionService) resolveRecordID(ctx context.Context, userID string, sellerID string) (string, error) {
	existingRecord, err := s.tokenStore.FindBySellerID(ctx, sellerID)
	if err != nil {
		return "", err
	}
	if existingRecord != nil {
		return existingRecord.ID, nil
	}

	existingRecord, err = s.tokenStore.FindByUserID(ctx, userID)
	if err != nil {
		return "", err
	}
	if existingRecord != nil {
		return existingRecord.ID, nil
	}

	return s.idGenerator.New()
}

func (s *AliExpressConnectionService) buildStatus(record clientaliexpress.SellerTokenRecord, now time.Time) *AliExpressConnectionStatus {
	reauthorizationRequired := false
	if record.RefreshTokenExpiresAt != nil && !record.RefreshTokenExpiresAt.After(now) {
		reauthorizationRequired = true
	}

	return &AliExpressConnectionStatus{
		UserID:                  record.UserID,
		SellerID:                record.SellerID,
		Account:                 record.Account,
		UserNick:                record.UserNick,
		Connected:               true,
		ReauthorizationRequired: reauthorizationRequired,
		AuthorizedAt:            record.AuthorizedAt,
		AccessTokenExpiresAt:    record.AccessTokenExpiresAt,
		RefreshTokenExpiresAt:   record.RefreshTokenExpiresAt,
		LastRefreshedAt:         record.LastRefreshedAt,
	}
}
