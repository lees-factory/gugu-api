package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/user"
	"github.com/ljj/gugu-api/internal/core/domain/verification"
)

type Service struct {
	userFinder              user.Finder
	userCreator             user.Creator
	userWriter              user.Writer
	emailVerificationFinder verification.Finder
	emailVerificationWriter verification.Writer
	oauthIdentityFinder     OAuthIdentityFinder
	oauthIdentityWriter     OAuthIdentityWriter
	sessionAppender         SessionAppender
	identityIDGenerator     IDGenerator
	tokenGenerator          TokenGenerator
	passwordHasher          PasswordHasher
	clock                   Clock
	emailSender             VerificationSender
}

func New(
	userFinder user.Finder,
	userCreator user.Creator,
	userWriter user.Writer,
	emailVerificationFinder verification.Finder,
	emailVerificationWriter verification.Writer,
	oauthIdentityFinder OAuthIdentityFinder,
	oauthIdentityWriter OAuthIdentityWriter,
	sessionAppender SessionAppender,
	identityIDGenerator IDGenerator,
	tokenGenerator TokenGenerator,
	passwordHasher PasswordHasher,
	clock Clock,
	emailSender VerificationSender,
) *Service {
	return &Service{
		userFinder:              userFinder,
		userCreator:             userCreator,
		userWriter:              userWriter,
		emailVerificationFinder: emailVerificationFinder,
		emailVerificationWriter: emailVerificationWriter,
		oauthIdentityFinder:     oauthIdentityFinder,
		oauthIdentityWriter:     oauthIdentityWriter,
		sessionAppender:         sessionAppender,
		identityIDGenerator:     identityIDGenerator,
		tokenGenerator:          tokenGenerator,
		passwordHasher:          passwordHasher,
		clock:                   clock,
		emailSender:             emailSender,
	}
}

func (s *Service) RegisterEmail(ctx context.Context, input RegisterEmailInput) (*RegisterEmailResult, error) {
	passwordHash, err := s.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	createdUser, err := s.createEmailUser(ctx, input, passwordHash)
	if err != nil {
		return nil, err
	}

	verificationResult, err := s.IssueEmailVerification(ctx, IssueEmailVerificationInput{
		UserID: createdUser.ID,
		Email:  createdUser.Email,
	})
	if err != nil {
		return nil, err
	}

	return &RegisterEmailResult{
		User:                   *createdUser,
		VerificationToken:      verificationResult.VerificationToken,
		VerificationDispatched: verificationResult.VerificationDispatched,
	}, nil
}

func (s *Service) LoginEmail(ctx context.Context, input LoginEmailInput) (*LoginResult, error) {
	foundUser, err := s.userFinder.FindByEmail(ctx, normalizeValue(input.Email))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser == nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.VerifyPassword(foundUser.PasswordHash, input.Password); err != nil {
		return nil, err
	}
	if !foundUser.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	session, err := s.CreateSession(ctx, foundUser.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:    *foundUser,
		Session: *session,
	}, nil
}

func (s *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) (*VerifyEmailResult, error) {
	verifyResult, err := s.VerifyEmailToken(ctx, VerifyEmailTokenInput{Token: input.Token})
	if err != nil {
		return nil, err
	}

	if err := s.userWriter.MarkEmailVerified(ctx, verifyResult.UserID, verifyResult.VerifiedAt); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}
	verifiedUser, err := s.userFinder.FindByID(ctx, verifyResult.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if verifiedUser == nil {
		return nil, ErrVerificationNotFound
	}

	return &VerifyEmailResult{User: *verifiedUser}, nil
}

