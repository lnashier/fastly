package store

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPut(t *testing.T) {
	m := Mock()

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

	k, err = m.Put([]byte("A"))
	assert.Nil(t, err)
	assert.Equal(t, "559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", k)

	// Allow overwrites
	k, err = m.Put([]byte("A"))
	assert.Nil(t, err)
	assert.Equal(t, "559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", k)
}

func TestGet(t *testing.T) {
	m := Mock()

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

	// Malformed key
	payload, err = m.Get("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	assert.Equal(t, ErrNotFound, err)
	assert.Nil(t, payload)

	// Simple object
	payload0 := []byte("A")
	k, _ := m.Put(payload0)
	payload, err = m.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, payload0, payload)

	// Max object
	payload0 = []byte{}
	for i := 0; i < MaxPayloadSize; i++ {
		payload0 = append(payload0, fmt.Sprintf("%d", (i%10))...)
	}
	k, err = m.Put(payload0)
	assert.Nil(t, err)
	assert.NotNil(t, k)
	payload, err = m.Get(k)
	assert.Nil(t, err)
	assert.Equal(t, payload0, payload)
}

func TestDelete(t *testing.T) {
	m := Mock()

	// Empty key
	err := m.Delete("")
	assert.Equal(t, ErrBadKey, err)

	// Malformed key
	err = m.Delete("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85")
	assert.Equal(t, ErrBadKey, err)

	// Malformed key
	err = m.Delete("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855.0")
	assert.Equal(t, ErrBadKey, err)

	// Simple object
	k, _ := m.Put([]byte("A"))
	err = m.Delete(k)
	assert.Nil(t, err)
	_, err = m.Get(k)
	assert.Equal(t, ErrNotFound, err)

	// Max object
	payload0 := []byte{}
	for i := 0; i < MaxPayloadSize; i++ {
		payload0 = append(payload0, fmt.Sprintf("%d", (i%10))...)
	}
	k, _ = m.Put(payload0)
	err = m.Delete(k)
	assert.Nil(t, err)
	_, err = m.Get(k)
	assert.Equal(t, ErrNotFound, err)
}
