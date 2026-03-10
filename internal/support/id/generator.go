package id

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type RandomHexGenerator struct {
	size int
}

func NewRandomHexGenerator(size int) RandomHexGenerator {
	return RandomHexGenerator{size: size}
}

func (g RandomHexGenerator) New() (string, error) {
	buf := make([]byte, g.size)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	return hex.EncodeToString(buf), nil
}
