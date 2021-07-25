package store

import (
	"github.com/pkg/errors"
)

const (
	// MaxPayloadSize is maximum payload size allowed
	// 50 mebibyte
	MaxPayloadSize = 50 * 1024 * 1024

	// MaxChunkSize is maximum size allowed for a single chunk
	// 1 mebibyte
	MaxChunkSize = 1 * 1024 * 1024

	// MaxValueSize is maximum size allowed for item value
	// https://github.com/memcached/memcached/blob/c472369fed5981ba8c004d426cee62d5165c47ca/proto_text.c#L1368
	// https://github.com/memcached/memcached/blob/c472369fed5981ba8c004d426cee62d5165c47ca/items.c#L366
	MaxValueSize = MaxChunkSize - 130
)

var (
	// ErrTooSmall means that item was too small <= 0
	ErrTooSmall = errors.New("store: item too small")

	// ErrTooLarge means that item was too large > MaxPayloadSize
	ErrTooLarge = errors.New("store: item too large")

	// ErrBadKey means that key is invalid
	ErrBadKey = errors.New("store: invalid key")

	// ErrCorruptedContent means that content got corrupted
	ErrCorruptedContent = errors.New("store: content corrupted")

	// ErrNoContent means that key does not exist
	ErrNoContent = errors.New("store: no content")
)

// Store defines functions to perform on the store
type Store interface {
	// Put saves object to underlying storage.
	// It generates the key using sha-256 algorithm.
	// Returns the key on success otherwise error
	Put(payload []byte) (string, error)
	// Get retrieves the object if object exists otherwise error is returned
	Get(key string) ([]byte, error)
	// Delete evicts the object from the store
	Delete(key string) error
	// Health tells if store is live and healthy
	Health() bool
}
