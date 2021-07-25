package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChunk(t *testing.T) {
	assert.Equal(t, 0, len(chunk([]byte{})))
	assert.Equal(t, 1, len(chunk([]byte("A"))))

	var payload []byte
	for i := 0; i < MaxValueSize+1; i++ {
		payload = append(payload, []byte("A")...)
	}
	assert.Equal(t, 2, len(chunk(payload)))
}

func TestCombine(t *testing.T) {
	assert.Equal(t, 0, len(combine([][]byte{})))

	assert.Equal(t, 2, len(combine([][]byte{
		[]byte("A"),
		[]byte("B"),
	})))

	assert.Equal(t, 3, len(combine([][]byte{
		[]byte("A"),
		[]byte("BC"),
	})))

	var payloadE []byte
	for i := 0; i < MaxValueSize+1; i++ {
		payloadE = append(payloadE, []byte("A")...)
	}
	payload := combine(chunk(payloadE))
	assert.Equal(t, payloadE, payload)
}
