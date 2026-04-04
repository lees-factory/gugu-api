package pricealert

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type IDGenerator interface {
	New() (string, error)
}

type Clock interface {
	Now() time.Time
}

type Service struct {
	finder      Finder
	writer      Writer
	idGenerator IDGenerator
	clock       Clock
}

func NewService(finder Finder, writer Writer, idGenerator IDGenerator, clock Clock) *Service {
	return &Service{
		finder:      finder,
		writer:      writer,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (s *Service) Register(ctx context.Context, userID string, skuID string, channel string) (*PriceAlert, error) {
	userID = strings.TrimSpace(userID)
	skuID = strings.TrimSpace(skuID)
	if channel == "" {
		channel = "EMAIL"
	}

	existing, err := s.finder.FindByUserIDAndSKUID(ctx, userID, skuID)
	if err != nil {
		return nil, fmt.Errorf("find existing alert: %w", err)
	}
	if existing != nil {
		if !existing.Enabled {
			if err := s.writer.UpdateEnabled(ctx, existing.ID, true); err != nil {
				return nil, fmt.Errorf("re-enable alert: %w", err)
			}
			existing.Enabled = true
		}
		return existing, nil
	}

	alertID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate alert id: %w", err)
	}

	alert := PriceAlert{
		ID:        alertID,
		UserID:    userID,
		SKUID:     skuID,
		Channel:   channel,
		Enabled:   true,
		CreatedAt: s.clock.Now(),
	}
	if err := s.writer.Create(ctx, alert); err != nil {
		return nil, fmt.Errorf("create alert: %w", err)
	}

	return &alert, nil
}

func (s *Service) Unregister(ctx context.Context, userID string, skuID string) error {
	found, err := s.finder.FindByUserIDAndSKUID(ctx, strings.TrimSpace(userID), strings.TrimSpace(skuID))
	if err != nil {
		return fmt.Errorf("find alert: %w", err)
	}
	if found == nil || !found.Enabled {
		return nil
	}
	if err := s.writer.UpdateEnabled(ctx, found.ID, false); err != nil {
		return fmt.Errorf("disable alert: %w", err)
	}
	return nil
}

func (s *Service) ListBySKUID(ctx context.Context, skuID string) ([]PriceAlert, error) {
	return s.finder.ListBySKUID(ctx, strings.TrimSpace(skuID))
}

func (s *Service) ListByUserID(ctx context.Context, userID string) ([]PriceAlert, error) {
	return s.finder.ListByUserID(ctx, strings.TrimSpace(userID))
}

func (s *Service) ListByProductID(ctx context.Context, productID string) ([]PriceAlert, error) {
	return s.finder.ListByProductID(ctx, strings.TrimSpace(productID))
}

func (s *Service) ListByProductIDs(ctx context.Context, productIDs []string) ([]PriceAlert, error) {
	return s.finder.ListByProductIDs(ctx, productIDs)
}
