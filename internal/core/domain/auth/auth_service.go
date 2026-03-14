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
	userFinder                user.Finder
	userCreator               user.Creator
	userWriter                user.Writer
	emailVerificationFinder   verification.Finder
	emailVerificationWriter   verification.Writer
	oauthIdentityFinder       OAuthIdentityFinder
	oauthIdentityWriter       OAuthIdentityWriter
	loginSessionReader        LoginSessionReader
	loginSessionWriter        LoginSessionWriter
	identityIDGenerator       IDGenerator
	sessionIDGenerator        IDGenerator
	tokenFamilyIDGenerator    IDGenerator
	verificationCodeGenerator TokenGenerator
	refreshTokenGenerator     TokenGenerator
	accessTokenIssuer         AccessTokenIssuer
	passwordHasher            PasswordHasher
	clock                     Clock
	emailSender               VerificationSender
	refreshTokenHasher        RefreshTokenHasher
}

func New(
	userFinder user.Finder,
	userCreator user.Creator,
	userWriter user.Writer,
	emailVerificationFinder verification.Finder,
	emailVerificationWriter verification.Writer,
	oauthIdentityFinder OAuthIdentityFinder,
	oauthIdentityWriter OAuthIdentityWriter,
	loginSessionReader LoginSessionReader,
	loginSessionWriter LoginSessionWriter,
	identityIDGenerator IDGenerator,
	sessionIDGenerator IDGenerator,
	tokenFamilyIDGenerator IDGenerator,
	verificationCodeGenerator TokenGenerator,
	refreshTokenGenerator TokenGenerator,
	accessTokenIssuer AccessTokenIssuer,
	passwordHasher PasswordHasher,
	clock Clock,
	emailSender VerificationSender,
	refreshTokenHasher RefreshTokenHasher,
) *Service {
	return &Service{
		userFinder:                userFinder,
		userCreator:               userCreator,
		userWriter:                userWriter,
		emailVerificationFinder:   emailVerificationFinder,
		emailVerificationWriter:   emailVerificationWriter,
		oauthIdentityFinder:       oauthIdentityFinder,
		oauthIdentityWriter:       oauthIdentityWriter,
		loginSessionReader:        loginSessionReader,
		loginSessionWriter:        loginSessionWriter,
		identityIDGenerator:       identityIDGenerator,
		sessionIDGenerator:        sessionIDGenerator,
		tokenFamilyIDGenerator:    tokenFamilyIDGenerator,
		verificationCodeGenerator: verificationCodeGenerator,
		refreshTokenGenerator:     refreshTokenGenerator,
		accessTokenIssuer:         accessTokenIssuer,
		passwordHasher:            passwordHasher,
		clock:                     clock,
		emailSender:               emailSender,
		refreshTokenHasher:        refreshTokenHasher,
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
		VerificationCode:       verificationResult.VerificationCode,
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

	tokens, err := s.issueLoginTokens(ctx, foundUser.ID, SessionMetadata{
		UserAgent:  input.UserAgent,
		ClientIP:   input.ClientIP,
		DeviceName: input.DeviceName,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:   *foundUser,
		Tokens: *tokens,
	}, nil
}

func (s *Service) VerifyEmail(ctx context.Context, input VerifyEmailInput) (*VerifyEmailResult, error) {
	verifyResult, err := s.VerifyEmailCode(ctx, VerifyEmailCodeInput{Code: input.Code})
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
			foundUser, err = s.userCreator.Create(ctx, user.CreateInput{
				Email:           emailValue,
				DisplayName:     input.DisplayName,
				AuthSource:      string(input.Provider),
				EmailVerified:   true,
				EmailVerifiedAt: new(s.Now()),
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

	tokens, err := s.issueLoginTokens(ctx, foundUser.ID, SessionMetadata{
		UserAgent:  input.UserAgent,
		ClientIP:   input.ClientIP,
		DeviceName: input.DeviceName,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		User:   *foundUser,
		Tokens: *tokens,
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
	code, err := s.verificationCodeGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate verification code: %w", err)
	}

	emailVerification := verification.EmailVerification{
		Code:      code,
		UserID:    input.UserID,
		Email:     input.Email,
		ExpiresAt: now.Add(10 * time.Minute),
		CreatedAt: now,
	}
	if err := s.emailVerificationWriter.Create(ctx, emailVerification); err != nil {
		return nil, fmt.Errorf("create verification: %w", err)
	}

	if err := s.emailSender.SendVerification(ctx, input.Email, code); err != nil {
		return nil, fmt.Errorf("send verification email: %w", err)
	}

	return &IssueEmailVerificationResult{
		VerificationCode:       code,
		VerificationDispatched: true,
	}, nil
}

func (s *Service) VerifyPassword(passwordHash string, password string) error {
	if err := s.passwordHasher.Verify(passwordHash, password); err != nil {
		return ErrInvalidCredentials
	}
	return nil
}

func (s *Service) VerifyEmailCode(ctx context.Context, input VerifyEmailCodeInput) (*VerifyEmailCodeResult, error) {
	foundVerification, err := s.emailVerificationFinder.FindByCode(ctx, input.Code)
	if err != nil {
		return nil, fmt.Errorf("find verification by code: %w", err)
	}
	if foundVerification == nil || foundVerification.UsedAt != nil || foundVerification.ExpiresAt.Before(s.clock.Now()) {
		return nil, ErrVerificationNotFound
	}

	now := s.clock.Now()
	if err := s.emailVerificationWriter.MarkUsed(ctx, input.Code, now); err != nil {
		return nil, fmt.Errorf("mark verification used: %w", err)
	}

	return &VerifyEmailCodeResult{
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

func (s *Service) IssueAuthTokens(userID string) (*AuthTokens, error) {
	return s.issueLoginTokens(context.Background(), userID, SessionMetadata{})
}

func (s *Service) RefreshTokens(ctx context.Context, input RefreshTokensInput) (*AuthTokens, error) {
	refreshTokenHash := s.refreshTokenHasher.Hash(input.RefreshToken)
	session, err := s.loginSessionReader.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("find login session by refresh token hash: %w", err)
	}
	if session == nil {
		return nil, ErrRefreshTokenInvalid
	}

	now := s.clock.Now()
	if session.RevokedAt != nil || session.ExpiresAt.Before(now) {
		return nil, ErrRefreshTokenInvalid
	}
	if session.RotatedAt != nil {
		if err := s.loginSessionWriter.MarkReuseDetected(ctx, session.ID, now); err != nil {
			return nil, fmt.Errorf("mark login session reuse detected: %w", err)
		}
		if err := s.loginSessionWriter.RevokeFamily(ctx, session.TokenFamilyID, now); err != nil {
			return nil, fmt.Errorf("revoke login session family: %w", err)
		}
		return nil, ErrRefreshTokenInvalid
	}

	if err := s.loginSessionWriter.UpdateLastSeen(ctx, session.ID, now); err != nil {
		return nil, fmt.Errorf("update login session last seen: %w", err)
	}

	newTokens, err := s.rotateLoginSession(ctx, *session, SessionMetadata{
		UserAgent:  input.UserAgent,
		ClientIP:   input.ClientIP,
		DeviceName: input.DeviceName,
	}, now)
	if err != nil {
		return nil, err
	}

	return newTokens, nil
}

func (s *Service) Logout(ctx context.Context, input LogoutInput) error {
	refreshTokenHash := s.refreshTokenHasher.Hash(input.RefreshToken)
	session, err := s.loginSessionReader.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return fmt.Errorf("find login session by refresh token hash: %w", err)
	}
	if session == nil {
		return ErrRefreshTokenInvalid
	}

	if err := s.loginSessionWriter.Revoke(ctx, session.ID, s.clock.Now()); err != nil {
		return fmt.Errorf("revoke login session: %w", err)
	}

	return nil
}

func (s *Service) issueLoginTokens(ctx context.Context, userID string, metadata SessionMetadata) (*AuthTokens, error) {
	now := s.clock.Now()
	tokenFamilyID, err := s.tokenFamilyIDGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate token family id: %w", err)
	}

	return s.createLoginSession(ctx, userID, tokenFamilyID, nil, metadata, now)
}

func (s *Service) rotateLoginSession(ctx context.Context, session LoginSession, metadata SessionMetadata, now time.Time) (*AuthTokens, error) {
	if err := s.loginSessionWriter.MarkRotated(ctx, session.ID, now); err != nil {
		return nil, fmt.Errorf("mark login session rotated: %w", err)
	}

	return s.createLoginSession(ctx, session.UserID, session.TokenFamilyID, &session.ID, metadata, now)
}

func (s *Service) createLoginSession(ctx context.Context, userID string, tokenFamilyID string, parentSessionID *string, metadata SessionMetadata, now time.Time) (*AuthTokens, error) {
	sessionID, err := s.sessionIDGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate login session id: %w", err)
	}
	refreshToken, err := s.refreshTokenGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}
	issuedAccessToken, err := s.accessTokenIssuer.IssueAccessToken(userID, now)
	if err != nil {
		return nil, fmt.Errorf("issue access token for user %s: %w", userID, err)
	}

	refreshExpiresAt := now.Add(14 * 24 * time.Hour)
	loginSession := LoginSession{
		ID:               sessionID,
		UserID:           userID,
		RefreshTokenHash: s.refreshTokenHasher.Hash(refreshToken),
		TokenFamilyID:    tokenFamilyID,
		ParentSessionID:  parentSessionID,
		UserAgent:        metadata.UserAgent,
		ClientIP:         metadata.ClientIP,
		DeviceName:       metadata.DeviceName,
		ExpiresAt:        refreshExpiresAt,
		LastSeenAt:       now,
		CreatedAt:        now,
	}
	if err := s.loginSessionWriter.Create(ctx, loginSession); err != nil {
		return nil, fmt.Errorf("create login session: %w", err)
	}

	return &AuthTokens{
		AccessToken:      issuedAccessToken.Token,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		AccessExpiresAt:  issuedAccessToken.ExpiresAt,
		RefreshExpiresAt: refreshExpiresAt,
	}, nil
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}
