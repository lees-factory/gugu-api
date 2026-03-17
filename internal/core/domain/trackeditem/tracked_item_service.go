package trackeditem

import (
	"context"
	"fmt"
	"strings"

	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type AddInput struct {
	UserID      string
	ProductID   string
	OriginalURL string
}

type AddResult struct {
	TrackedItem    TrackedItem
	AlreadyTracked bool
}

type Service struct {
	finder      Finder
	writer      Writer
	idGenerator IDGenerator
	clock       Clock
}

func NewService(
	finder Finder,
	writer Writer,
	idGenerator IDGenerator,
	clock Clock,
) *Service {
	return &Service{
		finder:      finder,
		writer:      writer,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (s *Service) Add(ctx context.Context, input AddInput) (*AddResult, error) {
	found, err := s.finder.FindByUserIDAndProductID(ctx, strings.TrimSpace(input.UserID), input.ProductID)
	if err != nil {
		return nil, fmt.Errorf("find tracked item by user id and product id: %w", err)
	}
	if found != nil {
		return &AddResult{
			TrackedItem:    *found,
			AlreadyTracked: true,
		}, nil
	}

	trackedItemID, err := s.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate tracked item id: %w", err)
	}

	tracked := TrackedItem{
		ID:          trackedItemID,
		UserID:      strings.TrimSpace(input.UserID),
		ProductID:   input.ProductID,
		OriginalURL: strings.TrimSpace(input.OriginalURL),
		CreatedAt:   s.clock.Now(),
	}
	if err := s.writer.Create(ctx, tracked); err != nil {
		return nil, fmt.Errorf("create tracked item: %w", err)
	}

	return &AddResult{
		TrackedItem:    tracked,
		AlreadyTracked: false,
	}, nil
}

func (s *Service) ListByUserID(ctx context.Context, userID string) ([]TrackedItem, error) {
	return s.finder.ListByUserID(ctx, strings.TrimSpace(userID))
}

func (s *Service) FindByUserIDAndProductID(ctx context.Context, userID string, productID string) (*TrackedItem, error) {
	return s.finder.FindByUserIDAndProductID(ctx, strings.TrimSpace(userID), strings.TrimSpace(productID))
}

func (s *Service) Delete(ctx context.Context, trackedItemID string, userID string) error {
	found, err := s.finder.FindByIDAndUserID(ctx, strings.TrimSpace(trackedItemID), strings.TrimSpace(userID))
	if err != nil {
		return fmt.Errorf("find tracked item by id and user id: %w", err)
	}
	if found == nil {
		return coreerror.New(coreerror.TrackedItemNotFound)
	}

	if err := s.writer.DeleteByIDAndUserID(ctx, found.ID, found.UserID); err != nil {
		return fmt.Errorf("delete tracked item by id and user id: %w", err)
	}

	return nil
}
