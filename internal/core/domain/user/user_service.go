package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/verification"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type RegisterEmailInput struct {
	Email       string
	Password    string
	DisplayName string
}

type RegisterEmailResult struct {
	User                   User
	VerificationCode       string
	VerificationDispatched bool
}

type VerifyEmailInput struct {
	Code string
}

type VerifyEmailResult struct {
	User User
}

type FindOrCreateOAuthUserInput struct {
	Email       string
	DisplayName string
	AuthSource  string
	VerifiedAt  time.Time
}

type Service struct {
	finder                    Finder
	creator                   Creator
	writer                    Writer
	verificationFinder        verification.Finder
	verificationWriter        verification.Writer
	verificationCodeGenerator CodeGenerator
	passwordHasher            PasswordHasher
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
	passwordHasher PasswordHasher,
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
		passwordHasher:            passwordHasher,
		mailer:                    mailer,
		clock:                     clock,
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

func (s *Service) RegisterEmail(ctx context.Context, input RegisterEmailInput) (*RegisterEmailResult, error) {
	passwordHash, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	createdUser, err := s.createEmailUser(ctx, input, passwordHash)
	if err != nil {
		return nil, err
	}

	verificationResult, err := s.issueEmailVerification(ctx, createdUser.ID, createdUser.Email)
	if err != nil {
		return nil, err
	}

	return &RegisterEmailResult{
		User:                   *createdUser,
		VerificationCode:       verificationResult.code,
		VerificationDispatched: verificationResult.dispatched,
	}, nil
}

func (s *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) (*VerifyEmailResult, error) {
	foundVerification, err := s.verificationFinder.FindByCode(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("find verification by code: %w", err)
	}
	if foundVerification == nil || foundVerification.UsedAt != nil || foundVerification.ExpiresAt.Before(s.clock.Now()) {
		return nil, coreerror.New(coreerror.VerificationNotFound)
	}

	now := s.clock.Now()
	if err := s.verificationWriter.MarkUsed(ctx, input.Code, now); err != nil {
		return nil, fmt.Errorf("mark verification used: %w", err)
	}

	if err := s.writer.MarkEmailVerified(ctx, foundVerification.UserID, now); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}

	verifiedUser, err := s.finder.FindByID(ctx, foundVerification.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if verifiedUser == nil {
		return nil, coreerror.New(coreerror.VerificationNotFound)
	}

	return &VerifyEmailResult{User: *verifiedUser}, nil
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

type issueVerificationResult struct {
	code       string
	dispatched bool
}

func (s *Service) issueEmailVerification(ctx context.Context, userID string, email string) (*issueVerificationResult, error) {
	now := s.clock.Now()
	code, err := s.verificationCodeGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate verification code: %w", err)
	}

	emailVerification := verification.EmailVerification{
		Code:      code,
		UserID:    userID,
		Email:     email,
		ExpiresAt: now.Add(10 * time.Minute),
		CreatedAt: now,
	}
	if err := s.verificationWriter.Create(ctx, emailVerification); err != nil {
		return nil, fmt.Errorf("create verification: %w", err)
	}

	if err := s.mailer.SendVerification(ctx, email, code); err != nil {
		return nil, fmt.Errorf("send verification email: %w", err)
	}

	return &issueVerificationResult{
		code:       code,
		dispatched: true,
	}, nil
}

func (s *Service) createEmailUser(ctx context.Context, input RegisterEmailInput, passwordHash string) (*User, error) {
	emailValue := normalizeEmail(input.Email)
	foundUser, err := s.finder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return nil, coreerror.New(coreerror.EmailAlreadyExists)
	}

	newUser, err := s.creator.Create(ctx, CreateInput{
		Email:         emailValue,
		DisplayName:   input.DisplayName,
		PasswordHash:  passwordHash,
		AuthSource:    "email",
		EmailVerified: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create email user: %w", err)
	}
	return newUser, nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
