package server

import (
	"fmt"
	"github.com/fastly/lib/store"
	gmux "github.com/gorilla/mux"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

// mux matches incoming client request to a list of registered handlers
type mux struct {
	*gmux.Router
	cfg *viper.Viper
	st  *store.Store
}

// init is to configure mux
func (m *mux) init() *mux {
	fmt.Println("mux@init enter")
	defer fmt.Println("mux@init exit")

	m.Methods(http.MethodGet).Path("/readiness").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := m.st.Health(); ok {
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusFailedDependency)
	})

	m.Methods(http.MethodPost).Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload []byte
		if r.Body != nil {
			// if content is more than allowed
			//r.Body = http.MaxBytesReader(w, r.Body, store.MaxPayloadSize+1024)
			contentType := r.Header.Get("Content-Type")
			switch {
			case strings.Contains(contentType, "multipart/form-data"):
				// no need to load all parts
				mr, err := r.MultipartReader()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// expecting one part
				part, err := mr.NextPart()
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				defer func() {
					_ = part.Close()
				}()
				// if it is not a file
				if len(part.FileName()) < 1 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				// get the contents
				payload, err = ioutil.ReadAll(part)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			case contentType == "application/octet-stream", contentType == "", contentType == "text/plain":
				var err error
				payload, err = ioutil.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					// we could return a status message too
					return
				}
			default:
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}
		}
		key, err := m.st.Put(payload)
		if err != nil {
			fmt.Printf("mux@post error %s\n", err.Error())
			switch err {
			case store.ErrTooSmall, store.ErrTooLarge, store.ErrNotStored:
				w.WriteHeader(http.StatusBadRequest)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			// we could return a status message too
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
			// we could return a status message too
			return
		}

		fmt.Printf("Get Key %v\n", key)

		payload, err := m.st.Get(key)
		if err != nil {
			fmt.Printf("mux@get error %s\n", err.Error())
			switch err {
			case store.ErrBadKey:
				w.WriteHeader(http.StatusBadRequest)
			case store.ErrCorruptedContent, store.ErrNotFound:
				// we should monitor ErrCorruptedContent
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			// we could return a status message too
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
			// we could return a status message too
			return
		}

		fmt.Printf("Delete Key %v\n", key)

		if err := m.st.Delete(key); err != nil {
			fmt.Printf("mux@delete error %s\n", err.Error())
			switch err {
			case store.ErrBadKey:
				w.WriteHeader(http.StatusBadRequest)
			case store.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			// we could return a status message too
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return m
}
