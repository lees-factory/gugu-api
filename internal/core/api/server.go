package api

import (
	"encoding/json"
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	apiauth "github.com/ljj/gugu-api/internal/core/api/v1/auth"
	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
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

func NewServer(_ config.Config) (*Server, error) {
	clock := timeutil.SystemClock{}
	userRepository := memoryuser.NewRepository()
	verificationRepository := memoryverification.NewRepository()
	oauthIdentityRepository := memoryauth.NewOAuthIdentityRepository()
	sessionRepository := memoryauth.NewSessionRepository()
	userIDGenerator := id.NewRandomHexGenerator(16)
	identityIDGenerator := id.NewRandomHexGenerator(16)
	tokenGenerator := security.NewRandomTokenGenerator(32)
	userWriter := domainuser.NewWriter(userRepository)
	userFinder := domainuser.NewFinder(userRepository)
	userCreator := domainuser.NewCreator(userWriter, userIDGenerator, clock)

	authService := domainauth.New(
		domainauth.NewUserAccountReader(userFinder),
		domainauth.NewEmailUserCreator(userFinder, userCreator),
		domainauth.NewOAuthUserResolver(userFinder, userCreator),
		domainauth.NewUserEmailVerifier(userWriter),
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

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	authController := apiauth.NewController(authService)
	router.Get("/health", func(w stdhttp.ResponseWriter, _ *stdhttp.Request) {
		writeJSON(w, stdhttp.StatusOK, map[string]string{"status": "ok"})
	})
	router.Route("/v1/auth", func(r chi.Router) {
		r.Post("/register/email", authController.RegisterEmail)
		r.Post("/login/email", authController.LoginEmail)
		r.Post("/verify-email", authController.VerifyEmail)
		r.Post("/oauth/login", authController.LoginOAuth)
	})

	return &Server{router: router}, nil
}

func (s *Server) Handler() stdhttp.Handler {
	return s.router
}

func writeJSON(w stdhttp.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
