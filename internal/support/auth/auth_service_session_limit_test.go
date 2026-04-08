package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	coreerror "github.com/ljj/gugu-api/internal/core/error"
)

func TestLoginEmail_EnforcesMaxActiveSessions(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user := domainuser.User{
		ID:            "user-1",
		Email:         "user@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	sessionRepo := newTestLoginSessionRepository()
	service := newTestAuthService(t, now, user, sessionRepo)

	ctx := context.Background()
	refreshTokens := make([]string, 0, 7)
	for i := 0; i < 7; i++ {
		loginResult, err := service.LoginEmail(ctx, LoginEmailInput{
			Email:    user.Email,
			Password: "pw-raw",
		})
		if err != nil {
			t.Fatalf("login #%d failed: %v", i+1, err)
		}
		refreshTokens = append(refreshTokens, loginResult.Tokens.RefreshToken)
	}

	activeCount, err := sessionRepo.CountActiveByUserID(ctx, user.ID, now)
	if err != nil {
		t.Fatalf("count active sessions: %v", err)
	}
	if activeCount != 5 {
		t.Fatalf("active session count = %d, want 5", activeCount)
	}

	revokedCount := 0
	for _, session := range sessionRepo.sessions {
		if session.RevokedAt != nil {
			revokedCount++
		}
	}
	if revokedCount != 2 {
		t.Fatalf("revoked session count = %d, want 2", revokedCount)
	}

	if _, err := service.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: refreshTokens[6]}); err != nil {
		t.Fatalf("latest token refresh should succeed: %v", err)
	}
}

func TestLoginEmail_RefreshAndLogoutRegression(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user := domainuser.User{
		ID:            "user-1",
		Email:         "user@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	service := newTestAuthService(t, now, user, newTestLoginSessionRepository())
	ctx := context.Background()

	loginResult, err := service.LoginEmail(ctx, LoginEmailInput{
		Email:    user.Email,
		Password: "pw-raw",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	refreshed, err := service.RefreshTokens(ctx, RefreshTokensInput{
		RefreshToken: loginResult.Tokens.RefreshToken,
	})
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	if err := service.Logout(ctx, LogoutInput{RefreshToken: refreshed.RefreshToken}); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err = service.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: refreshed.RefreshToken})
	if !errors.Is(err, coreerror.New(coreerror.RefreshTokenInvalid)) {
		t.Fatalf("expected invalid refresh token after logout, got: %v", err)
	}
}

func TestRefreshTokens_UpdatesLastSeenAndRotatesSession(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user := domainuser.User{
		ID:            "user-1",
		Email:         "user@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	sessionRepo := newTestLoginSessionRepository()
	service := newTestAuthService(t, now, user, sessionRepo)
	ctx := context.Background()

	loginResult, err := service.LoginEmail(ctx, LoginEmailInput{
		Email:    user.Email,
		Password: "pw-raw",
	})
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	originalHash := testRefreshTokenHasher{}.Hash(loginResult.Tokens.RefreshToken)
	originalSession, err := sessionRepo.FindByRefreshTokenHash(ctx, originalHash)
	if err != nil || originalSession == nil {
		t.Fatalf("original session lookup failed: %v", err)
	}

	refreshed, err := service.RefreshTokens(ctx, RefreshTokensInput{
		RefreshToken: loginResult.Tokens.RefreshToken,
	})
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}

	updatedOriginal, err := sessionRepo.FindByRefreshTokenHash(ctx, originalHash)
	if err != nil || updatedOriginal == nil {
		t.Fatalf("updated original session lookup failed: %v", err)
	}
	if updatedOriginal.RotatedAt == nil {
		t.Fatalf("expected original session to be marked rotated")
	}
	if sessionRepo.updateLastSeenCalled != 1 {
		t.Fatalf("updateLastSeenCalled = %d, want 1", sessionRepo.updateLastSeenCalled)
	}
	if sessionRepo.lastSeenUpdatedSessionID != originalSession.ID {
		t.Fatalf("lastSeen updated session id = %s, want %s", sessionRepo.lastSeenUpdatedSessionID, originalSession.ID)
	}

	newHash := testRefreshTokenHasher{}.Hash(refreshed.RefreshToken)
	newSession, err := sessionRepo.FindByRefreshTokenHash(ctx, newHash)
	if err != nil || newSession == nil {
		t.Fatalf("new session lookup failed: %v", err)
	}
	if newSession.ParentSessionID == nil || *newSession.ParentSessionID != originalSession.ID {
		t.Fatalf("new session parent id = %v, want %s", newSession.ParentSessionID, originalSession.ID)
	}
	if newSession.TokenFamilyID != originalSession.TokenFamilyID {
		t.Fatalf("new session family = %s, want %s", newSession.TokenFamilyID, originalSession.TokenFamilyID)
	}
}

func TestLogout_RevokesOnlyRequestedSession(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user := domainuser.User{
		ID:            "user-1",
		Email:         "user@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	sessionRepo := newTestLoginSessionRepository()
	service := newTestAuthService(t, now, user, sessionRepo)
	ctx := context.Background()

	login1, err := service.LoginEmail(ctx, LoginEmailInput{Email: user.Email, Password: "pw-raw"})
	if err != nil {
		t.Fatalf("first login failed: %v", err)
	}
	login2, err := service.LoginEmail(ctx, LoginEmailInput{Email: user.Email, Password: "pw-raw"})
	if err != nil {
		t.Fatalf("second login failed: %v", err)
	}

	if err := service.Logout(ctx, LogoutInput{RefreshToken: login1.Tokens.RefreshToken}); err != nil {
		t.Fatalf("logout failed: %v", err)
	}
	if sessionRepo.revokeCalled != 1 {
		t.Fatalf("revokeCalled = %d, want 1", sessionRepo.revokeCalled)
	}

	if _, err := service.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: login1.Tokens.RefreshToken}); !errors.Is(err, coreerror.New(coreerror.RefreshTokenInvalid)) {
		t.Fatalf("first token should be invalid after logout, got: %v", err)
	}

	if _, err := service.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: login2.Tokens.RefreshToken}); err != nil {
		t.Fatalf("second token should remain valid, got: %v", err)
	}
}

