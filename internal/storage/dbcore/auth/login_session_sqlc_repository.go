package auth

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ljj/gugu-api/internal/storage/dbcore/sqldb"
	supportauth "github.com/ljj/gugu-api/internal/support/auth"
)

type LoginSessionSQLCRepository struct {
	queries *sqldb.Queries
}

func NewLoginSessionSQLCRepository(db *sql.DB) *LoginSessionSQLCRepository {
	return &LoginSessionSQLCRepository{queries: sqldb.New(db)}
}

func (r *LoginSessionSQLCRepository) Create(ctx context.Context, session supportauth.LoginSession) error {
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

func (r *LoginSessionSQLCRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*supportauth.LoginSession, error) {
	row, err := r.queries.FindUserLoginSessionByRefreshTokenHash(ctx, refreshTokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return toDomainLoginSession(row), nil
}

func (r *LoginSessionSQLCRepository) ListActiveByUserID(ctx context.Context, userID string, now time.Time) ([]supportauth.LoginSession, error) {
	rows, err := r.queries.ListActiveUserLoginSessionsByUserID(ctx, sqldb.CountActiveUserLoginSessionsByUserIDParams{
		UserID:    userID,
		ExpiresAt: now,
	})
	if err != nil {
		return nil, err
	}

	items := make([]supportauth.LoginSession, 0, len(rows))
	for _, row := range rows {
		items = append(items, *toDomainLoginSession(row))
	}
	return items, nil
}

func (r *LoginSessionSQLCRepository) CountActiveByUserID(ctx context.Context, userID string, now time.Time) (int, error) {
	count, err := r.queries.CountActiveUserLoginSessionsByUserID(ctx, sqldb.CountActiveUserLoginSessionsByUserIDParams{
		UserID:    userID,
		ExpiresAt: now,
	})
	if err != nil {
		return 0, err
	}
	return int(count), nil
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

func (r *LoginSessionSQLCRepository) RevokeByUserSessionID(ctx context.Context, userID string, sessionID string, revokedAt time.Time) error {
	return r.queries.RevokeUserLoginSessionByUserIDSessionID(ctx, sqldb.RevokeUserLoginSessionByUserIDSessionIDParams{
		UserID:    userID,
		ID:        sessionID,
		RevokedAt: sql.NullTime{Time: revokedAt, Valid: true},
	})
}

func (r *LoginSessionSQLCRepository) RevokeOldestActiveByUserID(ctx context.Context, userID string, now time.Time, revokedAt time.Time) error {
	affected, err := r.queries.RevokeOldestActiveUserLoginSessionByUserID(ctx, sqldb.RevokeOldestActiveUserLoginSessionByUserIDParams{
		UserID:    userID,
		RevokedAt: sql.NullTime{Time: revokedAt, Valid: true},
		ExpiresAt: now,
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

func toDomainLoginSession(row sqldb.GuguUserLoginSession) *supportauth.LoginSession {
	session := &supportauth.LoginSession{
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
