package store

// Store defines functions to perform on the store
type Store interface {
	// PutWithKey saves to underlying storage with provided key
	// Returns error on failure
	PutWithKey(key string, payload []byte) error
	// Put saves object to underlying storage.
	// It generates the key using sha-256 algorithm.
	// Returns the key on success otherwise error
	Put(payload []byte) (string, error)
	// Get retrieves the object if object exists otherwise error is returned
	Get(key string) ([]byte, error)
	// Delete evicts the object from the store
	Delete(key string) error
}
