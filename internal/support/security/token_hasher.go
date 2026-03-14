package security

import (
	"crypto/sha256"
	"encoding/hex"
)

type TokenSHA256Hasher struct{}

func (h TokenSHA256Hasher) Hash(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
