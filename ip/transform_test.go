package ip

import "testing"
import "github.com/julienschmidt/httprouter"
import "github.com/stretchr/testify/assert"
import "net/http"
import "net/http/httptest"

func TestTransformPlaintext(t *testing.T) {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", Transform)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestTransformJson(t *testing.T) {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", Transform)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}

func TestTransformMultipleIp(t *testing.T) {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", Transform)

	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("x-forwarded-for", "192.168.1.124,10.0.0.1")
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "should return a 200 status")
	assert.Equal(t, "192.168.1.124", rr.Body.String(), "should return an ip address")
}

func TestTransformNoIp(t *testing.T) {
	router := httprouter.New()
	router.HandlerFunc("GET", "/", Transform)

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code, "should return a 400 status")
}
