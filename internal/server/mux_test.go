package server

import (
	"bytes"
	"fmt"
	"github.com/fastly/lib/store"
	gmux "github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupConfig(cfgPath string, env string) (*viper.Viper, error) {
	cfg := viper.New()
	cfg.AddConfigPath(cfgPath + "/envs/" + env)
	cfg.SetConfigName("app")
	if err := cfg.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "Failed to load config")
	}
	return cfg, nil
}

func TestMux(t *testing.T) {
	cfg, _ := setupConfig("../../web/configs", "local")
	m := &mux{
		Router: gmux.NewRouter(),
		cfg:    cfg,
		st:     store.Mock(),
	}

	m0 := m.init()
	assert.NotNil(t, m)
	assert.Equal(t, m, m0)

	var request *http.Request
	var response *httptest.ResponseRecorder

	request, _ = http.NewRequest(http.MethodGet, "/", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)

	request, _ = http.NewRequest(http.MethodPatch, "/", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)

	request, _ = http.NewRequest("GET", "/readiness", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)

	request, _ = http.NewRequest(http.MethodPost, "/", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodPost, "/", strings.NewReader("A"))
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.NotNil(t, response.Body)
	assert.Equal(t, "559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", string(response.Body.Bytes()))

	var payload []byte
	for i := 0; i < store.MaxPayloadSize+1; i++ {
		payload = append(payload, fmt.Sprintf("%d", (i%10))...)
	}
	request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewReader(payload))
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodGet, "/a", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodGet, "/959aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdfff", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusNotFound, response.Code)

	request, _ = http.NewRequest(http.MethodGet, "/559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd.0", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodGet, "/559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, []byte("A"), response.Body.Bytes())

	request, _ = http.NewRequest(http.MethodDelete, "/", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusMethodNotAllowed, response.Code)

	request, _ = http.NewRequest(http.MethodDelete, "/a", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodDelete, "/559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd.0", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusBadRequest, response.Code)

	request, _ = http.NewRequest(http.MethodDelete, "/559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusNoContent, response.Code)

	request, _ = http.NewRequest(http.MethodGet, "/559aead08264d5795d3909718cdd05abd49572e84fe55590eef31a88a08fdffd", nil)
	response = httptest.NewRecorder()
	m.ServeHTTP(response, request)
	assert.Equal(t, http.StatusNotFound, response.Code)
}
