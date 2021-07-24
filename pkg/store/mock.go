package store

import (
	"errors"
	"github.com/fastly/pkg/store/key"
)

// Mock to get an instance of mock Store
func Mock() Store {
	return &mockstore{
		c: make(map[string][]byte),
	}
}

type mockstore struct {
	c map[string][]byte
}

func (m mockstore) PutWithKey(k string, payload []byte) error {
	m.c[k] = payload
	return nil
}

func (m mockstore) Put(payload []byte) (string, error) {
	k := key.Get(payload)
	if _, ok := m.c[k]; ok {
		return "", errors.New("mock: item not stored")
	}
	if err := m.PutWithKey(k, payload); err != nil {
		return "", err
	}
	return k, nil
}

func (m mockstore) Get(k string) ([]byte, error) {
	value, ok := m.c[k]
	if !ok {
		return nil, errors.New("mock: cache miss")
	}
	return value, nil
}

func (m mockstore) Delete(k string) error {
	delete(m.c, k)
	return nil
}
