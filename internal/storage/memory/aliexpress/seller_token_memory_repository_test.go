package aliexpress

import (
	"context"
	"testing"
	"time"

	clientaliexpress "github.com/ljj/gugu-api/internal/clients/aliexpress"
)

func TestSellerTokenMemoryRepository(t *testing.T) {
	repository := NewRepository()
	now := time.Date(2026, 3, 15, 10, 0, 0, 0, time.UTC)
	refreshExpiresAt := now.Add(24 * time.Hour)

	token := clientaliexpress.SellerTokenRecord{
		ID:                    "ali-token-1",
		UserID:                "user-1",
		SellerID:              "seller-1",
		HavanaID:              "havana-1",
		AppUserID:             "app-user-1",
		UserNick:              "seller-nick",
		Account:               "seller@example.com",
		AccountPlatform:       "seller_center",
		Locale:                "ko_KR",
		SP:                    "ae",
		AccessToken:           "access-token",
		RefreshToken:          "refresh-token",
		AccessTokenExpiresAt:  now.Add(30 * time.Minute),
		RefreshTokenExpiresAt: &refreshExpiresAt,
		LastRefreshedAt:       now,
		AuthorizedAt:          now,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := repository.Upsert(context.Background(), token); err != nil {
		t.Fatalf("Upsert() error = %v", err)
	}

	foundByUserID, err := repository.FindByUserID(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}
	if foundByUserID == nil || foundByUserID.SellerID != "seller-1" {
		t.Fatalf("FindByUserID() = %#v", foundByUserID)
	}

	foundBySellerID, err := repository.FindBySellerID(context.Background(), "seller-1")
	if err != nil {
		t.Fatalf("FindBySellerID() error = %v", err)
	}
	if foundBySellerID == nil || foundBySellerID.UserID != "user-1" {
		t.Fatalf("FindBySellerID() = %#v", foundBySellerID)
	}

	expiringItems, err := repository.ListExpiringBefore(context.Background(), now.Add(1*time.Hour))
	if err != nil {
		t.Fatalf("ListExpiringBefore() error = %v", err)
	}
	if len(expiringItems) != 1 {
		t.Fatalf("len(ListExpiringBefore()) = %d, want 1", len(expiringItems))
	}
}
