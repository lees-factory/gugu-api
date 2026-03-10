package auth_test

import (
	"context"
	"testing"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	memoryauth "github.com/ljj/gugu-api/internal/storage/memory/auth"
	memoryuser "github.com/ljj/gugu-api/internal/storage/memory/user"
	memoryverification "github.com/ljj/gugu-api/internal/storage/memory/verification"
	"github.com/ljj/gugu-api/internal/support/email"
	"github.com/ljj/gugu-api/internal/support/id"
	"github.com/ljj/gugu-api/internal/support/security"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

type captureSender struct {
	token string
}

func (s *captureSender) SendVerification(_ context.Context, _ string, token string) error {
	s.token = token
	return nil
}

var _ domainauth.VerificationSender = (*captureSender)(nil)

func TestEmailSignupVerifyAndLogin(t *testing.T) {
	userRepository := memoryuser.NewRepository()
	verificationRepository := memoryverification.NewRepository()
	sessionRepository := memoryauth.NewSessionRepository()
	oauthIdentityRepository := memoryauth.NewOAuthIdentityRepository()
	clock := timeutil.SystemClock{}
	userIDGenerator := id.NewRandomHexGenerator(16)
	identityIDGenerator := id.NewRandomHexGenerator(16)
	tokenGenerator := security.NewRandomTokenGenerator(32)
	sender := &captureSender{}

	service := domainauth.New(
		domainuser.NewFinder(userRepository),
		domainuser.NewCreator(domainuser.NewWriter(userRepository), userIDGenerator, clock),
		domainuser.NewWriter(userRepository),
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		domainauth.NewOAuthIdentityFinder(oauthIdentityRepository),
		domainauth.NewOAuthIdentityWriter(oauthIdentityRepository),
		domainauth.NewSessionAppender(sessionRepository),
		identityIDGenerator,
		tokenGenerator,
		security.BcryptPasswordHasher{},
		clock,
		sender,
	)

	_, err := service.RegisterEmail(context.Background(), domainauth.RegisterEmailInput{
		Email:       "user@example.com",
		Password:    "secret123",
		DisplayName: "User",
	})
	if err != nil {
		t.Fatalf("register email: %v", err)
	}

	if sender.token == "" {
		t.Fatal("expected verification token to be dispatched")
	}

	if _, err := service.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "secret123",
	}); err != domainauth.ErrEmailNotVerified {
		t.Fatalf("expected email not verified, got %v", err)
	}

	if _, err := service.VerifyEmail(context.Background(), domainauth.VerifyEmailInput{Token: sender.token}); err != nil {
		t.Fatalf("verify email: %v", err)
	}

	loginResult, err := service.LoginEmail(context.Background(), domainauth.LoginEmailInput{
		Email:    "user@example.com",
		Password: "secret123",
	})
	if err != nil {
		t.Fatalf("login email: %v", err)
	}

	if loginResult.Session.Token == "" {
		t.Fatal("expected session token")
	}
}

func TestOAuthLoginCreatesAndReusesIdentity(t *testing.T) {
	userRepository := memoryuser.NewRepository()
	verificationRepository := memoryverification.NewRepository()
	sessionRepository := memoryauth.NewSessionRepository()
	oauthIdentityRepository := memoryauth.NewOAuthIdentityRepository()
	clock := timeutil.SystemClock{}
	userIDGenerator := id.NewRandomHexGenerator(16)
	identityIDGenerator := id.NewRandomHexGenerator(16)
	tokenGenerator := security.NewRandomTokenGenerator(32)

	service := domainauth.New(
		domainuser.NewFinder(userRepository),
		domainuser.NewCreator(domainuser.NewWriter(userRepository), userIDGenerator, clock),
		domainuser.NewWriter(userRepository),
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		domainauth.NewOAuthIdentityFinder(oauthIdentityRepository),
		domainauth.NewOAuthIdentityWriter(oauthIdentityRepository),
		domainauth.NewSessionAppender(sessionRepository),
		identityIDGenerator,
		tokenGenerator,
		security.BcryptPasswordHasher{},
		clock,
		email.LogSender{},
	)

	first, err := service.LoginOAuth(context.Background(), domainauth.OAuthLoginInput{
		Provider:    "google",
		Subject:     "google-subject",
		Email:       "oauth@example.com",
		DisplayName: "OAuth User",
	})
	if err != nil {
		t.Fatalf("first oauth login: %v", err)
	}

	second, err := service.LoginOAuth(context.Background(), domainauth.OAuthLoginInput{
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
}