func TestListMyActiveSessions_ReturnsOnlyActiveSessions(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user := domainuser.User{
		ID:            "user-1",
		Email:         "user@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	sessionRepo := newTestLoginSessionRepository()
	service := newTestAuthService(t, now, user, sessionRepo)
	ctx := context.Background()

	login1, err := service.LoginEmail(ctx, LoginEmailInput{Email: user.Email, Password: "pw-raw"})
	if err != nil {
		t.Fatalf("first login failed: %v", err)
	}
	login2, err := service.LoginEmail(ctx, LoginEmailInput{Email: user.Email, Password: "pw-raw"})
	if err != nil {
		t.Fatalf("second login failed: %v", err)
	}
	if _, err := service.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: login1.Tokens.RefreshToken}); err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if err := service.RevokeMySession(ctx, user.ID, sessionRepo.sessionIDByHash[testRefreshTokenHasher{}.Hash(login2.Tokens.RefreshToken)]); err != nil {
		t.Fatalf("revoke my session failed: %v", err)
	}

	sessions, err := service.ListMyActiveSessions(ctx, user.ID)
	if err != nil {
		t.Fatalf("list my active sessions failed: %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("active sessions len = %d, want 1", len(sessions))
	}
}

func TestRevokeMySession_DoesNotRevokeOtherUsersSession(t *testing.T) {
	now := time.Date(2026, 4, 8, 12, 0, 0, 0, time.UTC)
	user1 := domainuser.User{
		ID:            "user-1",
		Email:         "user1@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	user2 := domainuser.User{
		ID:            "user-2",
		Email:         "user2@example.com",
		PasswordHash:  "pw-hash",
		EmailVerified: true,
	}
	sessionRepo := newTestLoginSessionRepository()
	service1 := newTestAuthService(t, now, user1, sessionRepo)
	service2 := newTestAuthService(t, now, user2, sessionRepo)
	ctx := context.Background()

	login2, err := service2.LoginEmail(ctx, LoginEmailInput{Email: user2.Email, Password: "pw-raw"})
	if err != nil {
		t.Fatalf("user2 login failed: %v", err)
	}
	user2SessionID := sessionRepo.sessionIDByHash[testRefreshTokenHasher{}.Hash(login2.Tokens.RefreshToken)]

	if err := service1.RevokeMySession(ctx, user1.ID, user2SessionID); err != nil {
		t.Fatalf("user1 revoke attempt failed: %v", err)
	}

	if _, err := service2.RefreshTokens(ctx, RefreshTokensInput{RefreshToken: login2.Tokens.RefreshToken}); err != nil {
		t.Fatalf("user2 session should remain valid, got: %v", err)
	}
}

