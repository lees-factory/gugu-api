package product

import (
	"context"
	"fmt"
	"strings"
)

type CreateInput struct {
	Market            Market
	ExternalProductID string
	OriginalURL       string
	Title             string
	MainImageURL      string
	CurrentPrice      string
	Currency          string
	ProductURL        string
	CollectionSource  string
}

type Service struct {
	repository  Repository
	idGenerator IDGenerator
	clock       Clock
}

func NewService(repository Repository, idGenerator IDGenerator, clock Clock) *Service {
	return &Service{
		repository:  repository,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (*Product, error) {
	productID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate product id: %w", err)
	}

	now := s.clock.Now()
	item := Product{
		ID:                productID,
		Market:            input.Market,
		ExternalProductID: strings.TrimSpace(input.ExternalProductID),
		OriginalURL:       strings.TrimSpace(input.OriginalURL),
		Title:             strings.TrimSpace(input.Title),
		MainImageURL:      strings.TrimSpace(input.MainImageURL),
		CurrentPrice:      strings.TrimSpace(input.CurrentPrice),
		Currency:          strings.TrimSpace(input.Currency),
		ProductURL:        strings.TrimSpace(input.ProductURL),
		CollectionSource:  strings.TrimSpace(input.CollectionSource),
		LastCollectedAt:   now,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repository.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return &item, nil
}

func (s *Service) FindByMarketAndExternalProductID(ctx context.Context, market Market, externalProductID string) (*Product, error) {
	item, err := s.repository.FindByMarketAndExternalProductID(ctx, market, strings.TrimSpace(externalProductID))
	if err != nil {
		return nil, fmt.Errorf("find product by market and external product id: %w", err)
	}
	return item, nil
}
