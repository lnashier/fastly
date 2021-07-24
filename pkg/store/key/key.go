package key

import (
	"crypto/sha256"
	"encoding/hex"
)

// Get creates a key by performing sha256 on the payload
func Get(payload []byte) string {
	hash := sha256.Sum256(payload)
	return hex.EncodeToString(hash[:])
}
