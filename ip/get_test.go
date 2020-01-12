package ip

import "testing"
import "github.com/go-chi/chi"
import "github.com/stretchr/testify/assert"
import "net/http"
import "net/http/httptest"

func TestGetPlaintext(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Get)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestGetJson(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Get)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}

func TestGetMultipleIp(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Get)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124,10.0.0.1")
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestGetNoIp(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/", Get)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}
