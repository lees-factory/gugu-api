package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	memoryauth "github.com/ljj/gugu-api/internal/storage/memory/auth"
	memoryuser "github.com/ljj/gugu-api/internal/storage/memory/user"
	memoryverification "github.com/ljj/gugu-api/internal/storage/memory/verification"
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
	token string
	calls int
}

func (s *captureSender) SendVerification(_ context.Context, email string, token string) error {
	s.email = email
	s.token = token
	s.calls++
	return nil
}

var _ domainauth.IDGenerator = (*sequenceGenerator)(nil)
var _ domainauth.TokenGenerator = (*sequenceGenerator)(nil)
var _ domainauth.PasswordHasher = fakePasswordHasher{}
var _ domainauth.VerificationSender = (*captureSender)(nil)
var _ domainauth.Clock = (*fixedClock)(nil)

type authTestFixture struct {
	service                *domainauth.Service
	clock                  *fixedClock
	userRepository         *memoryuser.UserMemoryRepository
	verificationRepository *memoryverification.EmailVerificationMemoryRepository
	oauthRepository        *memoryauth.OAuthIdentityMemoryRepository
	sender                 *captureSender
}

func newAuthTestFixture() *authTestFixture {
	userRepository := memoryuser.NewRepository()
	verificationRepository := memoryverification.NewRepository()
	sessionRepository := memoryauth.NewSessionRepository()
	oauthRepository := memoryauth.NewOAuthIdentityRepository()
	clock := &fixedClock{now: time.Date(2026, time.March, 10, 9, 30, 0, 0, time.UTC)}
	sender := &captureSender{}

	userWriter := domainuser.NewWriter(userRepository)
	userFinder := domainuser.NewFinder(userRepository)
	userCreator := domainuser.NewCreator(
		userWriter,
		&sequenceGenerator{values: []string{"user-1", "user-2", "user-3"}},
		clock,
	)
	service := domainauth.New(
		domainauth.NewUserAccountReader(userFinder),
		domainauth.NewEmailUserCreator(userFinder, userCreator),
		domainauth.NewOAuthUserResolver(userFinder, userCreator),
		domainauth.NewUserEmailVerifier(userWriter),
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		domainauth.NewOAuthIdentityFinder(oauthRepository),
		domainauth.NewOAuthIdentityWriter(oauthRepository),
		domainauth.NewSessionAppender(sessionRepository),
		&sequenceGenerator{values: []string{"identity-1", "identity-2"}},
		&sequenceGenerator{values: []string{"token-1", "token-2", "token-3", "token-4", "token-5"}},
		fakePasswordHasher{},
		clock,
		sender,
	)

	return &authTestFixture{
		service:                service,
		clock:                  clock,
		userRepository:         userRepository,
		verificationRepository: verificationRepository,
		oauthRepository:        oauthRepository,
		sender:                 sender,
	}
}

func TestRegisterEmailCreatesVerificationAndDispatchesToken(t *testing.T) {
	fixture := newAuthTestFixture()

	result, err := fixture.service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
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
	if result.VerificationToken != "token-1" {
		t.Fatalf("expected verification token token-1, got %q", result.VerificationToken)
	}
	if !result.VerificationDispatched {
		t.Fatal("expected verification to be dispatched")
	}
	if fixture.sender.calls != 1 {
		t.Fatalf("expected one verification email, got %d", fixture.sender.calls)
	}
	if fixture.sender.email != "user@example.com" || fixture.sender.token != "token-1" {
		t.Fatalf("unexpected dispatched verification payload: email=%q token=%q", fixture.sender.email, fixture.sender.token)
	}

	foundVerification, err := fixture.verificationRepository.FindByToken(context.Background(), "token-1")
	if err != nil {
		t.Fatalf("find verification: %v", err)
	}
	if foundVerification == nil {
		t.Fatal("expected verification to be stored")
	}
	if foundVerification.ExpiresAt != fixture.clock.now.Add(24*time.Hour) {
		t.Fatalf("expected verification expiry %v, got %v", fixture.clock.now.Add(24*time.Hour), foundVerification.ExpiresAt)
	}
}

