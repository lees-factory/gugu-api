package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type RandomTokenGenerator struct {
	size int
}

func NewRandomTokenGenerator(size int) RandomTokenGenerator {
	return RandomTokenGenerator{size: size}
}

func (g RandomTokenGenerator) New() (string, error) {
	buf := make([]byte, g.size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
