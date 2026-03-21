package product

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func GenerateExternalSKUID(originSKUID string, skuProperties string) string {
	hash := sha256.Sum256([]byte(skuProperties))
	short := hex.EncodeToString(hash[:])[:8]
	return fmt.Sprintf("%s_%s", originSKUID, short)
}
