package server

import (
	"fmt"
	"github.com/fastly/pkg/store"
	gmux "github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// mux matches incoming client request to a list of registered handlers
type mux struct {
	*gmux.Router
	cfg *viper.Viper
}

// init is to configure mux
func (m *mux) init() *mux {
	fmt.Println("mux@init enter")
	defer fmt.Println("mux@init exit")

	st := store.Mock()
	/*
		st := memcached.New(memcached.WithStoreAddresses(
			m.cfg.GetStringSlice("store.addresses"),
		))
	*/

	fmt.Printf("mux@init store addresses %v\n", m.cfg.GetStringSlice("store.addresses"))

	m.Methods(http.MethodPost).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("Post Payload %v\n", string(payload))

		key, err := st.Put(payload)
		if err != nil {
			// we could return better status based of error type
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("Key %v\n", key)

		w.Header().Set("Content-Type", "plain/text; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// we could send json object as response if more details need to be shared
		_, _ = w.Write([]byte(key))
	})
	m.Methods(http.MethodGet).Path("/{key}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := gmux.Vars(r)
		key, ok := vars["key"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("Get Key %v\n", key)

		payload, err := st.Get(key)
		if err != nil {
			// we could return better status based of error type
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Printf("Payload %v\n", string(payload))

		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(payload)
	})

	m.Methods(http.MethodDelete).Path("/{key}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := gmux.Vars(r)
		key, ok := vars["key"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("Delete Key %v\n", key)

		if err := st.Delete(key); err != nil {
			// we could return better status based of error type
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return m
}
