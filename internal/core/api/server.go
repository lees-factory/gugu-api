package api

import (
	"database/sql"
	"fmt"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apiadvice "github.com/ljj/gugu-api/internal/core/api/advice"
	apiresponse "github.com/ljj/gugu-api/internal/core/api/response"
	apiauth "github.com/ljj/gugu-api/internal/core/api/v1/auth"
	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	dbcoreauth "github.com/ljj/gugu-api/internal/storage/dbcore/auth"
	dbcoreuser "github.com/ljj/gugu-api/internal/storage/dbcore/user"
	dbcoreverification "github.com/ljj/gugu-api/internal/storage/dbcore/verification"
	memoryauth "github.com/ljj/gugu-api/internal/storage/memory/auth"
	memoryuser "github.com/ljj/gugu-api/internal/storage/memory/user"
	memoryverification "github.com/ljj/gugu-api/internal/storage/memory/verification"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/email"
	"github.com/ljj/gugu-api/internal/support/id"
	"github.com/ljj/gugu-api/internal/support/security"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

type Server struct {
	router *chi.Mux
}

func NewServer(cfg config.Config, db *sql.DB) (*Server, error) {
	clock := timeutil.SystemClock{}
	userRepository, verificationRepository, oauthIdentityRepository := buildRepositories(db)
	loginSessionRepository := buildLoginSessionRepository(db)
	userIDGenerator := id.NewRandomHexGenerator(16)
	identityIDGenerator := id.NewRandomHexGenerator(16)
	sessionIDGenerator := id.NewRandomHexGenerator(16)
	tokenFamilyIDGenerator := id.NewRandomHexGenerator(16)
	verificationCodeGenerator := security.NewNumericCodeGenerator(6)
	refreshTokenGenerator := security.NewRandomTokenGenerator(32)
	authTokenIssuer := security.NewJWTTokenIssuer(cfg.JWTSecret, cfg.JWTIssuer)
	userWriter := domainuser.NewWriter(userRepository)
	userFinder := domainuser.NewFinder(userRepository)
	userCreator := domainuser.NewCreator(userWriter, userIDGenerator, clock)
	emailSender, err := email.NewSender(email.Config{
		Provider:     cfg.MailProvider,
		MailFrom:     cfg.MailFrom,
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUsername: cfg.SMTPUsername,
		SMTPPassword: cfg.SMTPPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("build email sender: %w", err)
	}

	authService := domainauth.New(
		userFinder,
		userCreator,
		userWriter,
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		domainauth.NewOAuthIdentityFinder(oauthIdentityRepository),
		domainauth.NewOAuthIdentityWriter(oauthIdentityRepository),
		domainauth.NewLoginSessionReader(loginSessionRepository),
		domainauth.NewLoginSessionWriter(loginSessionRepository),
		identityIDGenerator,
		sessionIDGenerator,
		tokenFamilyIDGenerator,
		verificationCodeGenerator,
		refreshTokenGenerator,
		authTokenIssuer,
		security.BcryptPasswordHasher{},
		clock,
		emailSender,
		security.TokenSHA256Hasher{},
	)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(CORSMiddleware(cfg.CORSAllowedOrigins))

	authController := apiauth.NewController(authService)
	router.Get("/health", apiadvice.Wrap(func(_ *stdhttp.Request) (int, any, error) {
		return stdhttp.StatusOK, apiresponse.SuccessWithData(map[string]string{"status": "ok"}), nil
	}))
	router.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register/email", apiadvice.Wrap(authController.RegisterEmail))
		r.Post("/login/email", apiadvice.Wrap(authController.LoginEmail))
		r.Post("/verify-email", apiadvice.Wrap(authController.VerifyEmail))
		r.Post("/oauth/login", apiadvice.Wrap(authController.LoginOAuth))
		r.Post("/refresh", apiadvice.Wrap(authController.Refresh))
		r.Post("/logout", apiadvice.Wrap(authController.Logout))
	})

	return &Server{router: router}, nil
}

func (s *Server) Handler() stdhttp.Handler {
	return s.router
}

func buildRepositories(db *sql.DB) (domainuser.Repository, domainverification.Repository, domainauth.OAuthIdentityRepository) {
	if db == nil {
		return memoryuser.NewRepository(), memoryverification.NewRepository(), memoryauth.NewOAuthIdentityRepository()
	}

	return dbcoreuser.NewRepository(db), dbcoreverification.NewRepository(db), dbcoreauth.NewOAuthIdentityRepository(db)
}

func buildLoginSessionRepository(db *sql.DB) domainauth.LoginSessionRepository {
	if db == nil {
		return memoryauth.NewLoginSessionRepository()
	}

	return dbcoreauth.NewLoginSessionRepository(db)
}
