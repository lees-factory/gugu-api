package integration

import (
	"context"
	"testing"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

func TestAliExpressConnectionServiceExchangeCode(t *testing.T) {
	now := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	store := &stubTokenStore{}
	service := NewAliExpressConnectionService(&stubAliExpressClient{
		tokenSet: &clientaliexpress.TokenSet{
			AccessToken:           "access-token",
			RefreshToken:          "refresh-token",
			ExpireTime:            now.Add(1 * time.Hour).UnixMilli(),
			RefreshTokenValidTime: now.Add(24 * time.Hour).UnixMilli(),
			HavanaID:              "havana-1",
			UserID:                "ae-user-1",
			SellerID:              "seller-1",
			UserNick:              "seller nick",
			Account:               "seller@example.com",
			AccountPlatform:       "seller_center",
			Locale:                "ko_KR",
			SP:                    "ae",
		},
	}, store, stubIDGenerator{id: "record-1"})
	service.now = func() time.Time { return now }

	result, err := service.ExchangeCode(context.Background(), ExchangeAliExpressCodeInput{
		Code: "code-1",
	})
	if err != nil {
		t.Fatalf("ExchangeCode() error = %v", err)
	}

	if result.SellerID != "seller-1" {
		t.Fatalf("SellerID = %q", result.SellerID)
	}
	if !result.Connected {
		t.Fatal("Connected = false, want true")
	}
	if store.lastUpsert == nil || store.lastUpsert.SellerID != "seller-1" {
		t.Fatalf("lastUpsert = %#v", store.lastUpsert)
	}
}

func TestAliExpressConnectionServiceGetConnectionStatusWithoutRecord(t *testing.T) {
	service := NewAliExpressConnectionService(&stubAliExpressClient{}, &stubTokenStore{}, stubIDGenerator{id: "record-1"})

	result, err := service.GetConnectionStatus(context.Background())
	if err != nil {
		t.Fatalf("GetConnectionStatus() error = %v", err)
	}

	if result.Connected {
		t.Fatal("Connected = true, want false")
	}
	if result.ReauthorizationRequired {
		t.Fatal("ReauthorizationRequired = true, want false")
	}
}

type stubAliExpressClient struct {
	tokenSet *clientaliexpress.TokenSet
}

func (c *stubAliExpressClient) BuildAuthorizationURL() (string, error) {
	return "https://api-sg.aliexpress.com/oauth/authorize", nil
}

func (c *stubAliExpressClient) ExchangeCode(_ context.Context, _ clientaliexpress.TokenExchangeInput) (*clientaliexpress.TokenSet, error) {
	return c.tokenSet, nil
}

func (c *stubAliExpressClient) RefreshAccessToken(_ context.Context, _ clientaliexpress.RefreshTokenInput) (*clientaliexpress.TokenSet, error) {
	return c.tokenSet, nil
}

func (c *stubAliExpressClient) GetProductSnapshot(_ context.Context, _ clientaliexpress.ProductLookupInput) (*clientaliexpress.ProductSnapshot, error) {
	return nil, nil
}

func (c *stubAliExpressClient) GetAffiliateProductDetail(_ context.Context, _ clientaliexpress.ProductDetailInput) (*clientaliexpress.ProductDetailResult, error) {
	return nil, nil
}

func (c *stubAliExpressClient) GetAffiliateProductSKUDetail(_ context.Context, _ clientaliexpress.ProductSKUDetailInput) (*clientaliexpress.ProductSKUDetailResult, error) {
	return nil, nil
}

type stubTokenStore struct {
	items      map[string]clientaliexpress.SellerTokenRecord
	lastUpsert *clientaliexpress.SellerTokenRecord
}

func (s *stubTokenStore) Upsert(_ context.Context, token clientaliexpress.SellerTokenRecord) error {
	if s.items == nil {
		s.items = map[string]clientaliexpress.SellerTokenRecord{}
	}
	s.items[token.SellerID] = token
	s.lastUpsert = &token
	return nil
}

func (s *stubTokenStore) FindOne(_ context.Context) (*clientaliexpress.SellerTokenRecord, error) {
	for _, item := range s.items {
		found := item
		return &found, nil
	}
	return nil, nil
}

func (s *stubTokenStore) FindBySellerID(_ context.Context, sellerID string) (*clientaliexpress.SellerTokenRecord, error) {
	if item, ok := s.items[sellerID]; ok {
		found := item
		return &found, nil
	}
	return nil, nil
}

func (s *stubTokenStore) ListExpiringBefore(_ context.Context, _ time.Time) ([]clientaliexpress.SellerTokenRecord, error) {
	return nil, nil
}

type stubIDGenerator struct {
	id string
}

func (g stubIDGenerator) New() (string, error) {
	return g.id, nil
}
