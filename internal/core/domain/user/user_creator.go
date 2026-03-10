package user

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type CreateInput struct {
	Email           string
	DisplayName     string
	PasswordHash    string
	AuthSource      string
	EmailVerified   bool
	EmailVerifiedAt *time.Time
}

type Creator interface {
	Create(ctx context.Context, input CreateInput) (*User, error)
}

type creator struct {
	writer      Writer
	idGenerator IDGenerator
	clock       Clock
}

func NewCreator(writer Writer, idGenerator IDGenerator, clock Clock) Creator {
	return &creator{
		writer:      writer,
		idGenerator: idGenerator,
		clock:       clock,
	}
}

func (c *creator) Create(ctx context.Context, input CreateInput) (*User, error) {
	userID, err := c.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate user id: %w", err)
	}

	now := c.clock.Now()
	newUser := User{
		ID:            userID,
		Email:         strings.TrimSpace(strings.ToLower(input.Email)),
		DisplayName:   strings.TrimSpace(input.DisplayName),
		PasswordHash:  input.PasswordHash,
		AuthSource:    strings.TrimSpace(strings.ToLower(input.AuthSource)),
		EmailVerified: input.EmailVerified,
		CreatedAt:     now,
	}
	if input.EmailVerified {
		verifiedAt := now
		if input.EmailVerifiedAt != nil {
			verifiedAt = *input.EmailVerifiedAt
		}
		newUser.EmailVerifiedAt = &verifiedAt
	}

	if err := c.writer.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &newUser, nil
}
