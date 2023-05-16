package store

// Option configures how we set up the client.
type Option interface {
	apply(*options)
}

// WithStoreAddresses returns an Option which sets the store address
func WithStoreAddresses(addresses []string) Option {
	return newFuncOption(func(o *options) {
		o.storeAddresses = addresses
	})
}

type options struct {
	storeAddresses []string
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fdo *funcOption) apply(do *options) {
	fdo.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}
