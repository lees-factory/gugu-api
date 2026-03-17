package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
	memoryauth "github.com/ljj/gugu-api/internal/storage/memory/auth"
	memoryuser "github.com/ljj/gugu-api/internal/storage/memory/user"
	memoryverification "github.com/ljj/gugu-api/internal/storage/memory/verification"
	"github.com/ljj/gugu-api/internal/support/security"
)

type fixedClock struct {
	now time.Time
}

func (c *fixedClock) Now() time.Time {
	return c.now
}

type sequenceGenerator struct {
	values []string
	index  int
}

func (g *sequenceGenerator) New() (string, error) {
	if g.index >= len(g.values) {
		return "", errors.New("no more generated values")
	}

	value := g.values[g.index]
	g.index++
	return value, nil
}

type fakePasswordHasher struct{}

func (h fakePasswordHasher) Hash(password string) (string, error) {
	return "hashed:" + password, nil
}

func (h fakePasswordHasher) Verify(hashedPassword string, rawPassword string) error {
	if hashedPassword != "hashed:"+rawPassword {
		return errors.New("password mismatch")
	}
	return nil
}

type captureSender struct {
	email string
	code  string
	calls int
}

type fakeAuthTokenIssuer struct {
	accessToken string
}

func (s *captureSender) SendVerification(_ context.Context, email string, code string) error {
	s.email = email
	s.code = code
	s.calls++
	return nil
}

func (i fakeAuthTokenIssuer) IssueAccessToken(userID string, now time.Time) (domainauth.IssuedAccessToken, error) {
	return domainauth.IssuedAccessToken{
		Token:     i.accessToken + ":" + userID,
		ExpiresAt: now.Add(15 * time.Minute),
	}, nil
}

var _ domainauth.IDGenerator = (*sequenceGenerator)(nil)
var _ domainauth.TokenGenerator = (*sequenceGenerator)(nil)
var _ domainauth.PasswordVerifier = fakePasswordHasher{}
var _ domainauth.AccessTokenIssuer = fakeAuthTokenIssuer{}
var _ domainauth.Clock = (*fixedClock)(nil)
var _ domainuser.Clock = (*fixedClock)(nil)
var _ domainuser.IDGenerator = (*sequenceGenerator)(nil)
var _ domainuser.CodeGenerator = (*sequenceGenerator)(nil)
var _ domainuser.PasswordHasher = fakePasswordHasher{}
var _ domainuser.VerificationMailer = (*captureSender)(nil)

type authTestFixture struct {
	authService            *domainauth.Service
	userService            *domainuser.Service
	clock                  *fixedClock
	verificationRepository *memoryverification.EmailVerificationMemoryRepository
	oauthRepository        *memoryauth.OAuthIdentityMemoryRepository
	loginSessionRepository *memoryauth.LoginSessionMemoryRepository
	sender                 *captureSender
}

func newAuthTestFixture() *authTestFixture {
	userRepository := memoryuser.NewRepository()
	verificationRepository := memoryverification.NewRepository()
	oauthRepository := memoryauth.NewOAuthIdentityRepository()
	loginSessionRepository := memoryauth.NewLoginSessionRepository()
	clock := &fixedClock{now: time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)}
	sender := &captureSender{}

	userWriter := domainuser.NewWriter(userRepository)
	userFinder := domainuser.NewFinder(userRepository)
	userCreator := domainuser.NewCreator(
		userWriter,
		&sequenceGenerator{values: []string{"user-1", "user-2", "user-3"}},
		clock,
	)

	userService := domainuser.NewService(
		userFinder,
		userCreator,
		userWriter,
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		&sequenceGenerator{values: []string{"123456", "654321", "111111", "222222", "333333"}},
		fakePasswordHasher{},
		sender,
		clock,
	)

	return &authTestFixture{
		authService: domainauth.New(
			userFinder,
			userCreator,
			domainauth.NewOAuthIdentityFinder(oauthRepository),
			domainauth.NewOAuthIdentityWriter(oauthRepository),
			domainauth.NewLoginSessionReader(loginSessionRepository),
			domainauth.NewLoginSessionWriter(loginSessionRepository),
			&sequenceGenerator{values: []string{"identity-1", "identity-2"}},
			&sequenceGenerator{values: []string{"session-1", "session-2", "session-3", "session-4"}},
			&sequenceGenerator{values: []string{"family-1", "family-2", "family-3"}},
			&sequenceGenerator{values: []string{"refresh-1", "refresh-2", "refresh-3", "refresh-4"}},
			fakeAuthTokenIssuer{accessToken: "access-token"},
			fakePasswordHasher{},
			clock,
			security.TokenSHA256Hasher{},
		),
		userService:            userService,
		clock:                  clock,
		verificationRepository: verificationRepository,
		oauthRepository:        oauthRepository,
		loginSessionRepository: loginSessionRepository,
		sender:                 sender,
	}
}

