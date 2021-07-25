package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPut(t *testing.T) {
	m := &memclient{}

	// Payload too small
	k, err := m.Put([]byte{})
	assert.Equal(t, ErrTooSmall, err)
	assert.Equal(t, "", k)

	// Payload too large
	var payload []byte
	for i := 0; i < MaxPayloadSize+1; i++ {
		payload = append(payload, fmt.Sprintf("%d", (i%10))...)
	}
	k, err = m.Put(payload)
	assert.Equal(t, ErrTooLarge, err)
	assert.Equal(t, "", k)
}

func TestGet(t *testing.T) {
	m := &memclient{}

	// Empty key
	payload, err := m.Get("")
	assert.Equal(t, ErrBadKey, err)
	assert.Nil(t, payload)

	// Malformed key
	payload, err = m.Get("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85")
	assert.Equal(t, ErrBadKey, err)
	assert.Nil(t, payload)

	// Malformed key
	payload, err = m.Get("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855.0")
	assert.Equal(t, ErrBadKey, err)
	assert.Nil(t, payload)
}

func TestDelete(t *testing.T) {
	m := &memclient{}

	// Empty key
	err := m.Delete("")
	assert.Equal(t, ErrBadKey, err)

	// Malformed key
	err = m.Delete("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85")
	assert.Equal(t, ErrBadKey, err)

	// Malformed key
	err = m.Delete("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855.0")
	assert.Equal(t, ErrBadKey, err)
}
