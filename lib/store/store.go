package store

import (
	"encoding/binary"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/pkg/errors"
	"math"
	"strings"
)

const (
	// MinPayloadSize is minimum payload size allowed
	// 1 byte
	MinPayloadSize = 1

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
	// ErrTooSmall means that item was too small < MinPayloadSize
	ErrTooSmall = errors.New("store: item too small")

	// ErrTooLarge means that item was too large > MaxPayloadSize
	ErrTooLarge = errors.New("store: item too large")

	// ErrBadKey means that key is invalid
	ErrBadKey = errors.New("store: invalid key")

	// ErrCorruptedContent means that content got corrupted
	ErrCorruptedContent = errors.New("store: content corrupted")

	// ErrNotStored means that payload was not stored
	ErrNotStored = errors.New("store: not stored")

	// ErrNotFound means that key does not exist
	ErrNotFound = errors.New("store: not found")

	// ErrOpFailed means that some error occurred
	ErrOpFailed = errors.New("store: op failed")
)

// New to get an instance of Store
func New(opts ...Option) *Store {
	m := &Store{
		opts: options{},
	}
	// apply the options
	for _, opt := range opts {
		opt.apply(&m.opts)
	}
	m.c = memcache.New(m.opts.storeAddresses...)
	return m
}

// Mock to get an instance of mocked Store
func Mock() *Store {
	return &Store{c: make(mockclient)}
}

// Store defines functions to perform on the store
type Store struct {
	c    clientI
	opts options
}

// Health tells if store is live and healthy
func (s Store) Health() bool {
	if err := s.c.Ping(); err != nil {
		return false
	}
	return true
}

// Put saves object to underlying storage.
// It generates the key using sha-256 algorithm.
// Returns the key on success otherwise error
func (s Store) Put(payload []byte) (string, error) {
	loadSize := binary.Size(payload)

	fmt.Printf("Payload size %d\n", loadSize)

	if loadSize < MinPayloadSize {
		return "", ErrTooSmall
	}
	if loadSize > MaxPayloadSize {
		return "", ErrTooLarge
	}
	k := createKey(payload)

	chunks := chunk(payload)

	fmt.Printf("Chunks %d\n", len(chunks))

	chunksCount := uint32(len(chunks))

	// All chunks are stored with their own keys.
	// Chunks-count is stored in all chunks as Flags.
	for i, c := range chunks {
		if err := s.c.Set(&memcache.Item{
			Key:   fmt.Sprintf("%s.%d", k, i),
			Value: c,
			Flags: chunksCount,
		}); err != nil {
			fmt.Printf("store#Put error %s\n", err.Error())
			switch err {
			case memcache.ErrNotStored:
				return "", ErrNotStored
			default:
				return "", ErrOpFailed
			}
		}
	}

	return k, nil
}

// Get retrieves the object if object exists otherwise error is returned
func (s Store) Get(k string) ([]byte, error) {
	if strings.Contains(k, ".") || len(k) != 64 {
		return nil, ErrBadKey
	}

	cobj0, err := s.c.Get(fmt.Sprintf("%s.%d", k, 0))
	if err == nil {
		// append first chunk
		chunks := [][]byte{cobj0.Value}
		// all chunks store chunks-count in Flags
		if int(cobj0.Flags) > 1 {
			var cks []string
			for i := 1; i < int(cobj0.Flags); i++ {
				cks = append(cks, fmt.Sprintf("%s.%d", k, i))
			}
			cobjs, err := s.c.GetMulti(cks)
			if err != nil {
				fmt.Printf("store#Get chunks error %s\n", err.Error())
				if err != memcache.ErrCacheMiss {
					return nil, ErrOpFailed
				}
			}
			// collect all the chunks in correct order
			for _, ck := range cks {
				chunks = append(chunks, cobjs[ck].Value)
			}
		}
		// combine to form the original payload
		payload := combine(chunks)
		// verify consistency of the payload by recreating the key
		if pk := createKey(payload); pk != k {
			return nil, ErrCorruptedContent
		}
		return payload, nil
	}
	fmt.Printf("store#Get first chunk error %s\n", err.Error())
	switch err {
	case memcache.ErrCacheMiss:
		return nil, ErrNotFound
	default:
		return nil, ErrOpFailed
	}
}

// Delete evicts the object from the store
func (s Store) Delete(k string) error {
	if strings.Contains(k, ".") || len(k) != 64 {
		return ErrBadKey
	}

	for i := 0; i <= int(math.Ceil(MaxPayloadSize/MaxValueSize)); i++ {
		// Let's delete all the chunks for this key
		if err := s.c.Delete(fmt.Sprintf("%s.%d", k, i)); err != nil {
			fmt.Printf("store#Delete chunk %d error %s\n", i, err.Error())
			if err == memcache.ErrCacheMiss {
				// if first chunk is not found that means key does not exist
				if i == 0 {
					switch err {
					case memcache.ErrCacheMiss:
						return ErrNotFound
					default:
						return ErrOpFailed
					}
				}
				// otherwise all possible chunks are deleted for this key
				return nil
			}
			return ErrOpFailed
		}
	}
	return nil
}
