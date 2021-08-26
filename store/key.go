package store

import (
	"crypto/sha256"
	"encoding/hex"
)

// createKey creates a key by performing sha256 on the payload
func createKey(payload []byte) string {
	h := sha256.Sum256(payload)
	return hex.EncodeToString(h[:])
}