func TestRegisterEmailCreatesUserAndDispatchesVerification(t *testing.T) {
	fixture := newAuthTestFixture()

	result, err := fixture.userService.RegisterEmail(context.Background(), domainuser.RegisterEmailInput{
		Email:       " User@Example.com ",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register email: %v", err)
	}

	if result.User.Email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", result.User.Email)
	}
	if result.User.AuthSource != "email" {
		t.Fatalf("expected auth source email, got %q", result.User.AuthSource)
	}
	if result.VerificationCode != "123456" {
		t.Fatalf("expected verification code 123456, got %q", result.VerificationCode)
	}
	if !result.VerificationDispatched {
		t.Fatal("expected verification dispatched")
	}
	if fixture.sender.calls != 1 || fixture.sender.email != "user@example.com" {
		t.Fatalf("unexpected verification dispatch: calls=%d email=%q", fixture.sender.calls, fixture.sender.email)
	}

	foundVerification, err := fixture.verificationRepository.FindByCode(context.Background(), "123456")
	if err != nil {
		t.Fatalf("find verification: %v", err)
	}
	if foundVerification == nil {
		t.Fatal("expected verification stored")
	}
	if foundVerification.ExpiresAt != fixture.clock.now.Add(10*time.Minute) {
		t.Fatalf("expected verification expiry %v, got %v", fixture.clock.now.Add(10*time.Minute), foundVerification.ExpiresAt)
	}
}

func TestUserServiceRejectsDuplicateEmail(t *testing.T) {
	fixture := newAuthTestFixture()

	_, err := fixture.userService.RegisterEmail(context.Background(), domainuser.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register initial user: %v", err)
	}

	_, err = fixture.userService.RegisterEmail(context.Background(), domainuser.RegisterEmailInput{
		Email:       " USER@EXAMPLE.COM ",
		Password:    "secret123",
		DisplayName: "User Two",
	})
	if !errors.Is(err, coreerror.New(coreerror.EmailAlreadyExists)) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestVerifyEmailCodeMarksUserVerified(t *testing.T) {
	fixture := newAuthTestFixture()
	code := registerEmailUser(t, fixture)

	verifyResult, err := fixture.userService.VerifyEmail(context.Background(), domainuser.VerifyEmailInput{Code: code})
	if err != nil {
		t.Fatalf("verify email: %v", err)
	}

	if !verifyResult.User.EmailVerified {
		t.Fatal("expected verified user")
	}
	if verifyResult.User.EmailVerifiedAt == nil || !verifyResult.User.EmailVerifiedAt.Equal(fixture.clock.now) {
		t.Fatalf("expected verified at %v, got %v", fixture.clock.now, verifyResult.User.EmailVerifiedAt)
	}
}

func TestVerifyEmailCodeRejectsFailureCases(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, fixture *authTestFixture) string
		want  error
	}{
		{
			name: "code not found",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				return "999999"
			},
			want: coreerror.New(coreerror.VerificationNotFound),
		},
		{
			name: "code already used",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				code := registerEmailUser(t, fixture)
				if _, err := fixture.userService.VerifyEmail(context.Background(), domainuser.VerifyEmailInput{Code: code}); err != nil {
					t.Fatalf("verify email setup: %v", err)
				}
				return code
			},
			want: coreerror.New(coreerror.VerificationNotFound),
		},
		{
			name: "code expired",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				code := registerEmailUser(t, fixture)
				fixture.clock.now = fixture.clock.now.Add(11 * time.Minute)
				return code
			},
			want: coreerror.New(coreerror.VerificationNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthTestFixture()
			code := tc.setup(t, fixture)

			_, err := fixture.userService.VerifyEmail(context.Background(), domainuser.VerifyEmailInput{Code: code})
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func TestEmailLoginFlow(t *testing.T) {
	fixture := newAuthTestFixture()
	registerAndVerifyEmailUser(t, fixture)

	loginResult, err := fixture.authService.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("login email: %v", err)
	}

	if loginResult.Tokens.AccessToken != "access-token:"+loginResult.User.ID {
		t.Fatalf("expected access token for %q, got %q", loginResult.User.ID, loginResult.Tokens.AccessToken)
	}
	if loginResult.Tokens.RefreshToken == "" {
		t.Fatal("expected refresh token to be issued")
	}
	if loginResult.Tokens.AccessExpiresAt != fixture.clock.now.Add(15*time.Minute) {
		t.Fatalf("expected access expiry %v, got %v", fixture.clock.now.Add(15*time.Minute), loginResult.Tokens.AccessExpiresAt)
	}
}

