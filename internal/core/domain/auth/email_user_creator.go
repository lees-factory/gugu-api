package auth

import (
	"context"
	"fmt"

	"github.com/ljj/gugu-api/internal/core/domain/user"
)

type CreateEmailUserInput struct {
	Email        string
	DisplayName  string
	PasswordHash string
}

type EmailUserCreator interface {
	Create(ctx context.Context, input CreateEmailUserInput) (*user.User, error)
}

type emailUserCreator struct {
	userFinder  user.Finder
	userCreator user.Creator
}

func NewEmailUserCreator(userFinder user.Finder, userCreator user.Creator) EmailUserCreator {
	return &emailUserCreator{
		userFinder:  userFinder,
		userCreator: userCreator,
	}
}

func (c *emailUserCreator) Create(ctx context.Context, input CreateEmailUserInput) (*user.User, error) {
	emailValue := normalizeEmail(input.Email)
	foundUser, err := c.userFinder.FindByEmail(ctx, emailValue)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	if foundUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	newUser, err := c.userCreator.Create(ctx, user.CreateInput{
		Email:         emailValue,
		DisplayName:   input.DisplayName,
		PasswordHash:  input.PasswordHash,
		AuthSource:    "email",
		EmailVerified: false,
	})
	if err != nil {
		return nil, fmt.Errorf("create email user: %w", err)
	}

	return newUser, nil
}
