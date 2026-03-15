package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	domainauth "github.com/ljj/gugu-api/internal/core/domain/auth"
	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
)

type LoginSessionSQLCRepository struct {
	queries *sqldb.Queries
}

func NewLoginSessionSQLCRepository(db *sql.DB) *LoginSessionSQLCRepository {
	return &LoginSessionSQLCRepository{queries: sqldb.New(db)}
}

func (r *LoginSessionSQLCRepository) Create(ctx context.Context, session domainauth.LoginSession) error {
	return r.queries.CreateUserLoginSession(ctx, sqldb.CreateUserLoginSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshTokenHash: session.RefreshTokenHash,
		TokenFamilyID:    session.TokenFamilyID,
		ParentSessionID:  nullString(session.ParentSessionID),
		UserAgent:        session.UserAgent,
		ClientIp:         session.ClientIP,
		DeviceName:       session.DeviceName,
		ExpiresAt:        session.ExpiresAt,
		LastSeenAt:       session.LastSeenAt,
		RotatedAt:        nullTime(session.RotatedAt),
		RevokedAt:        nullTime(session.RevokedAt),
		ReuseDetectedAt:  nullTime(session.ReuseDetectedAt),
		CreatedAt:        session.CreatedAt,
	})
}

func (r *LoginSessionSQLCRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*domainauth.LoginSession, error) {
	row, err := r.queries.FindUserLoginSessionByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomainLoginSession(row), nil
}

func (r *LoginSessionSQLCRepository) MarkRotated(ctx context.Context, sessionID string, rotatedAt time.Time) error {
	affected, err := r.queries.MarkUserLoginSessionRotated(ctx, sqldb.MarkUserLoginSessionRotatedParams{
		ID:        sessionID,
		RotatedAt: sql.NullTime{Time: rotatedAt, Valid: true},
	})
	return mapAffectedRows(affected, err)
}

func (r *LoginSessionSQLCRepository) Revoke(ctx context.Context, sessionID string, revokedAt time.Time) error {
	affected, err := r.queries.RevokeUserLoginSession(ctx, sqldb.RevokeUserLoginSessionParams{
		ID:        sessionID,
		RevokedAt: sql.NullTime{Time: revokedAt, Valid: true},
	})
	return mapAffectedRows(affected, err)
}

func (r *LoginSessionSQLCRepository) RevokeFamily(ctx context.Context, tokenFamilyID string, revokedAt time.Time) error {
	affected, err := r.queries.RevokeUserLoginSessionFamily(ctx, sqldb.RevokeUserLoginSessionFamilyParams{
		TokenFamilyID: tokenFamilyID,
		RevokedAt:     sql.NullTime{Time: revokedAt, Valid: true},
	})
	return mapAffectedRows(affected, err)
}

func (r *LoginSessionSQLCRepository) MarkReuseDetected(ctx context.Context, sessionID string, detectedAt time.Time) error {
	affected, err := r.queries.MarkUserLoginSessionReuseDetected(ctx, sqldb.MarkUserLoginSessionReuseDetectedParams{
		ID:              sessionID,
		ReuseDetectedAt: sql.NullTime{Time: detectedAt, Valid: true},
	})
	return mapAffectedRows(affected, err)
}

func (r *LoginSessionSQLCRepository) UpdateLastSeen(ctx context.Context, sessionID string, lastSeenAt time.Time) error {
	affected, err := r.queries.UpdateUserLoginSessionLastSeen(ctx, sqldb.UpdateUserLoginSessionLastSeenParams{
		ID:         sessionID,
		LastSeenAt: lastSeenAt,
	})
	return mapAffectedRows(affected, err)
}

func toDomainLoginSession(row sqldb.GuguUserLoginSession) *domainauth.LoginSession {
	session := &domainauth.LoginSession{
		ID:               row.ID,
		UserID:           row.UserID,
		RefreshTokenHash: row.RefreshTokenHash,
		TokenFamilyID:    row.TokenFamilyID,
		UserAgent:        row.UserAgent,
		ClientIP:         row.ClientIp,
		DeviceName:       row.DeviceName,
		ExpiresAt:        row.ExpiresAt,
		LastSeenAt:       row.LastSeenAt,
		CreatedAt:        row.CreatedAt,
	}
	if row.ParentSessionID.Valid {
		session.ParentSessionID = &row.ParentSessionID.String
	}
	if row.RotatedAt.Valid {
		session.RotatedAt = &row.RotatedAt.Time
	}
	if row.RevokedAt.Valid {
		session.RevokedAt = &row.RevokedAt.Time
	}
	if row.ReuseDetectedAt.Valid {
		session.ReuseDetectedAt = &row.ReuseDetectedAt.Time
	}
	return session
}

func mapAffectedRows(affected int64, err error) error {
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func nullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func nullTime(value *time.Time) sql.NullTime {
	if value == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *value, Valid: true}
}