func TestLoginEmailRejectsInvalidPassword(t *testing.T) {
	fixture := newAuthTestFixture()
	registerAndVerifyEmailUser(t, fixture)

	_, err := fixture.authService.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "wrong-password",
	})
	if !errors.Is(err, coreerror.New(coreerror.InvalidCredentials)) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestOAuthUserAndIdentityFlow(t *testing.T) {
	fixture := newAuthTestFixture()

	foundUser, err := fixture.userService.FindOrCreateOAuthUser(context.Background(), domainuser.FindOrCreateOAuthUserInput{
		Email:       "oauth@example.com",
		DisplayName: "OAuth User",
		AuthSource:  "google",
		VerifiedAt:  fixture.authService.Now(),
	})
	if err != nil {
		t.Fatalf("find or create oauth user: %v", err)
	}

	identity, err := fixture.authService.CreateOAuthIdentity(context.Background(), domainauth.CreateOAuthIdentityInput{
		UserID:   foundUser.ID,
		Provider: "google",
		Subject:  "google-subject",
		Email:    foundUser.Email,
	})
	if err != nil {
		t.Fatalf("create oauth identity: %v", err)
	}

	foundIdentity, err := fixture.authService.FindOAuthIdentity(context.Background(), "google", "google-subject")
	if err != nil {
		t.Fatalf("find oauth identity: %v", err)
	}
	if foundIdentity == nil || foundIdentity.ID != identity.ID {
		t.Fatalf("expected same oauth identity, got %+v", foundIdentity)
	}

	if err := fixture.authService.TouchOAuthIdentityLastLogin(context.Background(), "google", "google-subject"); err != nil {
		t.Fatalf("touch oauth last login: %v", err)
	}

	tokens, err := fixture.authService.IssueAuthTokens(foundUser.ID)
	if err != nil {
		t.Fatalf("issue oauth auth tokens: %v", err)
	}
	if tokens.AccessToken != "access-token:"+foundUser.ID {
		t.Fatalf("expected oauth access token for %q, got %q", foundUser.ID, tokens.AccessToken)
	}
}

func TestOAuthRejectsInvalidProvider(t *testing.T) {
	fixture := newAuthTestFixture()

	_, err := fixture.authService.FindOAuthIdentity(context.Background(), "   ", "subject")
	if !errors.Is(err, coreerror.New(coreerror.OAuthProviderInvalid)) {
		t.Fatalf("expected invalid provider error, got %v", err)
	}
}

func registerEmailUser(t *testing.T, fixture *authTestFixture) string {
	t.Helper()

	result, err := fixture.userService.RegisterEmail(context.Background(), domainuser.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register email user: %v", err)
	}

	return result.VerificationCode
}

func registerAndVerifyEmailUser(t *testing.T, fixture *authTestFixture) {
	t.Helper()

	code := registerEmailUser(t, fixture)
	if _, err := fixture.userService.VerifyEmail(context.Background(), domainuser.VerifyEmailInput{Code: code}); err != nil {
		t.Fatalf("verify email: %v", err)
	}
}