func newTestAuthService(t *testing.T, now time.Time, user domainuser.User, sessionRepo LoginSessionRepository) *Service {
	t.Helper()

	return NewService(
		testUserFinder{user: user},
		testUserCreator{},
		testOAuthIdentityRepository{},
		testOAuthIdentityRepository{},
		NewLoginSessionReader(sessionRepo),
		NewLoginSessionWriter(sessionRepo),
		&sequenceGenerator{prefix: "identity"},
		&sequenceGenerator{prefix: "session"},
		&sequenceGenerator{prefix: "family"},
		&sequenceGenerator{prefix: "refresh"},
		&testAccessTokenIssuer{},
		testPasswordVerifier{},
		testClock{now: now},
		testRefreshTokenHasher{},
	)
}

type testUserFinder struct {
	user domainuser.User
}

func (f testUserFinder) FindByEmail(_ context.Context, email string) (*domainuser.User, error) {
	if email != f.user.Email {
		return nil, nil
	}
	u := f.user
	return &u, nil
}

func (f testUserFinder) FindByID(_ context.Context, userID string) (*domainuser.User, error) {
	if userID != f.user.ID {
		return nil, nil
	}
	u := f.user
	return &u, nil
}

type testUserCreator struct{}

func (testUserCreator) Create(_ context.Context, _ domainuser.NewUser) (*domainuser.User, error) {
	return nil, errors.New("not implemented")
}

type testOAuthIdentityRepository struct{}

func (testOAuthIdentityRepository) FindByProviderSubject(_ context.Context, _, _ string) (*OAuthIdentity, error) {
	return nil, nil
}

func (testOAuthIdentityRepository) Create(_ context.Context, _ OAuthIdentity) error {
	return nil
}

func (testOAuthIdentityRepository) UpdateLastLogin(_ context.Context, _, _ string, _ time.Time) error {
	return nil
}

type sequenceGenerator struct {
	prefix string
	next   int
}

func (g *sequenceGenerator) New() (string, error) {
	g.next++
	return fmt.Sprintf("%s-%d", g.prefix, g.next), nil
}

type testAccessTokenIssuer struct {
	next int
}

func (i *testAccessTokenIssuer) IssueAccessToken(userID string, now time.Time) (IssuedAccessToken, error) {
	i.next++
	return IssuedAccessToken{
		Token:     fmt.Sprintf("access-%s-%d", userID, i.next),
		ExpiresAt: now.Add(15 * time.Minute),
	}, nil
}

type testPasswordVerifier struct{}

func (testPasswordVerifier) Verify(hashedPassword string, rawPassword string) error {
	if hashedPassword == "pw-hash" && rawPassword == "pw-raw" {
		return nil
	}
	return errors.New("invalid password")
}

type testClock struct {
	now time.Time
}

func (c testClock) Now() time.Time {
	return c.now
}

type testRefreshTokenHasher struct{}

func (testRefreshTokenHasher) Hash(value string) string {
	return "hash:" + value
}

type testLoginSessionRepository struct {
	sessions                 map[string]LoginSession
	sessionIDByHash          map[string]string
	sessionIDsByFamily       map[string]map[string]struct{}
	updateLastSeenCalled     int
	lastSeenUpdatedSessionID string
	revokeCalled             int
}

func newTestLoginSessionRepository() *testLoginSessionRepository {
	return &testLoginSessionRepository{
		sessions:           make(map[string]LoginSession),
		sessionIDByHash:    make(map[string]string),
		sessionIDsByFamily: make(map[string]map[string]struct{}),
	}
}