func (s *Service) LoginOAuth(ctx context.Context, input OAuthLoginInput) (*LoginResult, error) {
	foundIdentity, err := s.FindOAuthIdentity(ctx, input.Provider, input.Subject)
	if err != nil {
		return nil, err
	}

	var foundUser *user.User
	if foundIdentity != nil {
		if err := s.TouchOAuthIdentityLastLogin(ctx, input.Provider, input.Subject); err != nil {
			return nil, err
		}
		foundUser, err = s.userFinder.FindByID(ctx, foundIdentity.UserID)
		if err != nil {
			return nil, fmt.Errorf("find user by id: %w", err)
		}
		if foundUser == nil {
			return nil, ErrInvalidCredentials
		}
	} else {
		emailValue := normalizeValue(input.Email)
		foundUser, err = s.userFinder.FindByEmail(ctx, emailValue)
		if err != nil {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
		if foundUser == nil {
			now := s.Now()
			foundUser, err = s.userCreator.Create(ctx, user.CreateInput{
				Email:           emailValue,
				DisplayName:     input.DisplayName,
				AuthSource:      string(input.Provider),
				EmailVerified:   true,
				EmailVerifiedAt: &now,
			})
			if err != nil {
				return nil, fmt.Errorf("create oauth user: %w", err)
			}
		}
		if _, err := s.CreateOAuthIdentity(ctx, CreateOAuthIdentityInput{
			UserID:   foundUser.ID,
			Provider: input.Provider,
			Subject:  input.Subject,
			Email:    foundUser.Email,
		}); err != nil {
			return nil, err
		}
	}

	session, err := s.CreateSession(ctx, foundUser.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:    *foundUser,
		Session: *session,
	}, nil
}

func (s *Service) createEmailUser(ctx context.Context, input RegisterEmailInput, passwordHash string) (*user.User, error) {
	emailValue := normalizeValue(input.Email)
	foundUser, err := s.userFinder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return nil, user.ErrEmailAlreadyExists
	}

	newUser, err := s.userCreator.Create(ctx, user.CreateInput{
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

func (s *Service) IssueEmailVerification(ctx context.Context, input IssueEmailVerificationInput) (*IssueEmailVerificationResult, error) {
	now := s.clock.Now()
	token, err := s.tokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate verification token: %w", err)
	}

	emailVerification := verification.EmailVerification{
		Token:     token,
		UserID:    input.UserID,
		Email:     input.Email,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.emailVerificationWriter.Create(ctx, emailVerification); err != nil {
		return nil, fmt.Errorf("create verification: %w", err)
	}

	if err := s.emailSender.SendVerification(ctx, input.Email, token); err != nil {
		return nil, fmt.Errorf("send verification email: %w", err)
	}

	return &IssueEmailVerificationResult{
		VerificationToken:      token,
		VerificationDispatched: true,
	}, nil
}

func (s *Service) VerifyPassword(passwordHash string, password string) error {
	if err := s.passwordHasher.Verify(passwordHash, password); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func (s *Service) CreateSession(ctx context.Context, userID string) (*Session, error) {
	token, err := s.tokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	now := s.clock.Now()
	session := Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.sessionAppender.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &session, nil
}

func (s *Service) VerifyEmailToken(ctx context.Context, input VerifyEmailTokenInput) (*VerifyEmailTokenResult, error) {
	foundVerification, err := s.emailVerificationFinder.FindByToken(ctx, input.Token)
	if err != nil {
		return nil, fmt.Errorf("find verification by token: %w", err)
	}
	if foundVerification == nil || foundVerification.UsedAt != nil || foundVerification.ExpiresAt.Before(s.clock.Now()) {
		return nil, ErrVerificationNotFound
	}

	now := s.clock.Now()
	if err := s.emailVerificationWriter.MarkUsed(ctx, input.Token, now); err != nil {
		return nil, fmt.Errorf("mark verification used: %w", err)
	}

	return &VerifyEmailTokenResult{
		UserID:     foundVerification.UserID,
		VerifiedAt: now,
	}, nil
}

func (s *Service) FindOAuthIdentity(ctx context.Context, provider Provider, subject string) (*OAuthIdentity, error) {
	provider = Provider(normalizeValue(string(provider)))
	if provider == "" {
		return nil, ErrOAuthProviderInvalid
	}

	foundIdentity, err := s.oauthIdentityFinder.FindByProviderSubject(ctx, string(provider), subject)
	if err != nil {
		return nil, fmt.Errorf("find oauth identity: %w", err)
	}
	return foundIdentity, nil
}

func (s *Service) TouchOAuthIdentityLastLogin(ctx context.Context, provider Provider, subject string) error {
	provider = Provider(normalizeValue(string(provider)))
	if provider == "" {
		return ErrOAuthProviderInvalid
	}
	if err := s.oauthIdentityWriter.UpdateLastLogin(ctx, string(provider), subject, s.clock.Now()); err != nil {
		return fmt.Errorf("update oauth last login: %w", err)
	}
	return nil
}

type CreateOAuthIdentityInput struct {
	UserID   string
	Provider Provider
	Subject  string
	Email    string
}

func (s *Service) CreateOAuthIdentity(ctx context.Context, input CreateOAuthIdentityInput) (*OAuthIdentity, error) {
	provider := Provider(normalizeValue(string(input.Provider)))
	if provider == "" {
		return nil, ErrOAuthProviderInvalid
	}

	identityID, err := s.identityIDGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate oauth identity id: %w", err)
	}
	now := s.clock.Now()
	newIdentity := OAuthIdentity{
		ID:          identityID,
		UserID:      input.UserID,
		Provider:    string(provider),
		Subject:     input.Subject,
		Email:       input.Email,
		CreatedAt:   now,
		LastLoginAt: now,
	}
	if err := s.oauthIdentityWriter.Create(ctx, newIdentity); err != nil {
		return nil, fmt.Errorf("create oauth identity: %w", err)
	}
	return &newIdentity, nil
}

func (s *Service) HashPassword(password string) (string, error) {
	return s.passwordHasher.Hash(password)
}

func (s *Service) Now() time.Time {
	return s.clock.Now()
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
