package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type UserAccountReader interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	FindByID(ctx context.Context, userID string) (*user.User, error)
}

type userAccountReader struct {
	userFinder user.Finder
}

func NewUserAccountReader(userFinder user.Finder) UserAccountReader {
	return &userAccountReader{userFinder: userFinder}
}

func (r *userAccountReader) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	foundUser, err := r.userFinder.FindByEmail(ctx, normalizeEmail(email))
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return foundUser, nil
}

func (r *userAccountReader) FindByID(ctx context.Context, userID string) (*user.User, error) {
	foundUser, err := r.userFinder.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return foundUser, nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
