package auth

import (
	"database/sql"
	"fmt"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
	dbcoreauth "github.com/ljj/gugu-api/internal/storage/dbcore/auth"
	dbcoreuser "github.com/ljj/gugu-api/internal/storage/dbcore/user"
	dbcoreverification "github.com/ljj/gugu-api/internal/storage/dbcore/verification"
	memoryauth "github.com/ljj/gugu-api/internal/storage/memory/auth"
	memoryuser "github.com/ljj/gugu-api/internal/storage/memory/user"
	memoryverification "github.com/ljj/gugu-api/internal/storage/memory/verification"
	supportauth "github.com/ljj/gugu-api/internal/support/auth"
	"github.com/ljj/gugu-api/internal/support/config"
	"github.com/ljj/gugu-api/internal/support/id"
	"github.com/ljj/gugu-api/internal/support/mailer"
	"github.com/ljj/gugu-api/internal/support/security"
	timeutil "github.com/ljj/gugu-api/internal/support/time"
)

type Controllers struct {
	Auth *AuthController
	User *UserController
}

func Wire(cfg config.Config, db *sql.DB) (*Controllers, error) {
	clock := timeutil.SystemClock{}
	userRepository, verificationRepository, oauthIdentityRepository := buildAuthRepositories(db)
	loginSessionRepository := buildLoginSessionRepository(db)

	userWriter := domainuser.NewWriter(userRepository)
	userFinder := domainuser.NewFinder(userRepository)
	userCreator := domainuser.NewCreator(userWriter, id.NewRandomHexGenerator(16), clock)
	passwordHasher := security.BcryptPasswordHasher{}

	verificationMailer, err := mailer.NewSender(mailer.Config{
		Provider:     cfg.MailProvider,
		MailFrom:     cfg.MailFrom,
		SMTPHost:     cfg.SMTPHost,
		SMTPPort:     cfg.SMTPPort,
		SMTPUsername: cfg.SMTPUsername,
		SMTPPassword: cfg.SMTPPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("build mailer: %w", err)
	}

	userService := domainuser.NewService(
		userFinder,
		userCreator,
		userWriter,
		domainverification.NewFinder(verificationRepository),
		domainverification.NewWriter(verificationRepository),
		security.NewNumericCodeGenerator(6),
		verificationMailer,
		clock,
	)

	authService := supportauth.NewService(
		userFinder,
		userCreator,
		supportauth.NewOAuthIdentityFinder(oauthIdentityRepository),
		supportauth.NewOAuthIdentityWriter(oauthIdentityRepository),
		supportauth.NewLoginSessionReader(loginSessionRepository),
		supportauth.NewLoginSessionWriter(loginSessionRepository),
		id.NewRandomHexGenerator(16),
		id.NewRandomHexGenerator(16),
		id.NewRandomHexGenerator(16),
		security.NewRandomTokenGenerator(32),
		security.NewJWTTokenIssuer(cfg.JWTSecret, cfg.JWTIssuer),
		passwordHasher,
		clock,
		security.TokenSHA256Hasher{},
	)

	return &Controllers{
		Auth: NewAuthController(authService),
		User: NewUserController(userService, passwordHasher),
	}, nil
}

func buildAuthRepositories(db *sql.DB) (domainuser.Repository, domainverification.Repository, supportauth.OAuthIdentityRepository) {
	if db == nil {
		return memoryuser.NewRepository(), memoryverification.NewRepository(), memoryauth.NewOAuthIdentityRepository()
	}
	return dbcoreuser.NewSQLCRepository(db), dbcoreverification.NewSQLCRepository(db), dbcoreauth.NewOAuthIdentitySQLCRepository(db)
}

func buildLoginSessionRepository(db *sql.DB) supportauth.LoginSessionRepository {
	if db == nil {
		return memoryauth.NewLoginSessionRepository()
	}
	return dbcoreauth.NewLoginSessionSQLCRepository(db)
}
