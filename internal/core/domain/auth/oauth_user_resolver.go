package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type FindOrCreateOAuthUserInput struct {
	Email       string
	DisplayName string
	Provider    Provider
	VerifiedAt  time.Time
}

type OAuthUserResolver interface {
	FindOrCreate(ctx context.Context, input FindOrCreateOAuthUserInput) (*user.User, error)
}

type oauthUserResolver struct {
	userFinder  user.Finder
	userCreator user.Creator
}

func NewOAuthUserResolver(userFinder user.Finder, userCreator user.Creator) OAuthUserResolver {
	return &oauthUserResolver{
		userFinder:  userFinder,
		userCreator: userCreator,
	}
}

func (r *oauthUserResolver) FindOrCreate(ctx context.Context, input FindOrCreateOAuthUserInput) (*user.User, error) {
	emailValue := normalizeEmail(input.Email)
	foundUser, err := r.userFinder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return foundUser, nil
	}

	newUser, err := r.userCreator.Create(ctx, user.CreateInput{
		Email:           emailValue,
		DisplayName:     input.DisplayName,
		AuthSource:      string(input.Provider),
		EmailVerified:   true,
		EmailVerifiedAt: &input.VerifiedAt,
	})
	if err != nil {
		return nil, fmt.Errorf("create oauth user: %w", err)
	}

	return newUser, nil
}
