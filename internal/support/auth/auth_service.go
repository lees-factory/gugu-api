package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

type Service struct {
	userFinder             domainuser.Finder
	userCreator            domainuser.Creator
	oauthIdentityFinder    OAuthIdentityFinder
	oauthIdentityWriter    OAuthIdentityWriter
	loginSessionReader     LoginSessionReader
	loginSessionWriter     LoginSessionWriter
	identityIDGenerator    IDGenerator
	sessionIDGenerator     IDGenerator
	tokenFamilyIDGenerator IDGenerator
	refreshTokenGenerator  TokenGenerator
	accessTokenIssuer      AccessTokenIssuer
	passwordVerifier       PasswordVerifier
	clock                  Clock
	refreshTokenHasher     RefreshTokenHasher
}

func NewService(
	userFinder domainuser.Finder,
	userCreator domainuser.Creator,
	oauthIdentityFinder OAuthIdentityFinder,
	oauthIdentityWriter OAuthIdentityWriter,
	loginSessionReader LoginSessionReader,
	loginSessionWriter LoginSessionWriter,
	identityIDGenerator IDGenerator,
	sessionIDGenerator IDGenerator,
	tokenFamilyIDGenerator IDGenerator,
	refreshTokenGenerator TokenGenerator,
	accessTokenIssuer AccessTokenIssuer,
	passwordVerifier PasswordVerifier,
	clock Clock,
	refreshTokenHasher RefreshTokenHasher,
) *Service {
	return &Service{
		userFinder:             userFinder,
		userCreator:            userCreator,
		oauthIdentityFinder:    oauthIdentityFinder,
		oauthIdentityWriter:    oauthIdentityWriter,
		loginSessionReader:     loginSessionReader,
		loginSessionWriter:     loginSessionWriter,
		identityIDGenerator:    identityIDGenerator,
		sessionIDGenerator:     sessionIDGenerator,
		tokenFamilyIDGenerator: tokenFamilyIDGenerator,
		refreshTokenGenerator:  refreshTokenGenerator,
		accessTokenIssuer:      accessTokenIssuer,
		passwordVerifier:       passwordVerifier,
		clock:                  clock,
		refreshTokenHasher:     refreshTokenHasher,
	}
}

