package ip_test

import (
	"ifconfig/ip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestGetPlaintext(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Get)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestGetJson(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Get)

	req, _ := http.NewRequest("GET", "/?json", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "{\"ip\":\"192.168.1.124\"}\n", rr.Body.String(), "should return an ip address in json format")
}

func TestGetNoHeader(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Get)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}

func TestGetMultipleIp(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Get)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124,10.0.0.1")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestGetNoIp(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Get)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}