func TestRegisterEmailRejectsInvalidCredentials(t *testing.T) {
	testCases := []struct {
		name  string
		input domainauth.RegisterEmailInput
	}{
		{
			name: "empty email",
			input: domainauth.RegisterEmailInput{
				Email:       "   ",
				Password:    "secret123",
				DisplayName: "User",
			},
		},
		{
			name: "empty password",
			input: domainauth.RegisterEmailInput{
				Email:       "user@example.com",
				Password:    "",
				DisplayName: "User",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthTestFixture()

			_, err := fixture.service.RegisterEmail(context.Background(), tc.input)
			if !errors.Is(err, domainauth.ErrInvalidCredentials) {
				t.Fatalf("expected invalid credentials, got %v", err)
			}
		})
	}
}

func TestRegisterEmailRejectsDuplicateEmail(t *testing.T) {
	fixture := newAuthTestFixture()

	_, err := fixture.service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register initial email: %v", err)
	}

	_, err = fixture.service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
		Email:       " USER@EXAMPLE.COM ",
		Password:    "secret456",
		DisplayName: "User Two",
	})
	if !errors.Is(err, domainauth.ErrEmailAlreadyExists) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
}

func TestEmailSignupVerifyAndLogin(t *testing.T) {
	fixture := newAuthTestFixture()

	_, err := fixture.service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register email: %v", err)
	}

	if _, err := fixture.service.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "secret123",
	}); !errors.Is(err, domainauth.ErrEmailNotVerified) {
		t.Fatalf("expected email not verified, got %v", err)
	}

	verifyResult, err := fixture.service.VerifyEmail(context.Background(), domainauth.VerifyEmailInput{Token: fixture.sender.token})
	if err != nil {
		t.Fatalf("verify email: %v", err)
	}
	if !verifyResult.User.EmailVerified {
		t.Fatal("expected verified user after email verification")
	}
	if verifyResult.User.EmailVerifiedAt == nil || !verifyResult.User.EmailVerifiedAt.Equal(fixture.clock.now) {
		t.Fatalf("expected verified at %v, got %v", fixture.clock.now, verifyResult.User.EmailVerifiedAt)
	}

	loginResult, err := fixture.service.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("login email: %v", err)
	}

	if loginResult.Session.Token != "token-2" {
		t.Fatalf("expected session token token-2, got %q", loginResult.Session.Token)
	}
	if loginResult.Session.ExpiresAt != fixture.clock.now.Add(30*24*time.Hour) {
		t.Fatalf("expected session expiry %v, got %v", fixture.clock.now.Add(30*24*time.Hour), loginResult.Session.ExpiresAt)
	}
}

