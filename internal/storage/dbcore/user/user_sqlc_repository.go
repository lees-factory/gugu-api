package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainuser "github.com/ljj/gugu-api/internal/core/domain/user"
	"github.com/ljj/gugu-api/internal/storage/dbcore"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserSQLCRepository struct {
	queries *sqldb.Queries
}

func NewSQLCRepository(db *sql.DB) *UserSQLCRepository {
	return &UserSQLCRepository{queries: sqldb.New(db)}
}

func (r *UserSQLCRepository) FindByEmail(ctx context.Context, email string) (*domainuser.User, error) {
	row, err := r.queries.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainUser(row), nil
}

func (r *UserSQLCRepository) FindByID(ctx context.Context, userID string) (*domainuser.User, error) {
	row, err := r.queries.FindUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return toDomainUser(row), nil
}

func (r *UserSQLCRepository) Create(ctx context.Context, newUser domainuser.User) error {
	err := r.queries.CreateUser(ctx, sqldb.CreateUserParams{
		ID:              newUser.ID,
		Email:           newUser.Email,
		DisplayName:     newUser.DisplayName,
		PasswordHash:    newUser.PasswordHash,
		AuthSource:      newUser.AuthSource,
		EmailVerified:   newUser.EmailVerified,
		EmailVerifiedAt: nullTime(newUser.EmailVerifiedAt),
		CreatedAt:       newUser.CreatedAt,
	})
	if err != nil {
		if dbcore.IsUniqueViolation(err) {
			return ErrUserAlreadyExists
		}
		return err
	}

	return nil
}

func (r *UserSQLCRepository) MarkEmailVerified(ctx context.Context, userID string, verifiedAt time.Time) error {
	affected, err := r.queries.MarkUserEmailVerified(ctx, sqldb.MarkUserEmailVerifiedParams{
		ID:              userID,
		EmailVerifiedAt: sql.NullTime{Time: verifiedAt, Valid: true},
	})
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func toDomainUser(row sqldb.GuguAppUser) *domainuser.User {
	user := &domainuser.User{
		ID:            row.ID,
		Email:         row.Email,
		DisplayName:   row.DisplayName,
		PasswordHash:  row.PasswordHash,
		AuthSource:    row.AuthSource,
		EmailVerified: row.EmailVerified,
		CreatedAt:     row.CreatedAt,
	}
	if row.EmailVerifiedAt.Valid {
		user.EmailVerifiedAt = &row.EmailVerifiedAt.Time
	}
	return user
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *value, Valid: true}
}
