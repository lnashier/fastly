package service

import (
	"github.com/fastly/pkg/store"
	"io/ioutil"
	"net/http"
)

// Provider is the service provider
type Provider struct {
	Store store.Store
}

func (p Provider) Put(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	key, err := p.Store.Put(payload)
	if err != nil {
		// we could return better status based of error types
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	// we could send json object as response if more details need to be shared
	_, _ = w.Write([]byte(key))
}

func (p Provider) Get(w http.ResponseWriter, r *http.Request) {
	key := ""
	payload, err := p.Store.Get(key)
	if err != nil {
		// we could return better status based of error types
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(payload)
}

func (p Provider) Delete(w http.ResponseWriter, r *http.Request) {
	key := ""
	err := p.Store.Delete(key)
	if err != nil {
		// we could return better status based of error types
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