func TestLoginEmailRejectsFailureCases(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, fixture *authTestFixture)
		input domainauth.LoginEmailInput
		want  error
	}{
		{
			name: "unknown email",
			setup: func(t *testing.T, fixture *authTestFixture) {
				t.Helper()
			},
			input: domainauth.LoginEmailInput{
				Email:    "missing@example.com",
				Password: "secret123",
			},
			want: domainauth.ErrInvalidCredentials,
		},
		{
			name: "wrong password",
			setup: func(t *testing.T, fixture *authTestFixture) {
				t.Helper()
				registerAndVerifyEmailUser(t, fixture)
			},
			input: domainauth.LoginEmailInput{
				Email:    "user@example.com",
				Password: "wrong-password",
			},
			want: domainauth.ErrInvalidCredentials,
		},
		{
			name: "email not verified",
			setup: func(t *testing.T, fixture *authTestFixture) {
				t.Helper()
				registerEmailUser(t, fixture)
			},
			input: domainauth.LoginEmailInput{
				Email:    "user@example.com",
				Password: "secret123",
			},
			want: domainauth.ErrEmailNotVerified,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthTestFixture()
			tc.setup(t, fixture)

			_, err := fixture.service.LoginEmail(context.Background(), tc.input)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func TestVerifyEmailRejectsFailureCases(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, fixture *authTestFixture) string
		want  error
	}{
		{
			name: "token not found",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				return "missing-token"
			},
			want: domainauth.ErrVerificationNotFound,
		},
		{
			name: "token already used",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				token := registerEmailUser(t, fixture)
				_, err := fixture.service.VerifyEmail(context.Background(), domainauth.VerifyEmailInput{Token: token})
				if err != nil {
					t.Fatalf("verify email setup: %v", err)
				}
				return token
			},
			want: domainauth.ErrVerificationNotFound,
		},
		{
			name: "token expired",
			setup: func(t *testing.T, fixture *authTestFixture) string {
				t.Helper()
				token := registerEmailUser(t, fixture)
				fixture.clock.now = fixture.clock.now.Add(25 * time.Hour)
				return token
			},
			want: domainauth.ErrVerificationNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthTestFixture()
			token := tc.setup(t, fixture)

			_, err := fixture.service.VerifyEmail(context.Background(), domainauth.VerifyEmailInput{Token: token})
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func TestOAuthLoginCreatesAndReusesIdentity(t *testing.T) {
	fixture := newAuthTestFixture()

	first, err := fixture.service.LoginOAuth(context.Background(), domainauth.OAuthLoginInput{
		Provider:    "google",
		Subject:     "google-subject",
		Email:       "oauth@example.com",
		DisplayName: "OAuth User",
	})
	if err != nil {
		t.Fatalf("first oauth login: %v", err)
	}

	second, err := fixture.service.LoginOAuth(context.Background(), domainauth.OAuthLoginInput{
		Provider:    "google",
		Subject:     "google-subject",
		Email:       "oauth@example.com",
		DisplayName: "OAuth User",
	})
	if err != nil {
		t.Fatalf("second oauth login: %v", err)
	}

	if first.User.ID != second.User.ID {
		t.Fatalf("expected same user id, first=%s second=%s", first.User.ID, second.User.ID)
	}
	if first.Session.Token != "token-1" {
		t.Fatalf("expected first oauth session token token-1, got %q", first.Session.Token)
	}
	if second.Session.Token != "token-2" {
		t.Fatalf("expected second oauth session token token-2, got %q", second.Session.Token)
	}

	foundIdentity, err := fixture.oauthRepository.FindByProviderSubject(context.Background(), "google", "google-subject")
	if err != nil {
		t.Fatalf("find oauth identity: %v", err)
	}
	if foundIdentity == nil {
		t.Fatal("expected oauth identity to be stored")
	}
	if foundIdentity.ID != "identity-1" {
		t.Fatalf("expected identity id identity-1, got %q", foundIdentity.ID)
	}
}

func TestOAuthLoginRejectsFailureCases(t *testing.T) {
	testCases := []struct {
		name  string
		setup func(t *testing.T, fixture *authTestFixture)
		input domainauth.OAuthLoginInput
		want  error
	}{
		{
			name: "invalid provider",
			setup: func(t *testing.T, fixture *authTestFixture) {
				t.Helper()
			},
			input: domainauth.OAuthLoginInput{
				Provider:    "   ",
				Subject:     "subject",
				Email:       "oauth@example.com",
				DisplayName: "OAuth User",
			},
			want: domainauth.ErrOAuthProviderInvalid,
		},
		{
			name: "identity exists but user missing",
			setup: func(t *testing.T, fixture *authTestFixture) {
				t.Helper()
				err := fixture.oauthRepository.Create(context.Background(), domainauth.OAuthIdentity{
					ID:          "identity-existing",
					UserID:      "missing-user",
					Provider:    "google",
					Subject:     "subject",
					Email:       "oauth@example.com",
					CreatedAt:   fixture.clock.now,
					LastLoginAt: fixture.clock.now,
				})
				if err != nil {
					t.Fatalf("create oauth identity: %v", err)
				}
			},
			input: domainauth.OAuthLoginInput{
				Provider:    "google",
				Subject:     "subject",
				Email:       "oauth@example.com",
				DisplayName: "OAuth User",
			},
			want: domainauth.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixture := newAuthTestFixture()
			tc.setup(t, fixture)

			_, err := fixture.service.LoginOAuth(context.Background(), tc.input)
			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}

func registerEmailUser(t *testing.T, fixture *authTestFixture) string {
	t.Helper()

	result, err := fixture.service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register email: %v", err)
	}

	if fixture.sender.token == "" {
		t.Fatal("expected verification token to be dispatched")
	}
	if result.VerificationToken != fixture.sender.token {
		t.Fatalf("expected dispatched token %q, got %q", fixture.sender.token, result.VerificationToken)
	}

	return result.VerificationToken
}

func registerAndVerifyEmailUser(t *testing.T, fixture *authTestFixture) {
	t.Helper()

	token := registerEmailUser(t, fixture)
	if _, err := fixture.service.VerifyEmail(context.Background(), domainauth.VerifyEmailInput{Token: token}); err != nil {
		t.Fatalf("verify email: %v", err)
	}
}
