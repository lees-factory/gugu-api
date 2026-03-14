package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	"github.com/ljj/gugu-api/internal/storage/dbcore"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserSQLCRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *UserSQLCRepository {
	return &UserSQLCRepository{db: db}
}

func (r *UserSQLCRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	const query = `
SELECT
	id,
	email,
	display_name,
	password_hash,
	auth_source,
	email_verified,
	email_verified_at,
	created_at
	FROM gugu.app_users
WHERE email = $1
`

	return r.findOne(ctx, query, email)
}

func (r *UserSQLCRepository) FindByID(ctx context.Context, userID string) (*domainuser.User, error) {
	const query = `
SELECT
	id,
	email,
	display_name,
	password_hash,
	auth_source,
	email_verified,
	email_verified_at,
	created_at
	FROM gugu.app_users
WHERE id = $1
`

	return r.findOne(ctx, query, userID)
}

func (r *UserSQLCRepository) Create(ctx context.Context, newUser domainuser.User) error {
	const query = `
	INSERT INTO gugu.app_users (
	id,
	email,
	display_name,
	password_hash,
	auth_source,
	email_verified,
	email_verified_at,
	created_at
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8
)
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		newUser.ID,
		newUser.Email,
		newUser.DisplayName,
		newUser.PasswordHash,
		newUser.AuthSource,
		newUser.EmailVerified,
		newUser.EmailVerifiedAt,
		newUser.CreatedAt,
	)
	if err != nil {
		if dbcore.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserSQLCRepository) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error {
	const query = `
	UPDATE gugu.app_users
SET
	email_verified = TRUE,
	email_verified_at = $2
WHERE id = $1
`

	result, err := r.db.ExecContext(ctx, query, userID, verifiedAt)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *UserSQLCRepository) findOne(ctx context.Context, query string, arg string) (*domainuser.User, error) {
	var user domainuser.User
	var emailVerifiedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PasswordHash,
		&user.AuthSource,
		&user.EmailVerified,
		&emailVerifiedAt,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if emailVerifiedAt.Valid {
		user.EmailVerifiedAt = &emailVerifiedAt.Time
	}

	return &user, nil
}
