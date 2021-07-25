package store

import (
	"encoding/binary"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"math"
	"strings"
)

// New to get an instance of memcached Store
func New(opts ...Option) Store {
	m := &memclient{
		opts: options{},
	}

	// apply the options
	for _, opt := range opts {
		opt.apply(&m.opts)
	}

	return m.start()
}

// memclient is memcache.Client wrapper
type memclient struct {
	c    *memcache.Client
	opts options
}

func (m memclient) Health() bool {
	if err := m.c.Ping(); err != nil {
		return false
	}
	return true
}

func (m memclient) start() Store {
	// connect to memcache server
	m.c = memcache.New(m.opts.storeAddresses...)
	return m
}

func (m memclient) Put(payload []byte) (string, error) {
	loadSize := binary.Size(payload)

	fmt.Printf("Payload size %d\n", loadSize)

	if loadSize <= 0 {
		return "", ErrTooSmall
	}
	if loadSize > MaxPayloadSize {
		return "", ErrTooLarge
	}
	k := createKey(payload)

	chunks := chunk(payload)

	fmt.Printf("Chunks %d\n", len(chunks))

	var chunksCount uint32 = uint32(len(chunks))

	// All chunks are stored with their own keys.
	// Chunks-count is stored in all chunks as Flags.
	for i, c := range chunks {
		if err := m.c.Add(&memcache.Item{
			Key:   fmt.Sprintf("%s.%d", k, i),
			Value: c,
			Flags: chunksCount,
		}); err != nil {
			return "", err
		}
	}

	return k, nil
}

func (m memclient) Get(k string) ([]byte, error) {
	if strings.Contains(k, ".") || len(k) != 64 {
		return nil, ErrBadKey
	}

	cobj0, err := m.c.Get(fmt.Sprintf("%s.%d", k, 0))
	if err == nil {
		// append first chunk
		chunks := [][]byte{cobj0.Value}
		// all chunks store chunks-count in Flags
		if int(cobj0.Flags) > 1 {
			var cks []string
			for i := 1; i < int(cobj0.Flags); i++ {
				cks = append(cks, fmt.Sprintf("%s.%d", k, i))
			}
			cobjs, err := m.c.GetMulti(cks)
			if err != nil {
				return nil, err
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
	return nil, err
}

func (m memclient) Delete(k string) error {
	if strings.Contains(k, ".") || len(k) != 64 {
		return ErrBadKey
	}

	for i := 0; i < int(math.Ceil(MaxPayloadSize/MaxValueSize)); i++ {
		// Let's delete all the chunks for this key
		if err := m.c.Delete(fmt.Sprintf("%s.%d", k, i)); err != nil {
			if err == memcache.ErrCacheMiss {
				// if first chunk is not found that means key does not exist
				if i == 0 {
					return err
				}
				// otherwise all possible chunks are deleted for this key
				return nil
			}
			// something bad happened
			return err
		}
	}
	return nil
}
