package user

import (
	"context"
	"fmt"
	"strings"
)

type Creator interface {
	Create(ctx context.Context, newUser NewUser) (*User, error)
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

func (c *creator) Create(ctx context.Context, newUser NewUser) (*User, error) {
	userID, err := c.idGenerator.New()
	if err != nil {
		return nil, fmt.Errorf("generate user id: %w", err)
	}

	now := c.clock.Now()
	u := User{
		ID:            userID,
		Email:         strings.TrimSpace(strings.ToLower(newUser.Email)),
		DisplayName:   strings.TrimSpace(newUser.DisplayName),
		PasswordHash:  newUser.PasswordHash,
		AuthSource:    strings.TrimSpace(strings.ToLower(newUser.AuthSource)),
		EmailVerified: newUser.EmailVerified,
		CreatedAt:     now,
	}
	if newUser.EmailVerified {
		verifiedAt := now
		if newUser.EmailVerifiedAt != nil {
			verifiedAt = *newUser.EmailVerifiedAt
		}
		u.EmailVerifiedAt = &verifiedAt
	}

	if err := c.writer.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &u, nil
}
