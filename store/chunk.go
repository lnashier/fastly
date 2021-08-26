package store

// chunk divides payload into pieces.
// Each chunk can be up to MaxChunkSize.
func chunk(payload []byte) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(payload); i += MaxValueSize {
		chunkEnd := i + MaxValueSize
		if chunkEnd > len(payload) {
			chunkEnd = len(payload)
		}
		chunks = append(chunks, payload[i:chunkEnd])
	}
	return chunks
}

// combine merges chunks into single payload.
func combine(chunks [][]byte) []byte {
	var p []byte
	for _, c := range chunks {
		p = append(p, c...)
	}
	return p
}
