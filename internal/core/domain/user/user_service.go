package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/verification"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type Service struct {
	finder                    Finder
	creator                   Creator
	writer                    Writer
	verificationFinder        verification.Finder
	verificationWriter        verification.Writer
	verificationCodeGenerator CodeGenerator
	mailer                    VerificationMailer
	clock                     Clock
}

func NewService(
	finder Finder,
	creator Creator,
	writer Writer,
	verificationFinder verification.Finder,
	verificationWriter verification.Writer,
	verificationCodeGenerator CodeGenerator,
	mailer VerificationMailer,
	clock Clock,
) *Service {
	return &Service{
		finder:                    finder,
		creator:                   creator,
		writer:                    writer,
		verificationFinder:        verificationFinder,
		verificationWriter:        verificationWriter,
		verificationCodeGenerator: verificationCodeGenerator,
		mailer:                    mailer,
		clock:                     clock,
	}
}

func (s *Service) Create(ctx context.Context, newUser NewUser) (*User, error) {
	email := normalizeValue(newUser.Email)
	found, err := s.finder.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if found != nil {
		return nil, coreerror.New(coreerror.EmailAlreadyExists)
	}

	newUser.Email = email
	return s.creator.Create(ctx, newUser)
}

func (s *Service) FindByEmail(ctx context.Context, email string) (*User, error) {
	return s.finder.FindByEmail(ctx, normalizeValue(email))
}

func (s *Service) FindByID(ctx context.Context, userID string) (*User, error) {
	return s.finder.FindByID(ctx, userID)
}

func (s *Service) FindOrCreate(ctx context.Context, newUser NewUser) (*User, error) {
	email := normalizeValue(newUser.Email)
	found, err := s.finder.FindByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if found != nil {
		return found, nil
	}

	newUser.Email = email
	return s.creator.Create(ctx, newUser)
}

func (s *Service) VerifyEmail(ctx context.Context, code string) (*User, error) {
	foundVerification, err := s.verificationFinder.FindByCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("find verification by code: %w", err)
	}
	if foundVerification == nil || foundVerification.UsedAt != nil || foundVerification.ExpiresAt.Before(s.clock.Now()) {
		return nil, coreerror.New(coreerror.VerificationNotFound)
	}

	now := s.clock.Now()
	if err := s.verificationWriter.MarkUsed(ctx, code, now); err != nil {
		return nil, fmt.Errorf("mark verification used: %w", err)
	}

	if err := s.writer.MarkEmailVerified(ctx, foundVerification.UserID, now); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}

	return s.finder.FindByID(ctx, foundVerification.UserID)
}

func (s *Service) SendVerification(ctx context.Context, userID string, email string) (string, error) {
	now := s.clock.Now()
	code, err := s.verificationCodeGenerator.New()
	if err != nil {
		return "", fmt.Errorf("generate verification code: %w", err)
	}

	emailVerification := verification.EmailVerification{
		Code:      code,
		UserID:    userID,
		Email:     email,
		ExpiresAt: now.Add(10 * time.Minute),
		CreatedAt: now,
	}
	if err := s.verificationWriter.Create(ctx, emailVerification); err != nil {
		return "", fmt.Errorf("create verification: %w", err)
	}

	if err := s.mailer.SendVerification(ctx, email, code); err != nil {
		return "", fmt.Errorf("send verification email: %w", err)
	}

	return code, nil
}

func (s *Service) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) (*User, error) {
	if err := s.writer.MarkEmailVerified(ctx, userID, verifiedAt); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}

	return s.finder.FindByID(ctx, userID)
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