func (r *testLoginSessionRepository) Create(_ context.Context, session LoginSession) error {
	r.sessions[session.ID] = session
	r.sessionIDByHash[session.RefreshTokenHash] = session.ID
	if _, ok := r.sessionIDsByFamily[session.TokenFamilyID]; !ok {
		r.sessionIDsByFamily[session.TokenFamilyID] = make(map[string]struct{})
	}
	r.sessionIDsByFamily[session.TokenFamilyID][session.ID] = struct{}{}
	return nil
}

func (r *testLoginSessionRepository) FindByRefreshTokenHash(_ context.Context, refreshTokenHash string) (*LoginSession, error) {
	sessionID, ok := r.sessionIDByHash[refreshTokenHash]
	if !ok {
		return nil, nil
	}
	session := r.sessions[sessionID]
	return &session, nil
}

func (r *testLoginSessionRepository) ListActiveByUserID(_ context.Context, userID string, now time.Time) ([]LoginSession, error) {
	items := make([]LoginSession, 0)
	for _, session := range r.sessions {
		if session.UserID != userID {
			continue
		}
		if session.RevokedAt != nil || session.RotatedAt != nil || !session.ExpiresAt.After(now) {
			continue
		}
		items = append(items, session)
	}
	return items, nil
}

func (r *testLoginSessionRepository) CountActiveByUserID(_ context.Context, userID string, now time.Time) (int, error) {
	count := 0
	for _, session := range r.sessions {
		if session.UserID != userID {
			continue
		}
		if session.RevokedAt != nil || session.RotatedAt != nil || !session.ExpiresAt.After(now) {
			continue
		}
		count++
	}
	return count, nil
}

func (r *testLoginSessionRepository) MarkRotated(_ context.Context, sessionID string, rotatedAt time.Time) error {
	session := r.sessions[sessionID]
	session.RotatedAt = &rotatedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *testLoginSessionRepository) Revoke(_ context.Context, sessionID string, revokedAt time.Time) error {
	r.revokeCalled++
	session := r.sessions[sessionID]
	session.RevokedAt = &revokedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *testLoginSessionRepository) RevokeByUserSessionID(_ context.Context, userID string, sessionID string, revokedAt time.Time) error {
	session, ok := r.sessions[sessionID]
	if !ok {
		return nil
	}
	if session.UserID != userID || session.RevokedAt != nil {
		return nil
	}
	session.RevokedAt = &revokedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *testLoginSessionRepository) RevokeOldestActiveByUserID(_ context.Context, userID string, now time.Time, revokedAt time.Time) error {
	active := make([]LoginSession, 0)
	for _, session := range r.sessions {
		if session.UserID != userID {
			continue
		}
		if session.RevokedAt != nil || session.RotatedAt != nil || !session.ExpiresAt.After(now) {
			continue
		}
		active = append(active, session)
	}
	if len(active) == 0 {
		return sql.ErrNoRows
	}

	sort.Slice(active, func(i, j int) bool {
		return active[i].CreatedAt.Before(active[j].CreatedAt)
	})
	session := r.sessions[active[0].ID]
	session.RevokedAt = &revokedAt
	r.sessions[active[0].ID] = session
	return nil
}

func (r *testLoginSessionRepository) RevokeFamily(_ context.Context, tokenFamilyID string, revokedAt time.Time) error {
	for sessionID := range r.sessionIDsByFamily[tokenFamilyID] {
		session := r.sessions[sessionID]
		session.RevokedAt = &revokedAt
		r.sessions[sessionID] = session
	}
	return nil
}

func (r *testLoginSessionRepository) MarkReuseDetected(_ context.Context, sessionID string, detectedAt time.Time) error {
	session := r.sessions[sessionID]
	session.ReuseDetectedAt = &detectedAt
	r.sessions[sessionID] = session
	return nil
}

func (r *testLoginSessionRepository) UpdateLastSeen(_ context.Context, sessionID string, lastSeenAt time.Time) error {
	r.updateLastSeenCalled++
	r.lastSeenUpdatedSessionID = sessionID
	session := r.sessions[sessionID]
	session.LastSeenAt = lastSeenAt
	r.sessions[sessionID] = session
	return nil
}
