package server

import (
	"fmt"
	"github.com/fastly/lib/store"
	gmux "github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// mux matches incoming client request to a list of registered handlers
type mux struct {
	*gmux.Router
	cfg *viper.Viper
	st  store.Store
}

// init is to configure mux
func (m *mux) init() *mux {
	fmt.Println("mux@init enter")
	defer fmt.Println("mux@init exit")

	m.Methods(http.MethodPost).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		key, err := m.st.Put(payload)
		if err != nil {
			fmt.Printf("mux@post error %s\n", err.Error())
			// we could return better status based of error type
			w.WriteHeader(http.StatusInternalServerError)
			// we could return status message too
			return
		}

		fmt.Printf("Post Key %v\n", key)

		w.Header().Set("Content-Type", "plain/text; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		// we could send some serialized object in the response
		// if more details need to be shared
		_, _ = w.Write([]byte(key))
	})
	m.Methods(http.MethodGet).Path("/{key}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := gmux.Vars(r)
		key, ok := vars["key"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			// we could return status message too
			return
		}

		fmt.Printf("Get Key %v\n", key)

		payload, err := m.st.Get(key)
		if err != nil {
			fmt.Printf("mux@get error %s\n", err.Error())
			// we could return better status based of error type
			w.WriteHeader(http.StatusInternalServerError)
			// we could return status message too
			return
		}

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

		if err := m.st.Delete(key); err != nil {
			fmt.Printf("mux@delete error %s\n", err.Error())
			// we could return better status based of error type
			w.WriteHeader(http.StatusNotFound)
			// we could return status message too
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return m
}
