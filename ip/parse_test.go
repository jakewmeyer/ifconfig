package ip_test

import (
	"ifconfig/ip"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestReturnPlaintext(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Parse)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestParseJson(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Parse)

	req, _ := http.NewRequest("GET", "/?json", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "{\"ip\":\"192.168.1.124\"}\n", rr.Body.String(), "should return an ip address in json format")
}

func TestParseNoHeader(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Parse)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return a 500")
}

func TestParseMultipleIp(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Parse)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124,10.0.0.1")

	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestParseNoIp(t *testing.T) {
	t.Parallel()

	r := chi.NewRouter()
	r.Get("/", ip.Parse)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code, "should return a 500")
}
