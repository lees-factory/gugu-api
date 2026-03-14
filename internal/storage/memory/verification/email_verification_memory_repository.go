package verification

import (
	"context"
	"errors"
	"sync"
	"time"

	domainverification "github.com/ljj/gugu-api/internal/core/domain/verification"
)

type EmailVerificationMemoryRepository struct {
	mu            sync.RWMutex
	verifications map[string]domainverification.EmailVerification
}

func NewRepository() *EmailVerificationMemoryRepository {
	return &EmailVerificationMemoryRepository{
		verifications: make(map[string]domainverification.EmailVerification),
	}
}

func (r *EmailVerificationMemoryRepository) Create(_ context.Context, emailVerification domainverification.EmailVerification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.verifications[emailVerification.Code] = emailVerification
	return nil
}

func (r *EmailVerificationMemoryRepository) FindByCode(_ context.Context, code string) (*domainverification.EmailVerification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	foundVerification, ok := r.verifications[code]
	if !ok {
		return nil, nil
	}

	return &foundVerification, nil
}

func (r *EmailVerificationMemoryRepository) MarkUsed(_ context.Context, code string, usedAt time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	foundVerification, ok := r.verifications[code]
	if !ok {
		return errors.New("verification not found")
	}

	foundVerification.UsedAt = &usedAt
	r.verifications[code] = foundVerification
	return nil
}
