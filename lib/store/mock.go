package store

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type clientI interface {
	Set(item *memcache.Item) error
	Get(key string) (item *memcache.Item, err error)
	GetMulti(keys []string) (map[string]*memcache.Item, error)
	Delete(key string) error
	Ping() error
}

type mockclient map[string]*memcache.Item

func (m mockclient) Set(item *memcache.Item) error {
	m[item.Key] = item
	return nil
}

func (m mockclient) Get(k string) (item *memcache.Item, err error) {
	item, ok := m[k]
	if !ok {
		return nil, memcache.ErrCacheMiss
	}
	return item, nil
}

func (m mockclient) GetMulti(keys []string) (map[string]*memcache.Item, error) {
	items := make(map[string]*memcache.Item)
	for _, k := range keys {
		item, err := m.Get(k)
		if err != nil {
			return nil, err
		}
		items[k] = item
	}
	return items, nil
}

func (m mockclient) Delete(k string) error {
	if _, ok := m[k]; !ok {
		return memcache.ErrCacheMiss
	}
	delete(m, k)
	return nil
}

func (m mockclient) Ping() error {
	return nil
}
