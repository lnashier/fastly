package memcached

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/fastly/pkg/store"
	"github.com/fastly/pkg/store/key"
)

// New to get an instance of memcached Store
func New(opts ...Option) store.Store {
	c := &client{
		opts: options{},
	}

	// apply the options
	for _, opt := range opts {
		opt.apply(&c.opts)
	}

	return c.start()
}

// mem is memcache.Client wrapper
type client struct {
	c    *memcache.Client
	opts options
}

func (m client) start() store.Store {
	// connect to memcache server
	m.c = memcache.New(m.opts.storeAddresses...)
	return m
}

func (m client) PutWithKey(k string, payload []byte) error {
	if err := m.c.Add(&memcache.Item{
		Key:        k,
		Value:      payload,
		Expiration: 300,
	}); err != nil {
		return err
	}
	return nil
}

func (m client) Put(payload []byte) (string, error) {
	k := key.Get(payload)
	if err := m.PutWithKey(k, payload); err != nil {
		return "", err
	}
	return k, nil
}

func (m client) Get(k string) ([]byte, error) {
	// retrieve from memcache
	obj, err := m.c.Get(k)
	if err == nil {
		// found in store
		return obj.Value, nil
	}
	return nil, err
}

func (m client) Delete(k string) error {
	return m.c.Delete(k)
}
