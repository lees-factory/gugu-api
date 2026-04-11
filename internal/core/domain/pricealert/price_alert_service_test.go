package pricealert_test

import (
	"context"
	"testing"
	"time"

	domainpricealert "github.com/ljj/gugu-api/internal/core/domain/pricealert"
	memorypricealert "github.com/ljj/gugu-api/internal/storage/memory/pricealert"
)

type fixedIDGenerator struct{}

func (fixedIDGenerator) New() (string, error) {
	return "alert-fixed-id", nil
}

type fixedClock struct{}

func (fixedClock) Now() time.Time {
	return time.Date(2026, 4, 11, 0, 0, 0, 0, time.UTC)
}

func TestRegister_NormalizesDefaultChannelToEmail(t *testing.T) {
	repo := memorypricealert.NewRepository()
	service := domainpricealert.NewService(
		domainpricealert.NewFinder(repo),
		domainpricealert.NewWriter(repo),
		fixedIDGenerator{},
		fixedClock{},
	)

	alert, err := service.Register(context.Background(), "user-1", "sku-1", "email")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if alert.Channel != "EMAIL" {
		t.Fatalf("Register().Channel = %q, want EMAIL", alert.Channel)
	}
	if !alert.Enabled {
		t.Fatalf("Register().Enabled = false, want true")
	}
}

func TestRegister_ReenablesAndUpdatesChannelForExistingAlert(t *testing.T) {
	repo := memorypricealert.NewRepository()
	service := domainpricealert.NewService(
		domainpricealert.NewFinder(repo),
		domainpricealert.NewWriter(repo),
		fixedIDGenerator{},
		fixedClock{},
	)

	if _, err := service.Register(context.Background(), "user-1", "sku-1", "EMAIL"); err != nil {
		t.Fatalf("Register() initial error = %v", err)
	}
	if err := service.Unregister(context.Background(), "user-1", "sku-1"); err != nil {
		t.Fatalf("Unregister() error = %v", err)
	}

	updated, err := service.Register(context.Background(), "user-1", "sku-1", "sms")
	if err != nil {
		t.Fatalf("Register() update error = %v", err)
	}
	if updated.Channel != "SMS" {
		t.Fatalf("Register() updated channel = %q, want SMS", updated.Channel)
	}
	if !updated.Enabled {
		t.Fatalf("Register() updated enabled = false, want true")
	}
}
