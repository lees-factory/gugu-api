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
	emailValue := strings.TrimSpace(strings.ToLower(input.Email))
	if emailValue == "" || input.Password == "" {
		return nil, ErrInvalidCredentials
	}

	foundUser, err := s.userFinder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	passwordHash, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := s.clock.Now()
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

	token, err := s.tokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate verification token: %w", err)
	}

	emailVerification := verification.EmailVerification{
		Token:     token,
		UserID:    newUser.ID,
		Email:     emailValue,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.emailVerificationWriter.Create(ctx, emailVerification); err != nil {
		return nil, fmt.Errorf("create verification: %w", err)
	}

	if err := s.emailSender.SendVerification(ctx, emailValue, token); err != nil {
		return nil, fmt.Errorf("send verification email: %w", err)
	}

	return &RegisterEmailResult{
		User:                   *newUser,
		VerificationToken:      token,
		VerificationDispatched: true,
	}, nil
}

func (s *Service) LoginEmail(ctx context.Context, input LoginEmailInput) (*LoginResult, error) {
	foundUser, err := s.userFinder.FindByEmail(ctx, strings.TrimSpace(strings.ToLower(input.Email)))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser == nil {
		return nil, ErrInvalidCredentials
	}
	if err := s.passwordHasher.Verify(foundUser.PasswordHash, input.Password); err != nil {
		return nil, ErrInvalidCredentials
	}
	if !foundUser.EmailVerified {
		return nil, ErrEmailNotVerified
	}

	token, err := s.tokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate session token: %w", err)
	}

	now := s.clock.Now()
	session := Session{
		Token:     token,
		UserID:    foundUser.ID,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.sessionAppender.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &LoginResult{
		User:    *foundUser,
		Session: session,
	}, nil
}

func (s *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) (*VerifyEmailResult, error) {
	foundVerification, err := s.emailVerificationFinder.FindByToken(ctx, input.Token)
	if err != nil {
		return nil, fmt.Errorf("find verification by token: %w", err)
	}
	if foundVerification == nil || foundVerification.UsedAt != nil || foundVerification.ExpiresAt.Before(s.clock.Now()) {
		return nil, ErrVerificationNotFound
	}

	now := s.clock.Now()
	if err := s.userWriter.MarkEmailVerified(ctx, foundVerification.UserID, now); err != nil {
		return nil, fmt.Errorf("mark email verified: %w", err)
	}
	if err := s.emailVerificationWriter.MarkUsed(ctx, input.Token, now); err != nil {
		return nil, fmt.Errorf("mark verification used: %w", err)
	}

	foundUser, err := s.userFinder.FindByID(ctx, foundVerification.UserID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	if foundUser == nil {
		return nil, ErrVerificationNotFound
	}

	return &VerifyEmailResult{User: *foundUser}, nil
}

func (s *Service) LoginOAuth(ctx context.Context, input OAuthLoginInput) (*LoginResult, error) {
	provider := Provider(strings.TrimSpace(strings.ToLower(string(input.Provider))))
	if provider == "" {
		return nil, ErrOAuthProviderInvalid
	}

	now := s.clock.Now()
	foundIdentity, err := s.oauthIdentityFinder.FindByProviderSubject(ctx, string(provider), input.Subject)
	if err != nil {
		return nil, fmt.Errorf("find oauth identity: %w", err)
	}

	var foundUser *user.User
	if foundIdentity != nil {
		if err := s.oauthIdentityWriter.UpdateLastLogin(ctx, string(provider), input.Subject, now); err != nil {
			return nil, fmt.Errorf("update oauth last login: %w", err)
		}
		foundUser, err = s.userFinder.FindByID(ctx, foundIdentity.UserID)
		if err != nil {
			return nil, fmt.Errorf("find oauth user by id: %w", err)
		}
	} else {
		emailValue := strings.TrimSpace(strings.ToLower(input.Email))
		foundUser, err = s.userFinder.FindByEmail(ctx, emailValue)
		if err != nil {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
		if foundUser == nil {
			foundUser, err = s.userCreator.Create(ctx, user.CreateInput{
				Email:           emailValue,
				DisplayName:     input.DisplayName,
				AuthSource:      string(provider),
				EmailVerified:   true,
				EmailVerifiedAt: &now,
			})
			if err != nil {
				return nil, fmt.Errorf("create oauth user: %w", err)
			}
		}

		identityID, err := s.identityIDGenerator.New()
		if err != nil {
			return nil, fmt.Errorf("generate oauth identity id: %w", err)
		}
		newIdentity := OAuthIdentity{
			ID:          identityID,
			UserID:      foundUser.ID,
			Provider:    string(provider),
			Subject:     input.Subject,
			Email:       foundUser.Email,
			CreatedAt:   now,
			LastLoginAt: now,
		}
		if err := s.oauthIdentityWriter.Create(ctx, newIdentity); err != nil {
			return nil, fmt.Errorf("create oauth identity: %w", err)
		}
	}

	if foundUser == nil {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate oauth session token: %w", err)
	}

	session := Session{
		Token:     token,
		UserID:    foundUser.ID,
		ExpiresAt: now.Add(30 * 24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.sessionAppender.Create(ctx, session); err != nil {
		return nil, fmt.Errorf("create oauth session: %w", err)
	}

	return &LoginResult{
		User:    *foundUser,
		Session: session,
	}, nil
}
