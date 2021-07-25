package store

import (
	"encoding/binary"
	"errors"
	"fmt"
)

// Mock to get an instance of mock Store
func Mock() Store {
	return make(mockstore)
}

type mockstore map[string][]byte

func (m mockstore) Put(payload []byte) (string, error) {
	loadSize := binary.Size(payload)
	fmt.Printf("Payload size %d\n", loadSize)
	if loadSize <= 0 {
		return "", ErrTooSmall
	}
	if loadSize > MaxPayloadSize {
		return "", ErrTooLarge
	}
	k := createKey(payload)
	if _, ok := m[k]; ok {
		return "", errors.New("mock: item not stored")
	}
	m[k] = payload
	return k, nil
}

func (m mockstore) Get(k string) ([]byte, error) {
	payload, ok := m[k]
	if !ok {
		return nil, errors.New("mock: cache miss")
	}
	fmt.Printf("Payload size %d\n", binary.Size(payload))
	return payload, nil
}

func (m mockstore) Delete(k string) error {
	if _, ok := m[k]; !ok {
		return errors.New("mock: cache miss")
	}
	delete(m, k)
	return nil
}
