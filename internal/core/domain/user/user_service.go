package user

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type RegisterEmailInput struct {
	Email        string
	DisplayName  string
	PasswordHash string
}

type FindOrCreateOAuthUserInput struct {
	Email       string
	DisplayName string
	AuthSource  string
	VerifiedAt  time.Time
}

type Service struct {
	finder  Finder
	creator Creator
	writer  Writer
}

func NewService(finder Finder, creator Creator, writer Writer) *Service {
	return &Service{
		finder:  finder,
		creator: creator,
		writer:  writer,
	}
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*User, error) {
	foundUser, err := s.finder.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return foundUser, nil
}

func (s *Service) FindByID(ctx context.Context, userID string) (*User, error) {
	foundUser, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return foundUser, nil
}

func (s *Service) RegisterEmail(ctx context.Context, input RegisterEmailInput) (*User, error) {
	emailValue := normalizeEmail(input.Email)
	foundUser, err := s.finder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	newUser, err := s.creator.Create(ctx, CreateInput{
		Email:         emailValue,
		DisplayName:   input.DisplayName,
		PasswordHash:  input.PasswordHash,
		AuthSource:    "email",
		EmailVerified: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create email user: %w", err)
	}

	return newUser, nil
}

func (s *Service) FindOrCreateOAuthUser(ctx context.Context, input FindOrCreateOAuthUserInput) (*User, error) {
	emailValue := normalizeEmail(input.Email)
	foundUser, err := s.finder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return foundUser, nil
	}

	newUser, err := s.creator.Create(ctx, CreateInput{
		Email:           emailValue,
		DisplayName:     input.DisplayName,
		AuthSource:      strings.TrimSpace(strings.ToLower(input.AuthSource)),
		EmailVerified:   true,
		EmailVerifiedAt: &input.VerifiedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}

	return newUser, nil
}

func (s *Service) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) (*User, error) {
	if err := s.writer.MarkEmailVerified(ctx, userID, verifiedAt); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}

	foundUser, err := s.finder.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return foundUser, nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