func (s *Service) LoginEmail(ctx context.Context, input LoginEmailInput) (*LoginResult, error) {
	foundUser, err := s.userFinder.FindByEmail(ctx, normalizeValue(input.Email))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser == nil {
		return nil, coreerror.New(coreerror.InvalidCredentials)
	}
	if err := s.verifyPassword(foundUser.PasswordHash, input.Password); err != nil {
		return nil, err
	}
	if !foundUser.EmailVerified {
		return nil, coreerror.New(coreerror.EmailNotVerified)
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

func (s *Service) LoginOAuth(ctx context.Context, input OAuthLoginInput) (*LoginResult, error) {
	foundIdentity, err := s.FindOAuthIdentity(ctx, input.Provider, input.Subject)
	if err != nil {
		return nil, err
	}

	var foundUser *domainuser.User
	if foundIdentity != nil {
		if err := s.TouchOAuthIdentityLastLogin(ctx, input.Provider, input.Subject); err != nil {
			return nil, err
		}
		foundUser, err = s.userFinder.FindByID(ctx, foundIdentity.UserID)
		if err != nil {
			return nil, fmt.Errorf("find user by id: %w", err)
		}
		if foundUser == nil {
			return nil, coreerror.New(coreerror.InvalidCredentials)
		}
	} else {
		emailValue := normalizeValue(input.Email)
		foundUser, err = s.userFinder.FindByEmail(ctx, emailValue)
		if err != nil {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
		if foundUser == nil {
			now := s.clock.Now()
			foundUser, err = s.userCreator.Create(ctx, domainuser.NewUser{
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

func (s *Service) RefreshTokens(ctx context.Context, input RefreshTokensInput) (*AuthTokens, error) {
	refreshTokenHash := s.refreshTokenHasher.Hash(input.RefreshToken)
	session, err := s.loginSessionReader.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return nil, fmt.Errorf("find login session by refresh token hash: %w", err)
	}
	if session == nil {
		return nil, coreerror.New(coreerror.RefreshTokenInvalid)
	}

	now := s.clock.Now()
	if session.RevokedAt != nil || session.ExpiresAt.Before(now) {
		return nil, coreerror.New(coreerror.RefreshTokenInvalid)
	}
	if session.RotatedAt != nil {
		if err := s.loginSessionWriter.MarkReuseDetected(ctx, session.ID, now); err != nil {
			return nil, fmt.Errorf("mark login session reuse detected: %w", err)
		}
		if err := s.loginSessionWriter.RevokeFamily(ctx, session.TokenFamilyID, now); err != nil {
			return nil, fmt.Errorf("revoke login session family: %w", err)
		}
		return nil, coreerror.New(coreerror.RefreshTokenInvalid)
	}

	if err := s.loginSessionWriter.UpdateLastSeen(ctx, session.ID, now); err != nil {
		return nil, fmt.Errorf("update login session last seen: %w", err)
	}

	return s.rotateLoginSession(ctx, *session, SessionMetadata{
		UserAgent:  input.UserAgent,
		ClientIP:   input.ClientIP,
		DeviceName: input.DeviceName,
	}, now)
}

func (s *Service) Logout(ctx context.Context, input LogoutInput) error {
	refreshTokenHash := s.refreshTokenHasher.Hash(input.RefreshToken)
	session, err := s.loginSessionReader.FindByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		return fmt.Errorf("find login session by refresh token hash: %w", err)
	}
	if session == nil {
		return coreerror.New(coreerror.RefreshTokenInvalid)
	}

	if err := s.loginSessionWriter.Revoke(ctx, session.ID, s.clock.Now()); err != nil {
		return fmt.Errorf("revoke login session: %w", err)
	}

	return nil
}

func (s *Service) FindOAuthIdentity(ctx context.Context, provider OAuthProvider, subject string) (*OAuthIdentity, error) {
	provider = OAuthProvider(normalizeValue(string(provider)))
	if provider == "" {
		return nil, coreerror.New(coreerror.OAuthProviderInvalid)
	}

	foundIdentity, err := s.oauthIdentityFinder.FindByProviderSubject(ctx, string(provider), subject)
	if err != nil {
		return nil, fmt.Errorf("find oauth identity: %w", err)
	}
	return foundIdentity, nil
}

func (s *Service) TouchOAuthIdentityLastLogin(ctx context.Context, provider OAuthProvider, subject string) error {
	provider = OAuthProvider(normalizeValue(string(provider)))
	if provider == "" {
		return coreerror.New(coreerror.OAuthProviderInvalid)
	}
	if err := s.oauthIdentityWriter.UpdateLastLogin(ctx, string(provider), subject, s.clock.Now()); err != nil {
		return fmt.Errorf("update oauth last login: %w", err)
	}
	return nil
}

func (s *Service) CreateOAuthIdentity(ctx context.Context, input CreateOAuthIdentityInput) (*OAuthIdentity, error) {
	provider := OAuthProvider(normalizeValue(string(input.Provider)))
	if provider == "" {
		return nil, coreerror.New(coreerror.OAuthProviderInvalid)
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

func (s *Service) IssueAuthTokens(userID string) (*AuthTokens, error) {
	return s.issueLoginTokens(context.Background(), userID, SessionMetadata{})
}

func (s *Service) Now() time.Time {
	return s.clock.Now()
}

func (s *Service) verifyPassword(passwordHash string, password string) error {
	if err := s.passwordVerifier.Verify(passwordHash, password); err != nil {
		return coreerror.New(coreerror.InvalidCredentials)
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
