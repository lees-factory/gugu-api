package security

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

type NumericCodeGenerator struct {
	length int
}

func NewNumericCodeGenerator(length int) NumericCodeGenerator {
	return NumericCodeGenerator{length: length}
}

func (g NumericCodeGenerator) New() (string, error) {
	if g.length <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	var builder strings.Builder
	builder.Grow(g.length)

	for range g.length {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("generate numeric code: %w", err)
		}
		builder.WriteByte(byte('0' + n.Int64()))
	}

	return builder.String(), nil
}
